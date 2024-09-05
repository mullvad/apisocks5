package addrportxorv2

import (
	"bytes"
	"testing"
)

type addrPortTest struct {
	input   string
	xorKey  []byte
	encoded string
}

var addrPortXORV2Tests = []addrPortTest{
	{
		input:   "127.0.0.1:27500",
		xorKey:  []byte{0x01, 0x02, 0x03, 0x04, 0x5, 0x6},
		encoded: "2001:300:7f00:1:6c6b:102:304:506",
	},
	{
		input:   "192.168.1.1:27500",
		xorKey:  []byte{0x04, 0x02},
		encoded: "2001:300:c0a8:101:6c6b:402::",
	},
}

func TestAddrPortXORV2(t *testing.T) {
	for _, test := range addrPortXORV2Tests {
		t.Run(test.input, func(t *testing.T) {
			encoded, _ := Encode(test.input, test.xorKey)
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
			if !bytes.Equal(test.xorKey, decoded.XORKey) {
				t.Errorf("unexpected xor key")
				t.Logf("actual: %#v", decoded.XORKey)
				t.Logf("expected: %#v", test.xorKey)
			}
		})
	}
}
