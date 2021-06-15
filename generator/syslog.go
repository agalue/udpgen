package generator

import (
	"context"
	"fmt"
	"log"
	"log/syslog"
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
	gen.config.StartWorkers(ctx, gen.startWorker)
}

func (gen *Syslog) startWorker(ctx context.Context, stats *Stats) {
	addr := fmt.Sprintf("%s:%d", gen.config.Host, gen.config.Port)
	slog, err := syslog.Dial("udp", addr, syslog.Priority(gen.config.SyslogFacility), "udpgen")
	if err != nil {
		log.Printf("Cannot connect: %v", err)
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
			slog.Info(gen.config.SyslogMessage)
			stats.Inc()
		}
	}
}
