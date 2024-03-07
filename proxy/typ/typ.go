package typ

type Type int

const (
	Plain Type = 0x01
	XOR   Type = 0x02
)

func (t Type) String() string {
	switch t {
	case Plain:
		return "plain"
	case XOR:
		return "xor"
	default:
		return "Unknown"
	}
}
