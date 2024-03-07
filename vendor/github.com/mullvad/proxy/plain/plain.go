package plain

import (
	"io"
	"net"

	"github.com/mullvad/ipv6md/addrport"
	"github.com/mullvad/proxy/typ"
)

type plain struct {
	addrPort string
}

func New(ip net.IP) (*plain, error) {
	target, err := addrport.Decode(ip)
	if err != nil {
		return nil, err
	}

	return &plain{
		addrPort: target.String(),
	}, nil
}

func (p *plain) Address() string                       { return p.addrPort }
func (p *plain) Type() typ.Type                        { return typ.Plain }
func (p *plain) FromPeer(dst io.Writer, src io.Reader) { p.forward(dst, src) }
func (p *plain) ToPeer(dst io.Writer, src io.Reader)   { p.forward(dst, src) }

func (p *plain) forward(dst io.Writer, src io.Reader) {
	io.Copy(dst, src)

	tcpConn, ok := dst.(*net.TCPConn)
	if !ok {
		return
	}

	tcpConn.CloseWrite()
}
