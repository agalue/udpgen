package generator

import (
	"context"
	"fmt"
	"log"
	"net"
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

type Config struct {
	Host             string
	Port             int
	Workers          int
	PacketsPerSecond int
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
