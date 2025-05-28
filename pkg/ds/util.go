package ds

import (
	"fmt"
	"net"
)

// Convert an IP to a binary string key
func ConvertAddressToKey(ip net.IP) string {
	ip = ip.To16()
	b := ""
	for _, v := range ip {
		b += fmt.Sprintf("%08b", v)
	}
	return b
}

// Convert a network (IP + mask) to binary key up to prefix length
func ConvertNetworkToKey(network net.IPNet) string {
	// normalzie the prefix
	prefix := network.IP.To16()

	// convert the prefix to key
	prefixKey := ConvertAddressToKey(prefix)

	prefixLength, _ := network.Mask.Size()
	offset := 0

	if IsMappedToIPv6(prefix) {
		offset += 96
	}

	return prefixKey[offset : prefixLength+offset]
}

// Checks if the IP is a pure IPv6 or an IPv4-mapped IPv6
func IsMappedToIPv6(ip net.IP) bool {
	ip = ip.To16()

	// Check for IPv4-mapped IPv6 (::ffff:x.x.x.x)
	if ip.To4() != nil && ip[10] == 0xff && ip[11] == 0xff {
		return true
	}
	return false
}

// Convert net.IP and prefix length to *net.IPNet
func IPToNetwork(ip net.IP, prefixLength int) net.IPNet {
	var bits int
	var baseIP net.IP

	switch {
	// case IsMappedToIPv6(ip):
	// 	// If it's IPv4-mapped IPv6 (::ffff:a.b.c.d)
	// 	baseIP = ip.To4()
	// 	bits = 32
	case ip.To4() != nil:
		// Regular IPv4
		baseIP = ip.To4()
		bits = 32
	case ip.To16() != nil:
		// Pure IPv6
		baseIP = ip.To16()
		bits = 128
	default:
		return net.IPNet{}
	}

	mask := net.CIDRMask(prefixLength, bits)
	network := net.IPNet{
		IP:   baseIP.Mask(mask),
		Mask: mask,
	}

	return network
}
