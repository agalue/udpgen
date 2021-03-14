package generator

import (
	"context"
	"log"
	"sync"
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
