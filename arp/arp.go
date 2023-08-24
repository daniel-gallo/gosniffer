package arp

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

func GetARPRequestPacket(srcIP net.IP, srcMAC net.HardwareAddr, dstIP net.IP) []byte {
	broadcast := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	return getARPPacket(srcIP, srcMAC, dstIP, broadcast, layers.ARPRequest)
}

func GetARPReplyPacket(srcIP net.IP, srcMAC net.HardwareAddr, dstIP net.IP, dstMAC net.HardwareAddr) []byte {
	return getARPPacket(srcIP, srcMAC, dstIP, dstMAC, layers.ARPReply)
}

func getARPPacket(srcIP net.IP, srcMAC net.HardwareAddr, dstIP net.IP, dstMAC net.HardwareAddr, operation uint16) []byte {
	eth := layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeARP,
	}

	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         operation,
		SourceHwAddress:   srcMAC,
		SourceProtAddress: srcIP,
		DstHwAddress:      dstMAC,
		DstProtAddress:    dstIP,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	err := gopacket.SerializeLayers(buffer, options, &eth, &arp)
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
