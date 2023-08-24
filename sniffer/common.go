package sniffer

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

func GetMACAddress(packet gopacket.Packet) net.HardwareAddr {
	ethLayer := packet.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)
	return ethLayer.SrcMAC
}

func GetIPAddress(packet gopacket.Packet) net.IP {
	ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ipv4Layer != nil {
		return ipv4Layer.(*layers.IPv4).SrcIP
	}

	ipv6Layer := packet.Layer(layers.LayerTypeIPv6)
	if ipv6Layer != nil {
		return ipv6Layer.(*layers.IPv6).SrcIP
	}

	panic("No IP layer")
}
