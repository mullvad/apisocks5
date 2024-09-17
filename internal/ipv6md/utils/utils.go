package utils

import "net/netip"

func ToNetIPAddr(input []byte) netip.Addr {
	var addrBytes [4]byte
	copy(addrBytes[:], input)
	return netip.AddrFrom4(addrBytes)

}
