package structures

import (
	"encoding/csv"
	"os"
	"strconv"
)

// Struct representing a line of the CSV
type LinkRecord struct {
	ProbeProtocol  int
	ProbeSrcAddr   string
	ProbeDstPrefix string
	ProbeDstAddr   string
	ProbeSrcPort   int
	ProbeDstPort   int
	NearRound      int
	FarRound       int
	NearTTL        int
	FarTTL         int
	NearAddr       string
	FarAddr        string
}

func ReadLinkRecords(filename string, limit int, bufferSize int) <-chan LinkRecord {
	readCh := make(chan LinkRecord, bufferSize)
	go func() {
		defer close(readCh)
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader := csv.NewReader(file)

		i := 0

		for i < limit || limit == -1 {
			line, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				panic(err)
			}

			probeProtocol, err := strconv.Atoi(line[0])
			if err != nil {
				panic(err)
			}
			probeSrcAddr := line[1]
			probeDstPrefix := line[2]
			probeDstAddr := line[3]
			probeSrcPort, err := strconv.Atoi(line[4])
			if err != nil {
				panic(err)
			}
			probeDstPort, err := strconv.Atoi(line[5])
			if err != nil {
				panic(err)
			}
			nearRound, err := strconv.Atoi(line[6])
			if err != nil {
				panic(err)
			}
			farRound, err := strconv.Atoi(line[7])
			if err != nil {
				panic(err)
			}
			nearTTL, err := strconv.Atoi(line[8])
			if err != nil {
				panic(err)
			}
			farTTL, err := strconv.Atoi(line[9])
			if err != nil {
				panic(err)
			}
			nearAddr := line[10]
			farAddr := line[11]

			record := LinkRecord{
				ProbeProtocol:  probeProtocol,
				ProbeSrcAddr:   probeSrcAddr,
				ProbeDstPrefix: probeDstPrefix,
				ProbeDstAddr:   probeDstAddr,
				ProbeSrcPort:   probeSrcPort,
				ProbeDstPort:   probeDstPort,
				NearRound:      nearRound,
				FarRound:       farRound,
				NearTTL:        nearTTL,
				FarTTL:         farTTL,
				NearAddr:       nearAddr,
				FarAddr:        farAddr,
			}

			readCh <- record
			i += 1
		}
	}()

	return readCh
}
