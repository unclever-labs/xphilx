package xphilx

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

// TODO: move these to actual config
var (
	defaultRetries = 3
	defaultPeriod  = time.Duration(1 * time.Second)
)

// Run starts the exporter
func Run(cfg Config) {
	var err error
	var handle *pcap.Handle

	if err = validateConfig(cfg); err != nil {
		fmt.Println("Failed validating config:", err)
		return
	}

	// Set up pcap packet capture
	fmt.Println("Starting capture on interface:", cfg.Interface)
	handle, err = pcap.OpenLive(cfg.Interface, cfg.SnapLength, true, pcap.BlockForever)
	if err != nil {
		fmt.Println("Failed starting packet filter:", err)
		return
	}

	filter := "tcp and dst port " + cfg.Port
	if err := handle.SetBPFFilter(filter); err != nil {
		fmt.Println("Failed setting BPF filter:", err)
		return
	}

	// Create s3 uploader
	uploader := s3manager.NewUploader(session.Must(session.NewSession()))
	if err = testUpload(cfg, uploader); err != nil {
		fmt.Println("Failed testing s3 upload:", err)
		return
	}

	// Setup payload consumer
	payloadCh := make(chan []byte)
	go consumePayloads(cfg, uploader, payloadCh)

	// Set up assembly
	streamFactory := &httpStreamFactory{
		payloadCh: payloadCh,
	}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	fmt.Println("Reading in packets...")
	// Read in packets, pass to assembler.
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.
			if packet == nil {
				return
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Println("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}
