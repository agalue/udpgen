package generator

import (
	"context"
	"log"
	"time"
)

type Netflow9 struct {
	config *Config
}

func (gen *Netflow9) Init(cfg *Config) {
	gen.config = cfg
	if gen.config.Port == 0 {
		gen.config.Port = 9999
	}
}

func (gen *Netflow9) PacketDescription() string {
	return "Netflow 9 Flows"
}

func (gen *Netflow9) Start(ctx context.Context) {
	log.Printf("Not implemented yet; coming soon")
	gen.config.StartWorkers(ctx, gen.startWorker)
}

func (gen *Netflow9) startWorker(ctx context.Context, stats *Stats) {
	conn, err := gen.config.UDPConn()
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
		return
	}
	ticker := time.NewTicker(gen.config.TickDuration())
	packet := gen.buildNetflow9Packet()
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			conn.Close()
			return
		case <-ticker.C:
			conn.Write(packet)
			stats.Inc()
		}
	}
}

func (gen *Netflow9) buildNetflow9Packet() []byte {
	return []byte{} // FIXME
}
