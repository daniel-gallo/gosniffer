package lanscanner

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"gosniffer/arp"
	"gosniffer/iface"
	"net"
	"sort"
	"time"
)

const (
	timeBetweenScans = 5 * time.Second
)

type DeviceList []Device

type Scanner struct {
	iface             iface.Iface
	MacToManufacturer map[string]string
	// TODO: this map should have a ttl
	IpToDevice map[string]Device
}

func NewScanner(iface iface.Iface) Scanner {
	macToManufacturer := GetMACToManufacturer()
	return Scanner{
		iface,
		macToManufacturer,
		make(map[string]Device),
	}
}

func (scanner *Scanner) Scan(callback func(devices DeviceList)) {
	handle := scanner.iface.GetHandle()
	err := handle.SetBPFFilter("arp")
	if err != nil {
		panic(err)
	}

	go scanner.readARP(handle, callback)
	scanner.sendARPBroadcast(handle)
	scanner.sendARPBroadcast(handle)
	scanner.sendARPBroadcast(handle)
	for {
		scanner.sendARPBroadcast(handle)
		time.Sleep(timeBetweenScans)
	}
}

func (scanner *Scanner) sendARPBroadcast(handle *pcap.Handle) {
	for _, ip := range scanner.iface.GetAllIPs() {
		arpPacket := arp.GetARPRequestPacket(scanner.iface.IPAddr, scanner.iface.HardwareAddr, ip)

		err := handle.WritePacketData(arpPacket)
		if err != nil {
			panic(err)
		}
	}
}

func (scanner *Scanner) readARP(handle *pcap.Handle, callback func(devices DeviceList)) {
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)

	for packet := range src.Packets() {
		layer := packet.Layer(layers.LayerTypeARP)
		if layer == nil {
			continue
		}

		arpLayer := layer.(*layers.ARP)
		if arpLayer.Operation != layers.ARPReply {
			continue
		}

		ip := net.IP(arpLayer.SourceProtAddress)
		mac := net.HardwareAddr(arpLayer.SourceHwAddress)
		manufacturer := GetManufacturer(scanner.MacToManufacturer, mac)

		if scanner.isUsPoisoningSomeone(ip, mac) {
			continue
		}

		scanner.IpToDevice[ip.String()] = Device{ip, mac, manufacturer}

		deviceList := scanner.getDeviceList()
		callback(deviceList)
	}
}

func (scanner *Scanner) isUsPoisoningSomeone(ip net.IP, mac net.HardwareAddr) bool {
	return ip.String() == scanner.iface.GatewayIP.String() && mac.String() == scanner.iface.HardwareAddr.String()
}

func (scanner *Scanner) getDeviceList() DeviceList {
	deviceList := make(DeviceList, len(scanner.IpToDevice))

	i := 0
	for _, device := range scanner.IpToDevice {
		deviceList[i] = device
		i++
	}

	sort.Sort(ByIp(deviceList))

	return deviceList
}
