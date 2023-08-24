package sniffer

import (
	"github.com/google/gopacket"
	"testing"
)

type mockARPModule struct{}

func (mockARPModule) GetBPFFilter() string {
	return "arp"
}

func (mockARPModule) ProcessPacket(_ gopacket.Packet) {

}

type mockDNSModule struct{}

func (mockDNSModule) GetBPFFilter() string {
	return "udp and port 53"
}

func (mockDNSModule) ProcessPacket(_ gopacket.Packet) {

}

func TestBPFFilterCombinerWithMoreThanOneElement(t *testing.T) {
	modules := []Module{
		mockARPModule{},
		mockDNSModule{},
	}

	actualFilter := combineBpfFilters(modules)
	expectedFilter := "(arp) or (udp and port 53)"

	if actualFilter != expectedFilter {
		t.Errorf("Expected %s but got %s\n", expectedFilter, actualFilter)
	}
}

func TestBPFFilterCombinerWithOneElement(t *testing.T) {
	modules := []Module{
		mockARPModule{},
	}

	actualFilter := combineBpfFilters(modules)
	expectedFilter := "(arp)"

	if actualFilter != expectedFilter {
		t.Errorf("Expected %s but got %s\n", expectedFilter, actualFilter)
	}
}
