package addrport

import (
	"testing"
)

type addrPortTest struct {
	input   string
	encoded string
}

var addrPortTests = []addrPortTest{
	{
		input:   "127.0.0.1:1337",
		encoded: "2001:100:7f00:1:3905::",
	},
	{
		input:   "192.168.1.1:443",
		encoded: "2001:100:c0a8:101:bb01::",
	},
}

func TestAddrPort(t *testing.T) {
	for _, test := range addrPortTests {
		t.Run(test.input, func(t *testing.T) {
			encoded, _ := Encode(test.input)
			if test.encoded != encoded.String() {
				t.Errorf("unexpected encoding result")
				t.Logf("actual: %#v", encoded.String())
				t.Logf("expected: %#v", test.encoded)
			}

			decoded, _ := Decode(encoded)
			if test.input != decoded.String() {
				t.Errorf("unexpected decoding result")
				t.Logf("actual: %#v", decoded.String())
				t.Logf("expected: %#v", test.input)
			}
		})
	}
}
