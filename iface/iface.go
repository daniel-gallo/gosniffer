package iface

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/jackpal/gateway"
)

const (
	pcapSnaplen         = 65536
	pcapPromiscuousMode = true
	pcapTimeout         = 500 * time.Millisecond
)

type Iface struct {
	Name         string
	HardwareAddr net.HardwareAddr
	IPAddr       net.IP
	GatewayIP    net.IP
	Mask         net.IPMask
}

func GetIfaces() []Iface {
	filteredIfaces := make([]Iface, 0)

	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}

		for _, addr := range addrs {
			// Check if it is an IP address
			ipAddr, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			// Check if it is an IPv4 address
			ipv4Addr := ipAddr.IP.To4()
			if ipv4Addr == nil {
				continue
			}

			// Skip localhost
			if ipv4Addr.IsLoopback() {
				continue
			}

			gatewayIP, err := gateway.DiscoverGateway()
			if err != nil {
				panic(err)
			}
			gatewayIP = gatewayIP[12:]

			// TODO: can an iface have several IPv4 addresses?
			filteredIfaces = append(filteredIfaces, Iface{
				iface.Name,
				iface.HardwareAddr,
				ipv4Addr,
				gatewayIP,
				ipAddr.Mask,
			})
		}
	}

	return filteredIfaces
}

func GetIface(ifaceName string) Iface {
	ifaces := GetIfaces()
	ifaceNames := make([]string, len(ifaces))

	for i, iface := range ifaces {
		if iface.Name == ifaceName {
			return iface
		}

		ifaceNames[i] = iface.Name
	}

	panic(fmt.Sprintf("Interface %s not found (available interfaces: %v)\n", ifaceName, ifaceNames))
}

func (iface Iface) GetAllIPs() []net.IP {
	ip := binary.BigEndian.Uint32(iface.IPAddr)
	mask := binary.BigEndian.Uint32(iface.Mask)

	networkIP := ip & mask
	broadcastIP := networkIP + ^mask

	// The IPs will range from (networkIP + 1) to (broadcastIP - 1) (both included)
	numIPs := broadcastIP - networkIP - 1
	ips := make([]net.IP, numIPs)
	buffer := make([]byte, 4)
	i := 0
	for currentIP := networkIP + 1; currentIP <= broadcastIP-1; currentIP++ {
		binary.BigEndian.PutUint32(buffer, currentIP)
		ips[i] = net.IP{buffer[0], buffer[1], buffer[2], buffer[3]}

		i++
	}

	return ips
}

func (iface Iface) GetHandle() *pcap.Handle {
	handle, err := pcap.OpenLive(iface.Name, pcapSnaplen, pcapPromiscuousMode, pcapTimeout)
	if err != nil {
		panic(err)
	}

	return handle
}
