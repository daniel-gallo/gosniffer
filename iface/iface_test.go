package iface

import (
	"fmt"
	"net"
	"testing"
)

func TestGetAllIPs(t *testing.T) {
	iface := Iface{
		IPAddr: net.IP{192, 168, 1, 40},
		Mask:   net.IPv4Mask(255, 255, 255, 0),
	}

	actualIPs := iface.GetAllIPs()
	if len(actualIPs) != 254 {
		t.Errorf("The number of IPs returned should be 254\n")
	}

	for n := 1; n <= 254; n++ {
		expectedIP := fmt.Sprintf("192.168.1.%d", n)
		actualIP := actualIPs[n-1].String()

		if expectedIP != actualIP {
			t.Errorf("Expecting %s got got %s\n", expectedIP, actualIP)
		}
	}
}
