package addrportxor

import (
	"bytes"
	"testing"
)

type addrPortTest struct {
	input    string
	xorBytes uint16
	xorKey   []byte
	encoded  string
}

var addrPortXORTests = []addrPortTest{
	{
		input:    "127.0.0.1:1337",
		xorBytes: 10,
		xorKey:   []byte{0x01, 0x02, 0x03, 0x04},
		encoded:  "2001:200:7f00:1:3905:a00:102:304",
	},
	{
		input:    "192.168.1.1:443",
		xorBytes: 6667,
		xorKey:   []byte{0x04, 0x02},
		encoded:  "2001:200:c0a8:101:bb01:b1a:402:0",
	},
}

func TestAddrPortXOR(t *testing.T) {
	for _, test := range addrPortXORTests {
		t.Run(test.input, func(t *testing.T) {
			encoded, _ := Encode(test.input, test.xorBytes, test.xorKey)
			if test.encoded != encoded.String() {
				t.Errorf("unexpected encoding result")
				t.Logf("actual: %#v", encoded.String())
				t.Logf("expected: %#v", test.encoded)
			}

			decoded, _ := Decode(encoded)
			if test.input != decoded.AddrPort.String() {
				t.Errorf("unexpected decoding result")
				t.Logf("actual: %#v", decoded.AddrPort.String())
				t.Logf("expected: %#v", test.input)
			}
			if test.xorBytes != decoded.XORBytes {
				t.Errorf("unexpected xor bytes")
				t.Logf("actual: %#v", decoded.XORBytes)
				t.Logf("expected: %#v", test.xorBytes)
			}
			if !bytes.Equal(test.xorKey, decoded.XORKey) {
				t.Errorf("unexpected xor key")
				t.Logf("actual: %#v", decoded.XORKey)
				t.Logf("expected: %#v", test.xorKey)
			}
		})
	}
}
