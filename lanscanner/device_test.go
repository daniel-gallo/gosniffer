package lanscanner

import (
	"net"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	devices := []Device{
		{Ip: net.IP{192, 168, 1, 50}},
		{Ip: net.IP{192, 168, 1, 30}},
		{Ip: net.IP{192, 168, 1, 100}},
		{Ip: net.IP{192, 168, 1, 25}},
	}

	expectedIPs := []string{
		"192.168.1.25",
		"192.168.1.30",
		"192.168.1.50",
		"192.168.1.100",
	}

	sort.Sort(ByIp(devices))
	for i, device := range devices {
		expectedIp := expectedIPs[i]
		actualIp := device.Ip.String()

		if expectedIp != actualIp {
			t.Errorf("Expecting %s but got %s\n", expectedIp, actualIp)
		}
	}
}
