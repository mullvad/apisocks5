package xorv2

import (
	"bytes"
	"net"
	"strings"
	"testing"
)

// TestToPeer tests the ToPeer method of the xorv2 proxy. It specifically
// ensures that the index of the XOR key offset is maintained across multiple
// calls to the method.
func TestToPeer(t *testing.T) {
	xor, err := New(net.ParseIP("2001:300:7f00:1:d204:646f:6f6f:6f6d"))
	if err != nil {
		t.Errorf("unable to create new XOR proxy: %v", err)
	}

	var dst1 bytes.Buffer
	xor.ToPeer(&dst1, strings.NewReader("foo bar baz"))

	var dst2 bytes.Buffer
	xor.ToPeer(&dst2, strings.NewReader("foo bar baz"))

	if bytes.Equal(dst1.Bytes(), dst2.Bytes()) {
		t.Errorf("XOR:ed data is equal")
	}
}
