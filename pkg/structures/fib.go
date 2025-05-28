package structures

import (
	"errors"
	"net"

	"github.com/armon/go-radix"
)

// This is the main implementation of the forwarding info as stated in the design document.
// This struct uses a redix tree to perform lookups. Note that it has only one redix tree
// and for each entry there is a mapper called NFMap. That mapper maps the given near address
// to a set of far addresses.
//
// With this we only have one radix tree containing everyting.
type FIB struct {
	tree *radix.Tree
}

// Creates a new FIB.
func NewFIB() *FIB {
	return &FIB{
		tree: radix.New(),
	}
}

// Performs a longest prefix lookup and returns the *NFMap and *NHSet objects. Returns
// false and nil for all two of the objects.
func (f *FIB) Lookup(nearAddress, farAddress *net.IP) (*NFMap, *NHSet, bool, error) {
	// The IP address is onverted into a IPv6 String.
	_, item, found := f.tree.LongestPrefix(farAddress.To16().String())
	if !found {
		return nil, nil, false, nil
	}

	if mapObject, ok := item.(*NFMap); !ok {
		return nil, nil, false, errors.New("the expected type of NFMap is not an NFMap")
	} else {
		// The IP address is converted into a IPv6 only.
		setObject, found := mapObject.Lookup(nearAddress.To16())
		if !found {
			return nil, nil, false, nil
		}
		if setObject == nil {
			panic(errors.New("nfMapObject returned nil on lookup"))
		} else {
			return mapObject, setObject, true, nil
		}
	}
}

func (f *FIB) Insert(nearAddress *net.IP, destinationNetwork *net.IPNet, farAddress *net.IP) error {
	if nearAddress == nil || farAddress == nill || destinationNetwork == nil {
		return errors.New("given addresses or the network is nil")
	}
	mapObject, setObject, found, err := f.Lookup(nearAddress, farAddress)
	if err != nil {
		return err
	}

	if found {
	}
	return false
}

func (f *FIB) LookupMap(farAddress *net.IP) (*NFMap, bool, error) {
	// The IP address is onverted into a IPv6 String.
	key := farAddress.To16().String()

	match, item, found := f.tree.Get(key)
	if !found {
		return nil, false, nil
	}

	if mapObject, ok := item.(*NFMap); !ok {
		return nil, false, errors.New("the expected type of *NFMap is not an *NFMap")
	} else {
		return mapObject, true, nil
	}
}
