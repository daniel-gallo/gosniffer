package persistance

import (
	"net"
)

type Repository interface {
	Save(module string, ip net.IP, mac net.HardwareAddr, message string)
}
