package extractor

import (
	"io"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type DataSP struct {
	IncomingPackets int
	OutgoingPackets int
	IncomingLength  int
	OutgoingLength  int
	CumulSeq        map[int]int
	TestCumulSeq    [50]int
}

func Extract(pcapFile, localip string) ([50]int, error) {
	d, err := readThePcap(pcapFile, localip)
	if err != nil {
		return d.TestCumulSeq, err
	}
	d.TestCumulSeqMaker()
	return d.TestCumulSeq, nil
}

func readThePcap(pcapFile, localip string) (DataSP, error) {
	d := NewDataSP()
	// Open file instead of device
	handle, err := pcap.OpenOffline(pcapFile)
	if err != nil {
		return d, err
	}
	defer handle.Close()

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Initialize packet counter
	var pc int = 0
	// Initialize cumulative packets length
	var cl int = 0
	// Initialize Data
	d.IncomingPackets = 0
	d.OutgoingPackets = 0
	d.IncomingLength = 0
	d.OutgoingLength = 0
	// Flexible loop until EOF
	for {
		packet, err := packetSource.NextPacket()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error:", err)
			continue
		}
		if packet.ApplicationLayer() != nil {
			nl := gopacket.LayerString(packet.NetworkLayer())
			pl, err := packetLength(nl)
			if err != nil {
				continue
			}
			ifo, err := ifOutgoing(nl, localip)
			if err != nil {
				continue
			}
			if ifo {
				cl -= pl
				d.OutgoingPackets++
				d.OutgoingLength += pl
			} else {
				cl += pl
				d.IncomingPackets++
				d.IncomingLength += pl
			}
		}
		pc++
		d.CumulSeq[pc] = cl
	}
	return d, nil
}

func NewDataSP() DataSP {
	var d DataSP
	d.CumulSeq = make(map[int]int)
	return d
}

func ifOutgoing(s string, localip string) (bool, error) {
	pattern := regexp.MustCompile(`SrcIP=[\d.]+`)
	srcip := pattern.FindString(s)
	m, err := url.ParseQuery(srcip)
	if err != nil {
		return false, err
	}
	if m["SrcIP"][0] == localip {
		return true, nil
	}
	return false, nil
}

func packetLength(s string) (int, error) {
	pattern := regexp.MustCompile(`Length=[\d.]+`)
	srcip := pattern.FindString(s)
	m, err := url.ParseQuery(srcip)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(m["Length"][0])
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (d *DataSP) TestCumulSeqMaker() {
	keys := make([]float64, 0, len(d.CumulSeq))
	values := make([]float64, 0, len(d.CumulSeq))
	for k := range d.CumulSeq {
		keys = append(keys, float64(k))
	}
	sort.Float64s(keys)
	for _, k := range keys {
		values = append(values, float64(d.CumulSeq[int(k)]))
	}
	for i := 0; i < len(d.TestCumulSeq); i++ {
		x := float64(i) * (float64(len(keys)) / 50.0)
		x0 := int(x)
		x1 := int(x + 1.0)
		y0 := float64(d.CumulSeq[x0])
		y1 := float64(d.CumulSeq[x1])
		y := (y0*(float64(x1)-x) + y1*(x-float64(x0))) / (float64(x1) - float64(x0))
		d.TestCumulSeq[i] = int(y)
	}
}
