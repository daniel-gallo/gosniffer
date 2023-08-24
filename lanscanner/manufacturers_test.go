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
		{net.HardwareAddr{0xdc, 0x44, 0x27, 0x10, 0x00, 0x00}, "Tesla,Inc."},
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
