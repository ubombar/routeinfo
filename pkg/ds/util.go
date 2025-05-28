package ds

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strings"
)

// Convert the IP into a binary string. It automatically maps it into IPv6
// if it is a IPv4.
func ConvertAddressToKey(ip net.IP) string {
	ip = ip.To16()
	b := ""
	for _, v := range ip {
		b += fmt.Sprintf("%08b", v)
	}
	return b
}

// Convert the the given network into
func ConvertNetworkToKey(network net.IPNet) string {
	// normalzie the prefix
	prefix := network.IP.To16()

	// convert the prefix to key
	prefixKey := ConvertAddressToKey(prefix)

	prefixLength, _ := network.Mask.Size()
	extra := 0

	if IsMappedToIPv6(prefix) {
		extra += 96
	}

	return prefixKey[:prefixLength+extra]
}

// Checks if the IP is a pure IPv6 or an IPv4-mapped IPv6
func IsMappedToIPv6(ip net.IP) bool {
	ip = ip.To16()

	// Essentially means it follows the format ::ffff:X.X.X.X
	for i := 0; i < 10; i++ {
		if ip[i] != 0x00 {
			return false
		}
	}
	for i := 10; i < 12; i++ {
		if ip[i] != 0xff {
			return false
		}
	}
	return true
}

// Convert IP to IPNet. Note that if the given ip is a IPv4 format, then
// 96 is added to the prefix length.
func IPToNetwork(ip net.IP, prefixLength int) net.IPNet {
	ip = ip.To16()

	if IsMappedToIPv6(ip) {
		mask := net.CIDRMask(prefixLength+96, 128)
		return net.IPNet{
			IP:   ip.Mask(mask),
			Mask: mask,
		}
	} else {
		mask := net.CIDRMask(prefixLength, 128)
		return net.IPNet{
			IP:   ip.Mask(mask),
			Mask: mask,
		}
	}
}

// Adds zeros to the end to fix it to 128 bits.
func AddPaddingToKey(key string) (string, error) {
	if len(key) > 128 {
		return "", errors.New("given key is larger than 128 characters")
	}
	postfix := strings.Repeat("0", 128-len(key))

	return fmt.Sprintf("%s%s", key, postfix), nil
}

// Convert binary string (up to 128 bits) to IPv6 net.IP
func KeyToIP(key string, addIPv6MappingPrefix bool) (net.IP, error) {
	// Pad the key to get the 128 bit representation.
	key, err := AddPaddingToKey(key, addIPv6MappingPrefix)
	if err != nil {
		return net.IP{}, err
	}
	// Convert binary string to big.int
	// HELLO MY NAME IS MECOBAYN!
	n := new(big.Int)
	n, ok := n.SetString(key, 2)
	if !ok {
		return nil, fmt.Errorf("invalid binary string")
	}

	// Convert to hex string
	hexStr := n.Text(16)
	if len(hexStr) < 32 {
		hexStr = fmt.Sprintf("%032s", hexStr) // pad to 32 hex digits (128 bits)
	}

	// Convert hex string to bytes
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	return net.IP(bytes).To16(), nil
}
