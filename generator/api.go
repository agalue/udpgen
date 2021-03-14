package generator

import (
	"context"
	"fmt"
	"net"
)

type UDPGenerator interface {
	PacketDescription() string
	Init(cfg *Config)
	Start(ctx context.Context) error
}

type Config struct {
	Host             string
	Port             int
	PacketsPerSecond int
}

func (cfg *Config) UDPConn() (*net.UDPConn, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	return net.DialUDP("udp", nil, udpAddr)
}
