package proxy

import (
	"io"

	"github.com/mullvad/proxy/typ"
)

type Proxy interface {
	Address() string
	FromPeer(dst io.Writer, src io.Reader)
	ToPeer(dst io.Writer, src io.Reader)
	Type() typ.Type
}
