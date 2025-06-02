package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/ubombar/routeinfo/pkg/ds"
)

// "near_addr","far_addr","probe_dst_addr"
type NFPRecord struct {
	NearAddr     net.IP
	FarAddr      net.IP
	ProbeDstAddr net.IP
}

func (r *NFPRecord) ProbeDstNetwork(prefixLength int) (*net.IPNet, error) {
	return ds.IPToNetwork(&r.ProbeDstAddr, prefixLength)
}

// Read the recods and write them into a channel.
func ReadNFPRecordFromStdin(limit int, bufferSize int) <-chan NFPRecord {
	readCh := make(chan NFPRecord, bufferSize)
	go func() {
		defer close(readCh)
		reader := csv.NewReader(os.Stdin)

		for i := 0; i < limit || limit == -1; i++ {
			line, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				log.Printf("There was a problem trying top parse the line: %v.\n", err)
				continue
			}
			if !strings.HasPrefix(line[0], ":") { // probably the header
				continue
			}
			if len(line) != 3 {
				log.Println("There was a problem trying parse the line: number of columns are not 3.")
				continue
			}

			nearAddr := net.ParseIP(line[0]).To16()
			farAddr := net.ParseIP(line[1]).To16()
			probeDstAddr := net.ParseIP(line[2]).To16()

			record := NFPRecord{
				NearAddr:     nearAddr,
				FarAddr:      farAddr,
				ProbeDstAddr: probeDstAddr,
			}

			readCh <- record // how would that affect performance? Well, we can do it in parallel.
		}
	}()

	return readCh
}

func main() {
	prefixLength := 24  // always 24
	postfixLength := 8  // 32 - 24
	total := 1158313642 // hardcoded value for nfp recods of iris on 2025-05-05, 1.15 Billion recods.

	log.Printf("Starting to process NFP file, prefixlength=%v, total=%v.\n", prefixLength, total)

	f := ds.NewFIB(1000, true, 24)
	linksCh := ReadNFPRecordFromStdin(-1, 100)
	i := 0
	percent := 0.0
	startTime := time.Now()

	for l := range linksCh {
		// I am not doing zero check as they are assumed to be filtered out.

		if i%10000 == 0 {
			percent = 100 * float64(i) / float64(total)
			timePassedSeconds := time.Since(startTime).Seconds()
			totalTimeEstimateSeconds := (100 / percent) * timePassedSeconds
			remeaningTimeEstimateSeconds := totalTimeEstimateSeconds - timePassedSeconds

			timePassed := time.Duration(timePassedSeconds * float64(time.Second)).Truncate(time.Second)
			totalEstimation := time.Duration(totalTimeEstimateSeconds * float64(time.Second)).Truncate(time.Second)
			remeaning := time.Duration(remeaningTimeEstimateSeconds * float64(time.Second)).Truncate(time.Second)

			log.Printf("Progress: %v/%v [%.2f%%] %10v %10v %10v.\n", i, total, percent, timePassed, remeaning, totalEstimation)
		}

		destinationNetwork, err := l.ProbeDstNetwork(prefixLength)
		if err != nil {
			log.Println("There was a problem trying parse the line: number of columns are not 3.")
			continue
		}

		if err := f.Insert(&l.NearAddr, destinationNetwork, &l.FarAddr); err != nil {
			panic(err)
		}
		i += 1
	}

	log.Println("Done processing.")

	fmt.Printf("%v\n", f.ToIPInfo(postfixLength))
}
