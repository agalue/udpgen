package generator

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type UDPGenerator interface {
	PacketDescription() string
	Init(cfg *Config)
	Start(ctx context.Context)
}

type Run func(tx context.Context, stats *Stats)

type PropertiesFlag []string

func (p *PropertiesFlag) String() string {
	return strings.Join(*p, ", ")
}

func (p *PropertiesFlag) Set(value string) error {
	*p = append(*p, value)
	return nil
}

type Config struct {
	Host             string
	Port             int
	Workers          int
	PacketsPerSecond int

	// SNMP Trap Parameters
	TrapVersion        string
	TrapSource         string
	TrapID             string // Trap ID for v2c/v3, or Enterprise for v1
	TrapGeneric        int    // For v1
	TrapSpecific       int    // For v1
	TrapEngineID       string // For v3
	TrapUser           string // For v3
	TrapAuthPassphrase string // For v3
	TrapPrivPassphrase string // For v3

	TrapVarbinds PropertiesFlag

	// Syslog Parameters
	SyslogFacility int
	SyslogMessage  string
}

func (cfg *Config) IsSnmpV1() bool {
	return cfg.TrapVersion == "v1"
}

func (cfg *Config) TickDuration() time.Duration {
	return time.Duration((cfg.Workers * 1000000000) / cfg.PacketsPerSecond)
}

func (cfg *Config) UDPConn() (*net.UDPConn, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	return net.DialUDP("udp", nil, udpAddr)
}

func (cfg *Config) StartWorkers(ctx context.Context, run Run) {
	stats := new(Stats)
	go stats.Start(ctx)
	wg := new(sync.WaitGroup)
	wg.Add(cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		go func() {
			defer wg.Done()
			run(ctx, stats)
		}()
	}
	wg.Wait()
}

type Stats struct {
	Packets int64
}

func (s *Stats) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
		case <-ticker.C:
			log.Printf("Sent %d packets", s.Packets)
		}
	}
}

func (s *Stats) Inc() {
	atomic.AddInt64(&s.Packets, 1)
}
