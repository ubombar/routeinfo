package ds

import (
	"fmt"
	"log"
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
func (f *FIB) Get(address *net.IP) (*FT, bool, error) {
	key, err := IPToKey(address)
	if err != nil {
		return nil, false, err
	}
	if ft, ok := f.fibs[key]; !ok || ft == nil {
		return nil, false, nil
	} else {
		return ft, true, nil
	}
}

// Inserts a new forwarding info as defined in the forwarding info design document.
func (f *FIB) Insert(address *net.IP, network *net.IPNet, nexthop *net.IP) error {
	key, err := IPToKey(address)
	if err != nil {
		return err
	}
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
		nearAddress, err := KeyToIP(k)
		if err != nil {
			panic(err)
		}
		sb.WriteString(fmt.Sprintf("%v:\n%v", nearAddress, v))
	}

	return sb.String()
}

// This function computes the number of hosts and number of entries for a
// given address.
func (f *FIB) ToIPInfo(postfixLength int) string {
	var sb strings.Builder

	sb.WriteString("\"address\",\"num_networks\",\"num_hosts\"\n")

	for k, v := range f.fibs {
		nearAddress, err := KeyToIP(k)
		if err != nil {
			panic(err)
		}
		num_prefix := v.tree.Len()
		num_hosts := 1 << postfixLength
		sb.WriteString(fmt.Sprintf("\"%v\",\"%v\",\"%v\"\n", nearAddress, num_prefix, num_hosts))
	}

	return sb.String()
}

// To CSV
func (f *FIB) ToCSV() string {
	var sb strings.Builder

	for nearAddressKey, ftObj := range f.fibs {
		for networkKey, entry := range ftObj.tree.ToMap() {
			if farAddresses, ok := entry.(*FTEntry); !ok {
			} else {
				for _, farAddress := range farAddresses.dset {
					networkPrefix, err := KeyToIP(networkKey)
					if err != nil {
						log.Printf("An error occured while printing: %v.\n", err)
						continue
					}
					network, err := IPToNetwork(networkPrefix, int(f.defaultPrefixLength))
					if err != nil {
						log.Printf("An error occured while printing: %v.\n", err)
						continue
					}

					nearAddress, err := KeyToIP(nearAddressKey)
					if err != nil {
						log.Printf("An error occured while printing: %v.\n", err)
						continue
					}
					sb.WriteString(fmt.Sprintf("\"%v\",\"%v\",\"%v\"\n", nearAddress, network.String(), farAddress))
				}
			}
		}
		ftObj.tree.WalkPrefix("", func(networkKey string, ftEntry interface{}) bool {
			return true
		})
	}

	return sb.String()
}
