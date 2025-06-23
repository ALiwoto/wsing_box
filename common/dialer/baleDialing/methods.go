package baleDialing

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/gofrs/uuid/v5"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/singingEncoding"
	"github.com/sagernet/sing-box/log"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/pipe"
)

//---------------------------------------------------------

func (d *BaleDialerContainer) DialContext(
	ctx context.Context,
	network string,
	destination M.Socksaddr,
) (net.Conn, error) {
	pipe1, pipe2 := pipe.Pipe()
	connIdProvider, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	connId := strings.ReplaceAll(connIdProvider.String(), "-", "")
	baleConnectionsPool.Add(connId, &pipe2)
	bConn := &BaleConn{
		pipe:         pipe1,
		writeLock:    &sync.Mutex{},
		readLock:     &sync.Mutex{},
		ConnectionId: connId,
		botsPool:     d.getInsideBots(),
		Address: &BaleFakeAddr{
			NetworkStr: destination.Network(),
			AddressStr: destination.String(),
		},
		isClosed: &atomic.Bool{},
	}

	go bConn.bufferFlushWorker()
	return bConn, nil
}

func (d *BaleDialerContainer) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	return nil, nil
}

func (d *BaleDialerContainer) getInsideBots() []*BaleBotContainer {
	result := []*BaleBotContainer{}

	for current := range d.Bots.Pairs {
		result = append(result, d.Bots.Pairs[current].InsideBot)
	}

	return result
}

//---------------------------------------------------------

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *BaleConn) Read(b []byte) (n int, err error) {
	c.readLock.Lock()
	defer c.readLock.Unlock()

	return c.pipe.Read(b)
}

func (c *BaleConn) bufferFlushWorker() {
	for !c.isClosed.Load() {
		time.Sleep(bufferFlushInterval)

		c.writeLock.Lock()
		if len(c.bufferedData) == 0 {
			c.writeLock.Unlock()
			continue
		}

		// pass zero min lenth, since we want to flush everything
		err := c.writeBufferedData(0)
		if err != nil {
			log.Error("Failed to flush data:", err)
		}
		c.writeLock.Unlock()
	}
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (c *BaleConn) Write(b []byte) (n int, err error) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	// F-<conn_id>-<pckt_id> [data here]
	// R-<conn_id> [COMMAND]
	c.bufferedData = append(c.bufferedData, b...)
	err = c.writeBufferedData(MinCharLen)
	return len(b), err
}

func (c *BaleConn) writeBufferedData(minLen int) error {
	allData := singingEncoding.StdEncoding.EncodeToString(c.bufferedData)
	totalLen := utf8.RuneCountInString(allData)
	// packetUUID, _ := uuid.NewV4()
	// packetId := strings.ReplaceAll(packetUUID.String(), "-", "")
	if totalLen < minLen {
		// too small...it's better we let it get buffered
		return nil
	}

	if totalLen < MaxCharLen {
		err := c.getRandomBot().SendData("F-" + c.ConnectionId + " " + allData)
		return err
	}

	// the data is larger than we thought...we have to send it chunk by chunk
	allChucks := MakeChunks(allData, totalLen, MaxCharLen)
	for i := range allChucks {
		chunk := allChucks[i]
		// err := c.getRandomBot().SendData("F-" + c.ConnectionId + "-" + packetId + " " + chunk)
		err := c.getRandomBot().SendData("F-" + c.ConnectionId + " " + chunk)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *BaleConn) getRandomBot() *BaleBotContainer {
	return c.botsPool[rand.Intn(len(c.botsPool))]
}

// // Close closes the connection.
// // Any blocked Read or Write operations will be unblocked and return errors.
func (c *BaleConn) Close() error {
	err := c.getRandomBot().SendData("R-" + c.ConnectionId + baleCommandCloseConn)
	if err != nil {
		return err
	}
	return c.pipe.Close()
}

// // LocalAddr returns the local network address, if known.
func (c *BaleConn) LocalAddr() net.Addr {
	return c.Address
}

// // RemoteAddr returns the remote network address, if known.
func (c *BaleConn) RemoteAddr() net.Addr {
	return c.Address
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c *BaleConn) SetDeadline(t time.Time) error {
	return c.pipe.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *BaleConn) SetReadDeadline(t time.Time) error {
	return c.pipe.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *BaleConn) SetWriteDeadline(t time.Time) error {
	return c.pipe.SetWriteDeadline(t)
}

//---------------------------------------------------------

// Network is name of the network (for example, "tcp", "udp")
func (b *BaleFakeAddr) Network() string {
	return b.NetworkStr
}

// String is string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
func (b *BaleFakeAddr) String() string {
	return b.AddressStr
}

//---------------------------------------------------------

// SendData sends certain data to a random channel, hoping the other pair
// receives it.
func (c *BaleBotContainer) SendData(data string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	triesCount := 0
	for c.Bot != nil {
		_, err := c.Bot.SendMessage(
			c.BaleConfig.Chats[rand.Intn(len(c.BaleConfig.Chats))],
			data,
			&gotgbot.SendMessageOpts{},
		)

		if err != nil {
			if triesCount > MaxDataSendRetries {
				// just give up
				return err
			}

			errStr := err.Error()
			switch {
			case strings.Contains(errStr, TooManyRequestStr):
				_, after, found := strings.Cut(errStr, TooManyRequestStr)
				if found && after != "" {
					time.Sleep(time.Duration(ssg.ToInt(after)) * time.Second)
					triesCount++
				}
			default:
				return err
			}

			// the error needs retrying
			continue
		}
		return nil
	}

	return errors.New("no bot instance to send data")
}

//---------------------------------------------------------
