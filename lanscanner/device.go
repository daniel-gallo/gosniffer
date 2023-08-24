package lanscanner

import (
	"encoding/binary"
	"net"
)

type Device struct {
	Ip           net.IP
	Mac          net.HardwareAddr
	Manufacturer string
}

type ByIp []Device

func (devices ByIp) Len() int {
	return len(devices)
}

func (devices ByIp) Swap(i, j int) {
	devices[i], devices[j] = devices[j], devices[i]
}

func (devices ByIp) Less(i, j int) bool {
	firstIpAsUInt := binary.BigEndian.Uint32(devices[i].Ip)
	secondIpAsUInt := binary.BigEndian.Uint32(devices[j].Ip)

	return firstIpAsUInt < secondIpAsUInt
}
