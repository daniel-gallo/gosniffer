package modules

import (
	"MITM/persistance"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type TLS struct {
	persistance.Repository
}

func CreateTLSModule(repository persistance.Repository) TLS {
	return TLS{repository}
}

func (TLS) GetBPFFilter() string {
	return "tls"
}

func (m TLS) ProcessPacket(packet gopacket.Packet) {
	layer := packet.Layer(layers.LayerTypeTLS)
	if layer == nil {
		return
	}

	tlsLayer := layer.(*layers.TLS)
	_ = tlsLayer
}
