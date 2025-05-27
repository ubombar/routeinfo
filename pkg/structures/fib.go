package structures

import (
	"fmt"
	"net"
)

type FIB struct {
	// this is the root of the FIB.
	root *node
}

func NewFIB() *FIB {
	root := &node{}
	return &FIB{
		root: root,
	}
}

func (f *FIB) LookupString(nearAddressString, farAddressString string) (*nextHopStruct, bool) {
	nearAddress := net.IPAddr{IP: net.ParseIP(nearAddressString)}
	farAddress := net.IPAddr{IP: net.ParseIP(farAddressString)}
	return f.Lookup(&nearAddress, &farAddress)
}

func (f *FIB) Lookup(nearAddress, farAddress *net.IPAddr) (*nextHopStruct, bool) {
	if mapp, found := f.root.Search(farAddress.IP); !found {
		return nil, false
	} else {
		return mapp.Lookup(nearAddress.IP)
	}
}

func (f *FIB) InsertString(nearAddressString, farAddressString string, defaultPrefixLength int) bool {
	nearAddress := net.IPAddr{IP: net.ParseIP(nearAddressString)}
	farAddress := net.IPAddr{IP: net.ParseIP(farAddressString)}
	return f.Insert(&nearAddress, &farAddress, defaultPrefixLength)
}

func (f *FIB) Insert(nearAddress *net.IPAddr, farAddress *net.IPAddr, defaultPrefixLength int) bool {
	d, err := ipAddrToIPNet(*farAddress, defaultPrefixLength)
	if err != nil { // lol
		panic(err)
	}

	// Do a search on far address.
	mapp, found := f.root.Search(farAddress.IP)
	if !found {
		mapp = newNearFarsMap()
		f.root.Insert(d, mapp)
	}

	hopp, found := mapp.Lookup(nearAddress.IP)
	if !found {
		hopp = newNextHopStruct()
		mapp.Add(nearAddress.IP, hopp)
	}

	found = hopp.NextHopSet.Contains(farAddress)
	if !found {
		hopp.NextHopSet.Add(farAddress)
		return true
	} else {
		return false // if it exists, we don't add it again.
	}
}

type nextHopStruct struct {
	NextHopSet Set[*net.IPAddr]
}

func newNextHopStruct() *nextHopStruct {
	return &nextHopStruct{}
}

// nearFarsMap is a wrapper around a map from net.IP to X.
type nearFarsMap struct {
	data map[string]*nextHopStruct // string keys to handle net.IP equality and map lookups
}

// newNearFarsMap creates a new NearFarsMap.
func newNearFarsMap() *nearFarsMap {
	return &nearFarsMap{
		data: make(map[string]*nextHopStruct),
	}
}

// Add inserts an IP address with its associated X struct.
func (nfm *nearFarsMap) Add(ip net.IP, value *nextHopStruct) {
	nfm.data[ip.String()] = value
}

// Lookup returns the X associated with the given IP, or false if not found.
func (nfm *nearFarsMap) Lookup(ip net.IP) (*nextHopStruct, bool) {
	value, ok := nfm.data[ip.String()]
	return value, ok
}

// node represents a node in the radix tree optimized for net.IP.
type node struct {
	prefix   *net.IPNet
	value    *nearFarsMap
	children [2]*node
	isLeaf   bool
}

// Insert inserts a prefix with associated value.
func (n *node) Insert(ipNet *net.IPNet, value *nearFarsMap) {
	if n.prefix == nil {
		n.prefix = ipNet
		n.value = value
		n.isLeaf = true
		return
	}

	commonPrefixLen := commonPrefixLength(n.prefix, ipNet)

	if commonPrefixLen < maskLength(n.prefix) {
		// Split current node
		child := &node{
			prefix:   n.prefix,
			value:    n.value,
			children: n.children,
			isLeaf:   n.isLeaf,
		}
		n.prefix = &net.IPNet{
			IP:   n.prefix.IP.Mask(net.CIDRMask(commonPrefixLen, len(n.prefix.IP)*8)),
			Mask: net.CIDRMask(commonPrefixLen, len(n.prefix.IP)*8),
		}
		n.value = nil
		n.children = [2]*node{}
		n.isLeaf = false
		bit := getBit(child.prefix.IP, commonPrefixLen)
		n.children[bit] = child
	}

	if commonPrefixLen == maskLength(ipNet) {
		// Exact match or more specific
		n.value = value
		n.isLeaf = true
		return
	}

	bit := getBit(ipNet.IP, commonPrefixLen)
	if n.children[bit] == nil {
		n.children[bit] = &node{}
	}
	n.children[bit].Insert(ipNet, value)
}

// Search finds the longest prefix match for the given IP.
func (n *node) Search(ip net.IP) (*nearFarsMap, bool) {
	var result *nearFarsMap
	var found bool

	cur := n
	for cur != nil {
		if cur.prefix != nil && cur.prefix.Contains(ip) && cur.isLeaf {
			result = cur.value
			found = true
		}
		bit := getBit(ip, maskLength(cur.prefix))
		cur = cur.children[bit]
	}

	return result, found
}

// Helper functions

func getBit(ip net.IP, pos int) int {
	byteIndex := pos / 8
	if byteIndex >= len(ip) {
		return 0
	}
	bitIndex := 7 - (pos % 8)
	if (ip[byteIndex] & (1 << bitIndex)) != 0 {
		return 1
	}
	return 0
}

func commonPrefixLength(a, b *net.IPNet) int {
	maxLen := min(maskLength(a), maskLength(b))
	for i := 0; i < maxLen; i++ {
		if getBit(a.IP, i) != getBit(b.IP, i) {
			return i
		}
	}
	return maxLen
}

func maskLength(ipNet *net.IPNet) int {
	ones, _ := ipNet.Mask.Size()
	return ones
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Converts a net.IPAddr and prefix length to a net.IPNet
func ipAddrToIPNet(ipAddr net.IPAddr, prefixLen int) (*net.IPNet, error) {
	var maxBits int
	if ip4 := ipAddr.IP.To4(); ip4 != nil {
		maxBits = 32
	} else if ip6 := ipAddr.IP.To16(); ip6 != nil {
		maxBits = 128
	} else {
		return nil, fmt.Errorf("invalid IP address: %v", ipAddr.IP)
	}

	// Create a mask
	mask := net.CIDRMask(prefixLen, maxBits)

	// Apply the mask to get the network prefix
	networkIP := ipAddr.IP.Mask(mask)

	// Return net.IPNet
	return &net.IPNet{
		IP:   networkIP,
		Mask: mask,
	}, nil
}
