package generator

import (
	"context"
	"log"
	"sync"
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
	stats := new(Stats)
	go stats.Start(ctx)
	wg := new(sync.WaitGroup)
	for i := 0; i < gen.config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gen.startWorker(ctx, stats)
		}()
	}
	wg.Wait()
}

func (gen *Trap) startWorker(ctx context.Context, stats *Stats) {
	session := gen.createSession(ctx)
	if err := session.Connect(); err != nil {
		log.Fatalf("Cannot create SNMP session: %v", err)
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
			session.SendTrap(trap)
			stats.Inc()
		}
	}
}

func (gen *Trap) createSession(ctx context.Context) gosnmp.GoSNMP {
	return gosnmp.GoSNMP{
		Target:    gen.config.Host,
		Port:      uint16(gen.config.Port),
		Version:   gosnmp.Version1,
		Context:   ctx,
		Community: "public",
		Timeout:   2 * time.Second,
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
