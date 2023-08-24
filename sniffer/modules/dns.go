package modules

import (
	"MITM/persistance"
	"MITM/sniffer"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	_ "github.com/mattn/go-sqlite3"
)

type DNS struct {
	repository persistance.Repository
}

func CreateDNSModule(repository persistance.Repository) DNS {
	return DNS{repository}
}

func (DNS) GetBPFFilter() string {
	return "udp and port 53"
}

func (m DNS) ProcessPacket(packet gopacket.Packet) {
	if !isDNSQuery(packet) {
		return
	}

	macAddress := sniffer.GetMACAddress(packet)
	ipAddress := sniffer.GetIPAddress(packet)
	hostnames := getHostnames(packet)

	for _, hostname := range hostnames {
		m.repository.Save("dns", ipAddress, macAddress, hostname)
	}
}

func getHostnames(packet gopacket.Packet) []string {
	dnsLayer := packet.Layer(layers.LayerTypeDNS).(*layers.DNS)

	numNames := len(dnsLayer.Questions)
	names := make([]string, numNames)

	for i, question := range dnsLayer.Questions {
		names[i] = string(question.Name)
	}

	return names
}

func isDNSQuery(packet gopacket.Packet) bool {
	dnsLayer := packet.Layer(layers.LayerTypeDNS).(*layers.DNS)

	return dnsLayer.OpCode == layers.DNSOpCodeQuery && len(dnsLayer.Answers) == 0
}
