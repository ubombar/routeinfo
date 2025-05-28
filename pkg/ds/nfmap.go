package ds

import "net"

// NFMap maps the IP address to a NHSet object.
type NFMap struct {
	// NFMap uses the type map for mapping we don't expect many clashes here.
	dmap map[string]*FTEntry
}

// NewNFMap creates a new NearFarsMap.
func NewNFMap() *NFMap {
	return &NFMap{
		dmap: make(map[string]*FTEntry),
	}
}

// Add inserts an IP address with its associated X struct.
func (nfm *NFMap) Add(ip net.IP, value *FTEntry) {
	nfm.dmap[ip.String()] = value
}

// Lookup checks if the given IP address is contained in the struct
// and similar to map, it returns a pointer and bool. Bool is true
// if there is a object.
func (nfm *NFMap) Lookup(ip net.IP) (*FTEntry, bool) {
	value, ok := nfm.dmap[ip.String()]
	return value, ok
}
