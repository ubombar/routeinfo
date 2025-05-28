package ds

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/armon/go-radix"
)

const DefaultEntrySize = 1

// This is the main implementation of the forwarding info as stated in the design document.
// This struct uses a redix tree to perform lookups.
//
// With this we only have one radix tree containing everyting.
type FT struct {
	tree                *radix.Tree
	optimizeForIPv4     bool
	defaultPrefixLength uint
}

// Creates a new forwarding table.
func NewFowardingTable(optimizeForIPv4 bool, defaultPrefixLength uint) *FT {
	return &FT{
		tree:                radix.New(),
		optimizeForIPv4:     optimizeForIPv4,
		defaultPrefixLength: defaultPrefixLength,
	}
}

// Performs a longest prefix match to get the next hop of a given ip address. This
// speciality is used in the routers.
func (f *FT) Lookup(address net.IP) (*FTEntry, bool, error) {
	// The IP address is onverted into a IPv6 String.
	key := IPToKey(address)
	_, item, found := f.tree.LongestPrefix(key)
	if !found {
		return nil, false, nil
	}

	if setObj, ok := item.(*FTEntry); !ok {
		return nil, false, errors.New("the expected type of NFMap is not an NFMap")
	} else {
		return setObj, true, nil
	}
}

// Check if the given network is already registered.
func (f *FT) Contains(network net.IPNet) (*FTEntry, bool, error) {
	// The IP address is onverted into a IPv6 String.
	key := NetworkToKey(network)
	item, found := f.tree.Get(key)
	if !found {
		return nil, false, nil
	}

	if setObj, ok := item.(*FTEntry); !ok {
		return nil, false, errors.New("the expected type of NFMap is not an NFMap")
	} else {
		return setObj, true, nil
	}
}

// Inserts the nexthop address to the reverse forwarding table.
func (f *FT) Insert(network net.IPNet, nexthop net.IP) error {
	entry, found, err := f.Contains(network)
	if err != nil {
		return err
	}

	if !found {
		entry = newFTEntry(DefaultEntrySize)
	}

	entry.Add(nexthop)

	key := NetworkToKey(network)
	f.tree.Insert(key, entry)

	return nil
}

// Converts the forwarding table into a String
func (f *FT) String() string {
	var sb strings.Builder

	for networkKey, entry := range f.tree.ToMap() {
		networkPrefix, err := KeyToIP(networkKey)
		if err != nil {
			panic(err)
		}
		network := IPToNetwork(networkPrefix, int(f.defaultPrefixLength))
		sb.WriteString(fmt.Sprintf("\t%v -> %v\n", network.String(), entry))
	}

	return sb.String()
}

// FTEntry is an object that denotes the next hops of a router. It is a set because
// there can be ECPM and load balancing enabled.
type FTEntry struct {
	// IPs denote an array of net.IP addresses. Thanks to them we can perform
	// some dset operations like add, contains etc.
	dset []net.IP
}

// Creates a new NHSet struct.
func newFTEntry(size uint) *FTEntry {
	return &FTEntry{
		dset: make([]net.IP, 0, size),
	}
}

// Adds the ip address if it is not already in the set.
func (n *FTEntry) Add(ip net.IP) {
	if !n.Contains(ip) {
		n.dset = append(n.dset, ip)
	}
}

// Checks if the given IP address is already in the set.
func (n *FTEntry) Contains(ip net.IP) bool {
	for _, existing := range n.dset {
		if existing.Equal(ip) {
			return true
		}
	}
	return false
}

// ToString method returns a string representation of all IPs
func (ft *FTEntry) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, ip := range ft.dset {
		sb.WriteString(ip.String())
		if i < len(ft.dset)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}
