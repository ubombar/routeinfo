package ds

import (
	"errors"
	"net"

	"github.com/armon/go-radix"
)

const DefaultFtentrySize = 1

// This is the main implementation of the forwarding info as stated in the design document.
// This struct uses a redix tree to perform lookups. Note that it has only one redix tree
// and for each entry there is a mapper called NFMap. That mapper maps the given near address
// to a set of far addresses.
//
// With this we only have one radix tree containing everyting.
type FT struct {
	tree *radix.Tree
}

// Creates a new forwarding table.
func NewFowardingTable() *FT {
	return &FT{
		tree: radix.New(),
	}
}

// Performs a longest prefix lookup and returns the *NFMap and *NHSet objects. Returns
// false and nil for all two of the objects.
func (f *FT) Lookup(address net.IP) (*FTEntry, bool, error) {
	// The IP address is onverted into a IPv6 String.
	key := ConvertAddressToKey(address)
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

// Inserts the nexthop address to the reverse forwarding table.
func (f *FT) Insert(network net.IPNet, nexthop net.IP) error {
	entry, found, err := f.Lookup(nexthop)
	if err != nil {
		return err
	}

	if !found {
		entry = newFTEntry(DefaultFtentrySize)
	}

	entry.Add(nexthop)

	key := ConvertNetworkToKey(network)
	f.tree.Insert(key, entry)
	return nil
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
		dset: make([]net.IP, size),
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
