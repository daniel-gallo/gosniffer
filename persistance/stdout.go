package persistance

import (
	"fmt"
	"net"
	"time"
)

type Stdout struct{}

func (stdout Stdout) Save(module string, ip net.IP, mac net.HardwareAddr, message string) {
	fmt.Printf("[%v] %v: %v (%v) %v\n", time.Now(), module, ip, mac, message)
}
