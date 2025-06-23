package pptp

import (
	"context"
	"net"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/adapter/outbound"
	"github.com/sagernet/sing-box/common/dialer"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

func RegisterOutbound(registry *outbound.Registry) {
	outbound.Register(registry, C.TypeShadowTLS, NewOutbound)
}

type Outbound struct {
	adapter.Outbound
	tag        string
	dialer     N.Dialer
	serverAddr M.Socksaddr
	username   string
	password   string
	logger     log.ContextLogger
	client     *PPTPClient
}

func NewOutbound(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.PPTPOutboundOptions) (adapter.Outbound, error) {
	outbound := &Outbound{
		tag:      tag,
		username: options.Username,
		password: options.Password,
		logger:   logger,
	}

	if options.Server == "" {
		return nil, E.New("missing server address")
	}

	outbound.serverAddr = M.ParseSocksaddr(options.Server)
	if options.ServerPort != 0 {
		outbound.serverAddr.Port = options.ServerPort
	} else {
		// Default PPTP port
		outbound.serverAddr.Port = 1723
	}

	myDialer, err := dialer.NewWithOptions(dialer.Options{
		Context:        ctx,
		Options:        options.DialerOptions,
		RemoteIsDomain: options.ServerIsDomain(),
	})
	if err != nil {
		return nil, err
	}
	outbound.dialer = myDialer

	return outbound, nil
}

func (o *Outbound) Tag() string {
	return o.tag
}

func (o *Outbound) Type() string {
	return "pptp"
}

func (o *Outbound) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
	// PPTP can tunnel both TCP and UDP traffic
	o.client = NewPPTPClient(o.serverAddr, o.username, o.password, o.logger)
	return o.client.DialTunnel(ctx, destination)
}

func (o *Outbound) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	o.client = NewPPTPClient(o.serverAddr, o.username, o.password, o.logger)
	packetConn, err := o.client.ListenPacket(ctx, destination)
	if err != nil {
		return nil, err
	}
	return packetConn, nil
}

func (o *Outbound) NewConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) error {
	return E.New("not implemented")
}

func (o *Outbound) NewPacketConnection(ctx context.Context, conn N.PacketConn, metadata adapter.InboundContext) error {
	return E.New("not implemented")
}

func (o *Outbound) Close() error {
	return nil
}
