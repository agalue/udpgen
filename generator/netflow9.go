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

func (gen *Netflow9) Start(ctx context.Context) error {
	log.Printf("Not implemented yet; coming soon")
	conn, err := gen.config.UDPConn()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Duration(1000000000 / gen.config.PacketsPerSecond))
	packet := gen.buildNetflow9Packet()
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			conn.Close()
			return nil
		case <-ticker.C:
			conn.Write(packet)
		}
	}
}

func (gen *Netflow9) buildNetflow9Packet() []byte {
	return []byte{} // FIXME
}
