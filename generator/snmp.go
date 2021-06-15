package generator

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
)

type Trap struct {
	startTime time.Time
	config    *Config
}

func (gen *Trap) PacketDescription() string {
	return "SNMP Traps"
}

func (gen *Trap) Init(cfg *Config) {
	gen.config = cfg
	if gen.config.Port == 0 {
		gen.config.Port = 162
	}
	gen.startTime = time.Now()
}

func (gen *Trap) Start(ctx context.Context) {
	if gen.config.IsSnmpV2c() {
		log.Printf("Please remember to set use-address-from-varbind=true on trapd-configuration.xml")
	}
	gen.config.StartWorkers(ctx, gen.startWorker)
}

func (gen *Trap) startWorker(ctx context.Context, stats *Stats) {
	session := gen.createSession(ctx)
	if err := session.Connect(); err != nil {
		log.Printf("Cannot create SNMP session: %v", err)
		return
	}
	trap := gen.buildSnmpTrap()
	ticker := time.NewTicker(gen.config.TickDuration())
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			session.Conn.Close()
			return
		case <-ticker.C:
			if _, err := session.SendTrap(trap); err != nil {
				log.Printf("Cannot send trap: %v", err)
				return
			}
			stats.Inc()
		}
	}
}

func (gen *Trap) createSession(ctx context.Context) gosnmp.GoSNMP {
	target := gen.config.Host
	if addrs, err := net.LookupIP(gen.config.Host); err == nil {
		target = addrs[0].String()
	}
	version := gosnmp.Version1
	if gen.config.IsSnmpV2c() {
		version = gosnmp.Version2c
	}
	return gosnmp.GoSNMP{
		Target:    target,
		Port:      uint16(gen.config.Port),
		Version:   version,
		Context:   ctx,
		Community: "public",
		Timeout:   2 * time.Second,
	}
}

func (gen *Trap) buildSnmpTrap() gosnmp.SnmpTrap {
	source := "127.0.0.1"
	if addrs, err := net.LookupIP(gen.config.TrapSource); err == nil {
		source = addrs[0].String()
	}
	var trap gosnmp.SnmpTrap
	if gen.config.IsSnmpV2c() {
		trap = gosnmp.SnmpTrap{
			Variables: []gosnmp.SnmpPDU{
				{
					Name:  ".1.3.6.1.2.1.1.3.0",
					Type:  gosnmp.TimeTicks,
					Value: uint32(time.Since(gen.startTime).Seconds()),
				}, {
					Name:  ".1.3.6.1.6.3.1.1.4.1.0",
					Type:  gosnmp.ObjectIdentifier,
					Value: gen.config.TrapID,
				}, {
					Name:  ".1.3.6.1.6.3.18.1.3.0",
					Type:  gosnmp.IPAddress,
					Value: source,
				},
			},
		}
	} else {
		trap = gosnmp.SnmpTrap{
			Enterprise:   gen.config.TrapID,
			AgentAddress: source,
			GenericTrap:  gen.config.TrapGeneric,
			SpecificTrap: gen.config.TrapSpecific,
			Timestamp:    uint(time.Since(gen.startTime).Seconds()),
		}
	}
	if gen.config.TrapVarbinds != nil {
		for _, varbind := range gen.config.TrapVarbinds {
			data := strings.Split(varbind, "::")
			trap.Variables = append(trap.Variables, gosnmp.SnmpPDU{
				Name:  data[0],
				Type:  gosnmp.OctetString,
				Value: data[1],
			})
		}
	}
	return trap
}
