package generator

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"time"
)

const (
	UINT16_MAX     = 65535
	PAYLOAD_AVG_MD = 1024
	SECOND         = int64(time.Second)
)

type Netflow5Header struct {
	Version        uint16
	FlowCount      uint16
	SysUptime      uint32
	UnixSec        uint32
	UnixMsec       uint32
	FlowSequence   uint32
	EngineType     uint8
	EngineId       uint8
	SampleInterval uint16
}

type Netflow5Payload struct {
	SrcIP          uint32
	DstIP          uint32
	NextHopIP      uint32
	SnmpInIndex    uint16
	SnmpOutIndex   uint16
	NumPackets     uint32
	NumOctets      uint32
	SysUptimeStart uint32
	SysUptimeEnd   uint32
	SrcPort        uint16
	DstPort        uint16
	Padding1       uint8
	TcpFlags       uint8
	IpProtocol     uint8
	IpTos          uint8
	SrcAsNumber    uint16
	DstAsNumber    uint16
	SrcPrefixMask  uint8
	DstPrefixMask  uint8
	Padding2       uint16
}

type Netflow5Packet struct {
	Header  Netflow5Header
	Records []Netflow5Payload
}

type Netflow5 struct {
	config       *Config
	startTime    int64
	upTime       uint32
	flowSequence uint32
}

func (gen *Netflow5) PacketDescription() string {
	return "Netflow 5 Flows"
}

func (gen *Netflow5) Init(cfg *Config) {
	gen.startTime = time.Now().UnixNano()
	gen.config = cfg
	if gen.config.Port == 0 {
		gen.config.Port = 8877
	}
}

func (gen *Netflow5) Start(ctx context.Context) {
	gen.config.StartWorkers(ctx, gen.startWorker)
}

func (gen *Netflow5) startWorker(ctx context.Context, stats *Stats) {
	conn, err := gen.config.UDPConn()
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
		return
	}
	ticker := time.NewTicker(gen.config.TickDuration())
	packet, err := gen.buildNetflow5Packet(8)
	if err != nil {
		log.Fatalf("Cannot build: %v", err)
		return
	}
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

func (gen *Netflow5) buildNetflow5Packet(recordCount int) ([]byte, error) {
	buffer := new(bytes.Buffer)
	data := gen.createPacket(recordCount)
	err := binary.Write(buffer, binary.BigEndian, &data.Header)
	if err != nil {
		return nil, fmt.Errorf("Writing netflow header failed: %v", err)
	}
	for _, record := range data.Records {
		err := binary.Write(buffer, binary.BigEndian, &record)
		if err != nil {
			return nil, fmt.Errorf("Writing netflow record failed: %v", err)
		}
	}
	return buffer.Bytes(), nil
}

func (gen *Netflow5) createPacket(recordCount int) Netflow5Packet {
	records := []Netflow5Payload{}
	for i := 0; i < recordCount; i++ {
		records = append(records, gen.createPayload())
	}
	return Netflow5Packet{
		Header:  gen.createHeader(recordCount),
		Records: records,
	}
}

func (gen *Netflow5) createHeader(recordCount int) Netflow5Header {
	t := time.Now().UnixNano()
	sec := t / SECOND
	nsec := t - sec*SECOND
	gen.upTime = uint32((t-gen.startTime)/int64(time.Millisecond)) + 1000
	gen.flowSequence++
	return Netflow5Header{
		Version:        5,
		FlowCount:      uint16(recordCount),
		SysUptime:      gen.upTime,
		UnixSec:        uint32(sec),
		UnixMsec:       uint32(nsec),
		FlowSequence:   gen.flowSequence,
		EngineType:     1,
		EngineId:       0,
		SampleInterval: 0,
	}
}

func (gen *Netflow5) createPayload() Netflow5Payload {
	uptime := int(gen.upTime)
	uptimeEnd := uint32(uptime - gen.randomNum(10, 500))
	uptimeStart := uptimeEnd - uint32(gen.randomNum(10, 500))
	return Netflow5Payload{
		SrcIP:          rand.Uint32(),
		DstIP:          rand.Uint32(),
		NextHopIP:      rand.Uint32(),
		SrcPort:        uint16(rand.Intn(UINT16_MAX)),
		DstPort:        uint16(rand.Intn(UINT16_MAX)),
		SnmpInIndex:    uint16(rand.Intn(UINT16_MAX)),
		SnmpOutIndex:   uint16(rand.Intn(UINT16_MAX)),
		NumPackets:     uint32(rand.Intn(PAYLOAD_AVG_MD)),
		NumOctets:      uint32(rand.Intn(PAYLOAD_AVG_MD)),
		SysUptimeStart: uptimeStart,
		SysUptimeEnd:   uptimeEnd,
		Padding1:       0,
		IpProtocol:     6,
		IpTos:          0,
		SrcPrefixMask:  uint8(rand.Intn(32)),
		DstPrefixMask:  uint8(rand.Intn(32)),
		Padding2:       0,
	}
}

func (gen *Netflow5) randomNum(min, max int) int {
	return rand.Intn(max-min) + min
}
