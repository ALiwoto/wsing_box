package pptp

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"

	_ "github.com/sagernet/sing/common/buf"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/logger"
	M "github.com/sagernet/sing/common/metadata"
)

// PPTP message types
const (
	StartControlConnectionRequest uint16 = 1
	StartControlConnectionReply   uint16 = 2
	OutgoingCallRequest           uint16 = 7
	OutgoingCallReply             uint16 = 8
	CallClearRequest              uint16 = 12
	CallDisconnectNotify          uint16 = 13
)

const (
	GREProtocolPPP = 0x880B
)

// GRE header for PPTP (Enhanced GRE)
type GREHeader struct {
	FlagsAndVersion uint16 // Flags and version
	Protocol        uint16 // Protocol type
	PayloadLength   uint16 // Payload length
	CallID          uint16 // Call ID
	Seq             uint32 // Sequence number (optional)
	Ack             uint32 // Acknowledgment number (optional)
}

// PPTP header structure
type PPTPHeader struct {
	Length      uint16 // Total message length in octets
	MessageType uint16 // Type of message
	MagicCookie uint32 // Always 0x1A2B3C4D
	ControlType uint16 // Control message type
	Reserved0   uint16 // Reserved for future use (must be 0)
}

type PPTPClient struct {
	serverAddr M.Socksaddr
	username   string
	password   string
	logger     logger.Logger

	controlConn net.Conn
	mutex       sync.Mutex
	connected   bool
	callID      uint16
	peerCallID  uint16
	// Other necessary fields for PPTP
}

func NewPPTPClient(serverAddr M.Socksaddr, username, password string, logger logger.Logger) *PPTPClient {
	return &PPTPClient{
		serverAddr: serverAddr,
		username:   username,
		password:   password,
		logger:     logger,
		callID:     1, // Start with call ID 1
	}
}

func (c *PPTPClient) Connect(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.connected {
		return nil
	}

	// Connect to PPTP server on port 1723 (TCP)
	var err error
	c.controlConn, err = net.DialTimeout("tcp", c.serverAddr.String(), 10*time.Second)
	if err != nil {
		return E.Cause(err, "connect to PPTP server")
	}

	// PPTP control connection setup:
	// 1. Send Start-Control-Connection-Request
	if err := c.sendStartControlConnectionRequest(); err != nil {
		c.controlConn.Close()
		return err
	}

	// 2. Receive Start-Control-Connection-Reply
	if err := c.receiveStartControlConnectionReply(); err != nil {
		c.controlConn.Close()
		return err
	}

	// 3. Send Outgoing-Call-Request
	if err := c.sendOutgoingCallRequest(); err != nil {
		c.controlConn.Close()
		return err
	}

	// 4. Receive Outgoing-Call-Reply
	if err := c.receiveOutgoingCallReply(); err != nil {
		c.controlConn.Close()
		return err
	}

	c.connected = true

	// Start a goroutine to handle control connection messages
	go c.handleControlMessages()

	return nil
}

// Placeholder methods - these need to be implemented with actual PPTP protocol logic
func (c *PPTPClient) sendStartControlConnectionRequest() error {
	// Create the PPTP Start-Control-Connection-Request message
	var message struct {
		Header          PPTPHeader
		ProtocolVersion uint16   // PPTP protocol version
		Reserved1       uint16   // Reserved (must be 0)
		FramingCap      uint32   // Framing capabilities
		BearerCap       uint32   // Bearer capabilities
		MaxChannels     uint16   // Maximum channels
		FirmwareRev     uint16   // Firmware revision
		HostName        [64]byte // Host name
		VendorString    [64]byte // Vendor string
	}

	// Fill header
	message.Header.Length = 156             // Total message length (156 bytes)
	message.Header.MessageType = 1          // Control message
	message.Header.MagicCookie = 0x1A2B3C4D // Magic cookie
	message.Header.ControlType = StartControlConnectionRequest

	// Fill message fields
	message.ProtocolVersion = 0x0100 // Version 1.0
	message.FramingCap = 1           // Async framing supported
	message.BearerCap = 1            // Analog bearer supported
	message.MaxChannels = 1          // Request single channel
	message.FirmwareRev = 1          // Firmware revision 1

	// Copy hostname (up to 64 bytes, null-terminated)
	hostname := []byte("sing-box-pptp-client")
	copy(message.HostName[:], hostname)

	// Copy vendor string (up to 64 bytes, null-terminated)
	vendor := []byte("sing-box")
	copy(message.VendorString[:], vendor)

	// Send the message
	return binary.Write(c.controlConn, binary.LittleEndian, message)
}

