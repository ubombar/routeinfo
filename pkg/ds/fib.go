package ds

import (
	"fmt"
	"net"
	"strings"
)

// FI stands for forwarding information base
type FIB struct {
	fibs                map[string]*FT
	optimizeForIPv4     bool
	defaultPrefixLength uint
}

// Creates a new forwarding information base.
func NewFIB(size uint, optimizeForIPv4 bool, defaultPrefixLength uint) *FIB {
	return &FIB{
		fibs:                make(map[string]*FT, size),
		optimizeForIPv4:     optimizeForIPv4,
		defaultPrefixLength: defaultPrefixLength,
	}
}

// Gets the FI of the router address.
func (f *FIB) Lookup(address net.IP) (*FT, bool) {
	key := address.To16().String()
	if ft, ok := f.fibs[key]; !ok || ft == nil {
		return nil, false
	} else {
		return ft, true
	}
}

// Inserts a new forwarding info as defined in the forwarding info design document.
func (f *FIB) Insert(address net.IP, network net.IPNet, nexthop net.IP) error {
	key := address.To16().String()
	if ft, ok := f.fibs[key]; !ok || ft == nil {
		entry := NewFowardingTable(f.optimizeForIPv4, f.defaultPrefixLength)
		err := entry.Insert(network, nexthop)
		if err != nil {
			return err
		}

		f.fibs[key] = entry
		return nil
	} else {
		return ft.Insert(network, nexthop)
	}
}

// Converts the forwarding table into a String
func (f *FIB) String() string {
	var sb strings.Builder

	for k, v := range f.fibs {
		sb.WriteString(fmt.Sprintf("%v:\n%v\n", k, v))
	}

	return sb.String()
}

// To CSV
func (f *FIB) ToCSV() string {
	var sb strings.Builder

	for nearAddress, ftObj := range f.fibs {
		ftObj.tree.Walk(func(networkKey string, ftEntry interface{}) bool {
			if entry, ok := ftEntry.(*FTEntry); !ok {
				return true
			} else {
				for _, farAddress := range entry.dset {
					networkPrefix, err := KeyToIP(networkKey, f.optimizeForIPv4)
					if err != nil {
						return false
					}
					network := IPToNetwork(networkPrefix, int(f.defaultPrefixLength))
					sb.WriteString(fmt.Sprintf("%v, %v, %v\n", nearAddress, network, farAddress))
				}
			}

			return true
		})
	}

	return sb.String()
}
