package main

import (
	"log"

	"github.com/gopacket/gopacket/pcap"
	"github.com/mikezappa87/bgp-back/pkg/lab1"
)

func main() {
	handle, err := pcap.OpenLive("eth0", 4096, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	filter := "tcp and port 179"

	defer handle.Close()

	if filter != "" {
		log.Println("applying filter ", filter)
		err := handle.SetBPFFilter(filter)
		if err != nil {
			log.Fatalf("error applying BPF Filter %s - %v", filter, err)
		}
	}

	lab1.Setup()
	lab1.Listen(handle)
}
