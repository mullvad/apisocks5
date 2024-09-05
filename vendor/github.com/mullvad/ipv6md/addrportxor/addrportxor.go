// Encapsulate an address and port together with a basic XOR definition
//
// Format description:
//
// Offsets      Description
// -------      -----------
// 00 - 02      Dummy header to make it look like a real IPv6 address
// 02 - 04      Type (little endian)
// 04 - 08      IPv4 address (big endian)
// 08 - 10      Port number (little endian)
// 10 - 12      Number of bytes to XOR; 0x00 means XOR everything
// 12 - 16      Key, 0x00 bytes will be skipped

package addrportxor

import (
	"encoding/binary"
	"errors"
	"net"
	"net/netip"

	"github.com/mullvad/ipv6md"
	"github.com/mullvad/ipv6md/addrport"
	"github.com/mullvad/ipv6md/utils"
)

var (
	ErrInvalidKeyLength = errors.New("invalid key length")
	ErrInvalidKey       = errors.New("invalid key")
)

// DecodedAddrPortXOR is used by the Decode function to encapsulate the
// returned values.
type DecodedAddrPortXOR struct {
	AddrPort netip.AddrPort
	XORBytes uint16
	XORKey   []byte
}

// Encode encodes the given address, port and XOR encryption details in an IPv6
// formatted slice of bytes.
func Encode(addrPort string, xorBytes uint16, xorKey []byte) (net.IP, error) {
	if len(xorKey) == 0 || len(xorKey) > 4 {
		return nil, ErrInvalidKeyLength
	}

	data, err := addrport.Encode(addrPort)
	if err != nil {
		return nil, err
	}

	binary.LittleEndian.PutUint16(data[2:4], ipv6md.AddrPortXOR.ToUint16())
	binary.LittleEndian.PutUint16(data[10:12], xorBytes)
	copy(data[12:16], xorKey)

	return net.IP(data[:]), nil
}

// Decode assumes an IPv4 address and port, along with a 16-bit integer and a
// 1-4 bytes long key, has been encoded within the IPv6 address. It returns a
// netip.AddrPort with the information and the XOR bytes and key.
func Decode(ip net.IP) (*DecodedAddrPortXOR, error) {
	if ip == nil {
		return nil, addrport.ErrAddrPortInvalidIP
	}

	data := []byte(ip.To16())
	if len(data) != 16 {
		return nil, addrport.ErrAddrPortInvalidIPLen
	}

	typ := binary.LittleEndian.Uint16(data[2:4])
	if ipv6md.Type(typ) != ipv6md.AddrPortXOR {
		return nil, ipv6md.ErrUnexpectedType
	}

	addr := utils.ToNetIPAddr(data[4:8])
	port := binary.LittleEndian.Uint16(data[8:10])
	ap := netip.AddrPortFrom(addr, port)

	xorBytes := binary.LittleEndian.Uint16(data[10:12])

	var xorKey []byte
	for _, b := range data[12:16] {
		if b == 0x00 {
			continue
		}

		xorKey = append(xorKey, b)
	}
	if len(xorKey) == 0 {
		return nil, ErrInvalidKey
	}

	return &DecodedAddrPortXOR{
		AddrPort: ap,
		XORBytes: xorBytes,
		XORKey:   xorKey,
	}, nil
}
