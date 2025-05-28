package main

import (
	"fmt"
	"net"

	"github.com/ubombar/routeinfo/pkg/ds"
)

func main() {
	f := ds.NewFowardingTable()

	linksCh := ds.ReadLinkRecords("./data/links__58cb52ec_5ee7_45be_8797_e019a2815a2b__f82cf048_aff0_4ead_96f7_3e05aa4b9b14.csv", -1, 100)
	i := 0
	total := 39829550

	for l := range linksCh {
		if l.NearAddr == "::" || l.FarAddr == "::" {
			continue
		}

		if i%10000 == 0 {
			fmt.Printf("Done %v/%v (%02f%%)\n", i, total, 100*float64(i)/float64(total))
		}

		destinationAddress := net.ParseIP(l.ProbeDstAddr).To16()
		prefixLength := 24

		destinationNetwork := ds.IPToNetwork(destinationAddress, prefixLength)
		// nearAddress := net.ParseIP(l.NearAddr).To16()
		farAddress := net.ParseIP(l.FarAddr).To16()

		f.Insert(destinationNetwork, farAddress)
		i += 1
	}

	fmt.Println("Loaded")

	fmt.Println(f.Lookup(net.ParseIP("::ffff:223.255.246.224")))

	// for {
	// }

	// root := &structures.Node{}
	//
	// // Insert some IP prefixes
	// _, ipNet1, _ := net.ParseCIDR("192.168.1.0/24")
	// _, ipNet2, _ := net.ParseCIDR("192.168.0.0/16")
	// _, ipNet3, _ := net.ParseCIDR("10.0.0.0/8")
	//
	// root.Insert(ipNet1, "Office Network")
	// root.Insert(ipNet2, "Local Network")
	// root.Insert(ipNet3, "Private Network")
	//
	// // Search for IP addresses
	// ip := net.ParseIP("192.168.1.42")
	// if value, found := root.Search(ip); found {
	// 	fmt.Printf("%s belongs to %v\n", ip, value)
	// } else {
	// 	fmt.Printf("%s not found in any prefix\n", ip)
	// }
	//
	// ip = net.ParseIP("10.1.2.3")
	// if value, found := root.Search(ip); found {
	// 	fmt.Printf("%s belongs to %v\n", ip, value)
	// } else {
	// 	fmt.Printf("%s not found in any prefix\n", ip)
	// }
	//
	// ip = net.ParseIP("8.8.8.8")
	// if value, found := root.Search(ip); found {
	// 	fmt.Printf("%s belongs to %v\n", ip, value)
	// } else {
	// 	fmt.Printf("%s not found in any prefix\n", ip)
	// }
}

// func main() {
// 	// linksCh := structures.ReadLinkRecords("./data/links__58cb52ec_5ee7_45be_8797_e019a2815a2b__f82cf048_aff0_4ead_96f7_3e05aa4b9b14.csv", 1000, 100)
// 	//
// 	// for l := range linksCh {
// 	// 	if l.NearAddr == "::" || l.FarAddr == "::" {
// 	// 		continue
// 	// 	}
// 	//
// 	// 	fmt.Printf("l: %v\n", l)
// 	// }
// 	// tree := structures.NewFIB()
// 	//
// 	// // Insert some CIDRs
// 	// tree.Insert("192.168.0.0/16", "Network A")
// 	// tree.Insert("192.168.1.0/24", "Network B")
// 	// tree.Insert("10.0.0.0/8", "Network C")
// 	//
// 	// testIPs := []string{
// 	// 	"192.168.1.100",
// 	// 	"192.168.2.1",
// 	// 	"10.5.6.7",
// 	// 	"8.8.8.8",
// 	// }
// 	//
// 	// for _, ipStr := range testIPs {
// 	// 	ip := net.ParseIP(ipStr)
// 	// 	if value, found := tree.LongestPrefixMatch(ip); found {
// 	// 		fmt.Printf("IP %s matched to %v\n", ipStr, value)
// 	// 	} else {
// 	// 		fmt.Printf("IP %s not found\n", ipStr)
// 	// 	}
// 	// }
// }
