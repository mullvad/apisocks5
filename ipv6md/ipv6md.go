package ipv6md

import (
	"encoding/binary"
	"errors"
	"net"
)

var (
	ErrUnknownData    = errors.New("unknown data")
	ErrUnknownType    = errors.New("unknown type")
	ErrUnexpectedType = errors.New("unexpected type")
)

type Type uint16

const (
	Unknown       Type = 0x00
	AddrPort      Type = 0x01
	AddrPortXOR   Type = 0x02
	AddrPortXORV2 Type = 0x03
)

func (t Type) ToUint16() uint16 {
	return uint16(t)
}

func (t Type) String() string {
	switch t {
	case AddrPort:
		return "AddrPort"
	case AddrPortXOR:
		return "AddrPortXOR"
	case AddrPortXORV2:
		return "AddrPortXORV2"
	default:
		return "Unknown"
	}
}

// ipv6Prefix contains the first two bytes of the IPv6 address, set to 0x20 and
// 0x01 to better masquerade the address as a "real" address.
var IPv6Prefix = []byte{0x20, 0x01}

// GetType returns the type of the given IP address or an error if the data
// contains an unknown type.
func GetType(ip net.IP) (Type, error) {
	data := ip.To16()
	if data == nil {
		return Unknown, ErrUnknownData
	}

	switch Type(binary.LittleEndian.Uint16(data[2:4])) {
	case AddrPort:
		return AddrPort, nil
	case AddrPortXOR:
		return AddrPortXOR, nil
	case AddrPortXORV2:
		return AddrPortXORV2, nil
	}

	return Unknown, ErrUnknownType
}
