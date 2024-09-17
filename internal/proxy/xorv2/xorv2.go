package xorv2

import (
	"io"
	"net"

	"github.com/mullvad/apisocks5/internal/ipv6md/addrportxorv2"
	"github.com/mullvad/apisocks5/internal/proxy/typ"
)

type XORV2 struct {
	addrPort string
	xorKey   []byte
}

func New(ip net.IP) (*XORV2, error) {
	target, err := addrportxorv2.Decode(ip)
	if err != nil {
		return nil, err
	}

	return &XORV2{
		addrPort: target.AddrPort.String(),
		xorKey:   target.XORKey,
	}, nil
}

func (x *XORV2) Address() string                       { return x.addrPort }
func (x *XORV2) Type() typ.Type                        { return typ.XORV2 }
func (x *XORV2) FromPeer(dst io.Writer, src io.Reader) { x.forward(dst, src) }
func (x *XORV2) ToPeer(dst io.Writer, src io.Reader)   { x.forward(dst, src) }

func (x *XORV2) forward(dst io.Writer, src io.Reader) {
	buf := make([]byte, 1024*64)
	offset := 0

	for {
		nr, err := src.Read(buf)
		if err != nil || nr <= 0 {
			break
		}

		for i := 0; i < nr; i++ {
			buf[i] ^= x.xorKey[offset]
			offset = (offset + 1) % len(x.xorKey)
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
