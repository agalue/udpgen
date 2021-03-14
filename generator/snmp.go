package generator

import (
	"context"
	"log"
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

func (gen *Trap) Start(ctx context.Context) error {
	session := &gosnmp.GoSNMP{
		Target:    gen.config.Host,
		Port:      uint16(gen.config.Port),
		Version:   gosnmp.Version1,
		Context:   ctx,
		Community: "public",
		Timeout:   2 * time.Second,
	}
	if err := session.Connect(); err != nil {
		return err
	}
	trap := gen.buildSnmpTrap()
	ticker := time.NewTicker(time.Duration(1000000000 / gen.config.PacketsPerSecond))
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			session.Conn.Close()
			return nil
		case <-ticker.C:
			if _, err := session.SendTrap(trap); err != nil {
				log.Printf("Cannot send trap: %v", err)
			}
		}
	}
}

func (gen *Trap) buildSnmpTrap() gosnmp.SnmpTrap {
	return gosnmp.SnmpTrap{
		Variables: []gosnmp.SnmpPDU{
			{
				Name:  ".1.3.6.1.6.3.1.1.5.1",
				Type:  gosnmp.OctetString,
				Value: "ABCDEF",
			},
		},
		Enterprise:   ".1.3.6.1.1.6.3.1.1.5",
		AgentAddress: gen.config.Host,
		GenericTrap:  6,
		SpecificTrap: 1,
		Timestamp:    uint(time.Since(gen.startTime).Seconds()),
	}
}