func (c *PPTPClient) receiveStartControlConnectionReply() error {
	// Create a structure to hold the reply
	var reply struct {
		Header          PPTPHeader
		ProtocolVersion uint16
		Result          uint8
		ErrorCode       uint8
		FramingCap      uint32
		BearerCap       uint32
		MaxChannels     uint16
		FirmwareRev     uint16
		HostName        [64]byte
		VendorString    [64]byte
	}

	// Read the reply from the server
	if err := binary.Read(c.controlConn, binary.LittleEndian, &reply); err != nil {
		return E.Cause(err, "read start control connection reply")
	}

	// Verify the message type
	if reply.Header.MessageType != 1 || reply.Header.ControlType != StartControlConnectionReply {
		return E.New("unexpected PPTP control message type")
	}

	// Check the result code
	if reply.Result != 1 { // 1 = success
		errorMessages := map[uint8]string{
			2: "general error",
			3: "command channel already exists",
			4: "not authorized",
			5: "unsupported protocol version",
		}
		errorMsg, ok := errorMessages[reply.Result]
		if !ok {
			errorMsg = "unknown error"
		}
		return E.New("PPTP connection failed: ", errorMsg)
	}

	c.logger.Debug("PPTP control connection established")
	return nil
}

func (c *PPTPClient) sendOutgoingCallRequest() error {
	// Create the PPTP Outgoing-Call-Request message
	var message struct {
		Header          PPTPHeader
		CallID          uint16   // Call ID assigned by client
		SerialNumber    uint16   // Serial number for this request
		MinBPS          uint32   // Minimum BPS (bits per second)
		MaxBPS          uint32   // Maximum BPS
		BearerType      uint32   // Bearer type
		FramingType     uint32   // Framing type
		RecvWindowSize  uint16   // Receive window size
		ProcessingDelay uint16   // Processing delay
		PhoneNumberLen  uint16   // Length of phone number
		Reserved1       uint16   // Reserved
		PhoneNumber     [64]byte // Phone number
		Subaddress      [64]byte // Subaddress
	}

	// Fill header
	message.Header.Length = 168             // Total message length
	message.Header.MessageType = 1          // Control message
	message.Header.MagicCookie = 0x1A2B3C4D // Magic cookie
	message.Header.ControlType = OutgoingCallRequest

	// Fill message fields
	message.CallID = c.callID
	message.SerialNumber = 1
	message.MinBPS = 0          // No minimum BPS requirement
	message.MaxBPS = 0          // No maximum BPS requirement
	message.BearerType = 3      // Any bearer type
	message.FramingType = 3     // Any framing type
	message.RecvWindowSize = 64 // Typical value
	message.ProcessingDelay = 0 // No delay
	message.PhoneNumberLen = 0  // No phone number for VPN

	// Send the message
	return binary.Write(c.controlConn, binary.LittleEndian, message)
}

func (c *PPTPClient) receiveOutgoingCallReply() error {
	// Create a structure to hold the reply
	var reply struct {
		Header          PPTPHeader
		CallID          uint16 // Call ID assigned by client
		PeerCallID      uint16 // Call ID assigned by server
		Result          uint8  // Result code
		ErrorCode       uint8  // Error code
		CauseCode       uint16 // Cause code
		ConnectSpeed    uint32 // Connect speed
		RecvWindowSize  uint16 // Receive window size
		ProcessingDelay uint16 // Processing delay
		PhysicalChannel uint32 // Physical channel ID
	}

	// Read the reply from the server
	if err := binary.Read(c.controlConn, binary.LittleEndian, &reply); err != nil {
		return E.Cause(err, "read outgoing call reply")
	}

	// Verify the message type
	if reply.Header.MessageType != 1 || reply.Header.ControlType != OutgoingCallReply {
		return E.New("unexpected PPTP control message type")
	}

	// Check the result code
	if reply.Result != 1 { // 1 = success
		errorMessages := map[uint8]string{
			2: "general error",
			3: "no carrier",
			4: "busy signal",
			5: "no dial tone",
			6: "timeout",
			7: "not available",
		}
		errorMsg, ok := errorMessages[reply.Result]
		if !ok {
			errorMsg = "unknown error"
		}
		return E.New("PPTP outgoing call failed: ", errorMsg)
	}

	// Store the peer's Call ID for future reference
	c.peerCallID = reply.PeerCallID

	c.logger.Debug("PPTP tunnel established")
	return nil
}

