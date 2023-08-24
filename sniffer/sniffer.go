package sniffer

import (
	"MITM/iface"
	"github.com/google/gopacket"
)

func Sniff(iface iface.Iface, modules []Module) {
	handle := iface.GetHandle()

	bpfFilter := combineBpfFilters(modules)
	err := handle.SetBPFFilter(bpfFilter)
	if err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		for _, module := range modules {
			module.ProcessPacket(packet)
		}
	}
}

func combineBpfFilters(modules []Module) string {
	bpfFilter := ""
	for i, module := range modules {
		bpfFilter += "(" + module.GetBPFFilter() + ")"

		if i < len(modules)-1 {
			// There are more elements
			bpfFilter += " or "
		}
	}

	return bpfFilter
}
