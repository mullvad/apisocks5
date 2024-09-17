package typ

type Type int

const (
	Plain Type = 0x01
	XOR   Type = 0x02
	XORV2 Type = 0x03
)

func (t Type) String() string {
	switch t {
	case Plain:
		return "plain"
	case XOR:
		return "xor"
	case XORV2:
		return "xor-v2"
	default:
		return "Unknown"
	}
}
