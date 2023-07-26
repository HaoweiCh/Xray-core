package singbridge

import (
	"context"
	"os"

	"github.com/haoweich/xray-core/common/net"
	"github.com/haoweich/xray-core/common/net/cnc"
	"github.com/haoweich/xray-core/common/session"
	"github.com/haoweich/xray-core/proxy"
	"github.com/haoweich/xray-core/transport"
	"github.com/haoweich/xray-core/transport/internet"
	"github.com/haoweich/xray-core/transport/pipe"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

var _ N.Dialer = (*XrayDialer)(nil)

type XrayDialer struct {
	internet.Dialer
}

func NewDialer(dialer internet.Dialer) *XrayDialer {
	return &XrayDialer{dialer}
}

func (d *XrayDialer) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
	return d.Dialer.Dial(ctx, ToDestination(destination, ToNetwork(network)))
}

func (d *XrayDialer) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	return nil, os.ErrInvalid
}

type XrayOutboundDialer struct {
	outbound proxy.Outbound
	dialer   internet.Dialer
}

func NewOutboundDialer(outbound proxy.Outbound, dialer internet.Dialer) *XrayOutboundDialer {
	return &XrayOutboundDialer{outbound, dialer}
}

func (d *XrayOutboundDialer) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
	ctx = session.ContextWithOutbound(context.Background(), &session.Outbound{
		Target: ToDestination(destination, ToNetwork(network)),
	})
	opts := []pipe.Option{pipe.WithSizeLimit(64 * 1024)}
	uplinkReader, uplinkWriter := pipe.New(opts...)
	downlinkReader, downlinkWriter := pipe.New(opts...)
	conn := cnc.NewConnection(cnc.ConnectionInputMulti(downlinkWriter), cnc.ConnectionOutputMulti(uplinkReader))
	go d.outbound.Process(ctx, &transport.Link{Reader: downlinkReader, Writer: uplinkWriter}, d.dialer)
	return conn, nil
}

func (d *XrayOutboundDialer) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	return nil, os.ErrInvalid
}
