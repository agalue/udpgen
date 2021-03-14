package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/agalue/udpgen/generator"
)

func main() {
	log.SetOutput(os.Stdout)

	var payload string
	cfg := new(generator.Config)
	flag.StringVar(&payload, "x", "syslog", "Type of payload: snmp, syslog, netflow5, netflow9")
	flag.StringVar(&cfg.Host, "h", "127.0.0.1", "Target host / IP address")
	flag.IntVar(&cfg.Port, "p", 0, "Target port (default depends on mode)")
	flag.IntVar(&cfg.PacketsPerSecond, "r", 10000, "Number of packets per second to generate")
	flag.IntVar(&cfg.Workers, "w", 1, "Number of workers (concurrent go-routines)")
	flag.Parse()

	if cfg.PacketsPerSecond <= 0 {
		log.Fatalln("Packet rate cannot be zero.")
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
	}()
	go func() {
		select {
		case <-signalChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	var udpgen generator.UDPGenerator
	switch payload {
	case "snmp":
		udpgen = new(generator.Trap)
	case "syslog":
		udpgen = new(generator.Syslog)
	case "netflow5":
		udpgen = new(generator.Netflow5)
	case "netflow9":
		udpgen = new(generator.Netflow9)
	default:
		log.Fatalf("Invalid payload type: %s", payload)
	}

	udpgen.Init(cfg)
	log.Printf("Sending %s to %s:%d at target rate of %d packets per seconds across %d worker(s).", udpgen.PacketDescription(), cfg.Host, cfg.Port, cfg.PacketsPerSecond, cfg.Workers)
	udpgen.Start(ctx)
	log.Println("Good bye!")
}
