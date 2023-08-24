package sniffer

import "github.com/google/gopacket"

type Module interface {
	GetBPFFilter() string
	ProcessPacket(packet gopacket.Packet)
}
