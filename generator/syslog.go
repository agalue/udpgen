package generator

import (
	"context"
	"fmt"
	"log"
	"log/syslog"
	"sync"
	"time"
)

type Syslog struct {
	config *Config
}

func (gen *Syslog) PacketDescription() string {
	return "Syslog Messages"
}

func (gen *Syslog) Init(cfg *Config) {
	gen.config = cfg
	if gen.config.Port == 0 {
		gen.config.Port = 514
	}
}

func (gen *Syslog) Start(ctx context.Context) {
	stats := new(Stats)
	go stats.Start(ctx)
	wg := new(sync.WaitGroup)
	wg.Add(gen.config.Workers)
	for i := 0; i < gen.config.Workers; i++ {
		go func() {
			defer wg.Done()
			gen.startWorker(ctx, stats)
		}()
	}
	wg.Wait()
}

func (gen *Syslog) startWorker(ctx context.Context, stats *Stats) {
	addr := fmt.Sprintf("%s:%d", gen.config.Host, gen.config.Port)
	slog, err := syslog.Dial("udp", addr, syslog.LOG_LOCAL7, "udpgen")
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
		return
	}
	ticker := time.NewTicker(gen.config.TickDuration())
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			slog.Close()
			return
		case <-ticker.C:
			slog.Info("%%SEC-6-IPACCESSLOGP: list in110 denied tcp 10.99.99.1(63923) -> 10.98.98.1(1521), 1 packet")
			stats.Inc()
		}
	}
}
