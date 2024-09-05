package xorv2

import (
	"io"
	"net"

	"github.com/mullvad/ipv6md/addrportxorv2"
	"github.com/mullvad/proxy/typ"
)

type xor struct {
	addrPort     string
	xorKey       []byte
	xorKeyOffset int
}

func New(ip net.IP) (*xor, error) {
	target, err := addrportxorv2.Decode(ip)
	if err != nil {
		return nil, err
	}

	return &xor{
		addrPort: target.AddrPort.String(),
		xorKey:   target.XORKey,
	}, nil
}

func (x *xor) Address() string                       { return x.addrPort }
func (x *xor) Type() typ.Type                        { return typ.XORV2 }
func (x *xor) FromPeer(dst io.Writer, src io.Reader) { x.forward(dst, src) }
func (x *xor) ToPeer(dst io.Writer, src io.Reader)   { x.forward(dst, src) }

func (x *xor) forward(dst io.Writer, src io.Reader) {
	buf := make([]byte, 1024*64)

	for {
		nr, err := src.Read(buf)
		if err != nil || nr <= 0 {
			break
		}

		for i := 0; i < nr; i++ {
			buf[i] ^= x.xorKey[x.xorKeyOffset]
			x.xorKeyOffset = (x.xorKeyOffset + 1) % len(x.xorKey)
		}

		nw, err := dst.Write(buf[0:nr])
		if (err != nil) || (nr != nw) {
			break
		}
	}

	tcpConn, ok := dst.(*net.TCPConn)
	if !ok {
		return
	}

	tcpConn.CloseWrite()
}
