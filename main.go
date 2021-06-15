package main

import (
	"context"
	"flag"
	"log"
	"log/syslog"
	"os"
	"os/signal"

	"github.com/agalue/udpgen/generator"
)

func main() {
	log.SetOutput(os.Stdout)

	var payload string
	cfg := new(generator.Config)

	flag.StringVar(&payload, "x", "syslog", "Type of payload: snmp, syslog, netflow5, netflow9")
	flag.StringVar(&cfg.Host, "h", "127.0.0.1", "Target Hostname or IP address")
	flag.IntVar(&cfg.Port, "p", 0, "Target port (default depends on mode)")
	flag.IntVar(&cfg.PacketsPerSecond, "r", 10000, "Number of packets per second to generate")
	flag.IntVar(&cfg.Workers, "w", 1, "Number of workers (concurrent go-routines)")

	flag.StringVar(&cfg.TrapVersion, "trap-version", "v1", "SNMP Trap Version: v1 or v2c")
	flag.StringVar(&cfg.TrapSource, "trap-host", "127.0.0.1", "IP Address or Hostname of the Trap Sender")
	flag.StringVar(&cfg.TrapID, "trap-id", ".1.3.6.1.1.6.3.1.1.5", "SNMPv1 Trap Enterprise or SNMPv2c Trap ID")
	flag.IntVar(&cfg.TrapSpecific, "trap-specific", 1, "SNMPv1 Trap Specific")
	flag.IntVar(&cfg.TrapGeneric, "trap-generic", 6, "SNMPv1 Trap Generic")
	flag.Var(&cfg.TrapVarbinds, "trap-varbind", "An SNMP trap varbind (can be used multiple times)\nfor instance: .1.3.6.1.6.3.1.1.5.1::ABC (octet-string assume)")

	flag.IntVar(&cfg.SyslogFacility, "syslog-facility", int(syslog.LOG_LOCAL7), "Syslog Facility, from /usr/include/sys/syslog.h")
	flag.StringVar(&cfg.SyslogMessage, "syslog-message", "%%SEC-6-IPACCESSLOGP: list in110 denied tcp 10.99.99.1(63923) -> 10.98.98.1(1521), 1 packet", "Syslog Message")

	flag.Parse()

	if cfg.PacketsPerSecond <= 0 {
		log.Fatalln("Packet rate cannot be zero.")
	}
	if cfg.TrapVarbinds == nil {
		cfg.TrapVarbinds = append(cfg.TrapVarbinds, ".1.3.6.1.6.3.1.1.5.1::ABC")
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
