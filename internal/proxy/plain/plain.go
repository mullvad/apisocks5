package plain

import (
	"io"
	"net"

	"github.com/mullvad/apisocks5/internal/ipv6md/addrport"
	"github.com/mullvad/apisocks5/internal/proxy/typ"
)

type Plain struct {
	addrPort string
}

func New(ip net.IP) (*Plain, error) {
	target, err := addrport.Decode(ip)
	if err != nil {
		return nil, err
	}

	return &Plain{
		addrPort: target.String(),
	}, nil
}

func (p *Plain) Address() string                       { return p.addrPort }
func (p *Plain) Type() typ.Type                        { return typ.Plain }
func (p *Plain) FromPeer(dst io.Writer, src io.Reader) { p.forward(dst, src) }
func (p *Plain) ToPeer(dst io.Writer, src io.Reader)   { p.forward(dst, src) }

func (p *Plain) forward(dst io.Writer, src io.Reader) {
	io.Copy(dst, src)

	tcpConn, ok := dst.(*net.TCPConn)
	if !ok {
		return
	}

	tcpConn.CloseWrite()
}
