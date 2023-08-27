package lanscanner

import (
	"net"
	"testing"
)

func TestGetManufacturer(t *testing.T) {
	macToManufacturer := GetMACToManufacturer()

	testCases := []struct {
		mac                  net.HardwareAddr
		expectedManufacturer string
	}{
		{net.HardwareAddr{0xe8, 0xf7, 0x91, 0x00, 0x00, 0x00}, "Xiaomi Communications Co Ltd"},
		{net.HardwareAddr{0xf4, 0xf5, 0xd8, 0x00, 0x00, 0x00}, "Google, Inc."},
		{net.HardwareAddr{0x4c, 0x4f, 0xee, 0x00, 0x00, 0x00}, "OnePlus Technology (Shenzhen) Co., Ltd"},
	}

	for _, testCase := range testCases {
		actualManufacturer := GetManufacturer(macToManufacturer, testCase.mac)

		if actualManufacturer != testCase.expectedManufacturer {
			t.Errorf("Expecting %s but got %s\n", testCase.expectedManufacturer, actualManufacturer)
		}
	}
}
