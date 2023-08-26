package persistance

import (
	"net"
	"time"
)

type Log struct {
	Module    string
	Timestamp time.Time
	Ip        net.IP
	Mac       net.HardwareAddr
	Message   string
}

type Repository interface {
	Save(module string, ip net.IP, mac net.HardwareAddr, message string)
	Load(numMessages int) []Log
}