func (c *PPTPClient) handleControlMessages() {
	defer c.Close()

	// Buffer to read the header
	headerBuf := make([]byte, 8)

	for {
		// Read just the header first to determine message length
		_, err := io.ReadFull(c.controlConn, headerBuf)
		if err != nil {
			c.logger.Error("PPTP control connection error:", err)
			return
		}

		var header PPTPHeader
		if err := binary.Read(bytes.NewReader(headerBuf), binary.LittleEndian, &header); err != nil {
			c.logger.Error("Failed to parse PPTP header:", err)
			return
		}

		// Calculate remaining bytes to read
		remaining := int(header.Length) - 8
		if remaining < 0 {
			c.logger.Error("Invalid PPTP message length")
			return
		}

		// Read the rest of the message
		messageBuf := make([]byte, remaining)
		if remaining > 0 {
			if _, err := io.ReadFull(c.controlConn, messageBuf); err != nil {
				c.logger.Error("Failed to read PPTP message body:", err)
				return
			}
		}

		// Handle different message types
		switch header.ControlType {
		case 6: // Echo-Request
			c.handleEchoRequest(header, messageBuf)
		case 13: // Call-Disconnect-Notify
			c.logger.Info("Received Call-Disconnect-Notify from server")
			return
		case 14: // WAN-Error-Notify
			c.logger.Error("Received WAN-Error-Notify from server")
			// Could parse error details from messageBuf
		default:
			c.logger.Debug("Received unhandled PPTP control message type:", header.ControlType)
		}
	}
}

func (c *PPTPClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.connected {
		return nil
	}

	// Send Call-Disconnect-Notify if needed

	if c.controlConn != nil {
		c.controlConn.Close()
		c.controlConn = nil
	}

	c.connected = false
	return nil
}

func (c *PPTPClient) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	// Ensure control connection is established
	if err := c.Connect(ctx); err != nil {
		return nil, err
	}

	// Create a packet connection through the PPTP tunnel
	// This is a placeholder - you'll need to implement the actual UDP over GRE tunneling
	return nil, E.New("UDP over PPTP not implemented yet")
}

// This method will create a tunnel for a specific connection
func (c *PPTPClient) DialTunnel(ctx context.Context, destination M.Socksaddr) (net.Conn, error) {
	// Ensure control connection is established
	if err := c.Connect(ctx); err != nil {
		return nil, err
	}

	// TODO: Implement GRE tunnel setup and return a connection that uses it

	// This is a placeholder - you'll need to implement the actual GRE tunneling
	return nil, E.New("GRE tunneling not implemented yet")
}

func (c *PPTPClient) handleEchoRequest(header PPTPHeader, messageBody []byte) {
	// Extract the identifier from the Echo-Request
	if len(messageBody) < 4 || header.Length < 1 {
		c.logger.Error("Echo-Request message too short")
		return
	}

	identifier := binary.LittleEndian.Uint32(messageBody[:4])

	// Create and send Echo-Reply
	var reply struct {
		Header     PPTPHeader
		Identifier uint32
		Result     uint8
		Reserved   [3]byte
	}

	reply.Header.Length = 16
	reply.Header.MessageType = 1
	reply.Header.MagicCookie = 0x1A2B3C4D
	reply.Header.ControlType = 7 // Echo-Reply
	reply.Identifier = identifier
	reply.Result = 1 // Success

	if err := binary.Write(c.controlConn, binary.LittleEndian, reply); err != nil {
		c.logger.Error("Failed to send Echo-Reply:", err)
	}
}
