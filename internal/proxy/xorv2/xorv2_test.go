package xorv2

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

// TestToPeer tests the ToPeer method of the xorv2 proxy. It specifically
// ensures that the index of the XOR key offset is maintained across multiple
// calls to the method.
func TestToPeer(t *testing.T) {
	xor, err := New(net.ParseIP("2001:300:7f00:1:d204:646f:6f6f:6f6d"))
	if err != nil {
		t.Errorf("unable to create new XOR proxy: %v", err)
	}

	pr, pw := io.Pipe()
	var dst bytes.Buffer

	go func() {
		pw.Write([]byte("foo bar baz"))
		time.Sleep(10 * time.Millisecond)
		pw.Write([]byte("foo bar baz"))
		pw.Close()
	}()

	xor.ToPeer(&dst, pr)
	data := dst.Bytes()

	if bytes.Equal(data[:11], data[11:]) {
		t.Errorf("XOR:ed data is equal")
	}
}
