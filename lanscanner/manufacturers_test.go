package lanscanner

import (
	"net"
	"slices"
	"sort"
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

func TestPrefixLengths(t *testing.T) {
	if !sort.IsSorted(sort.Reverse(sort.IntSlice(PrefixLengths))) {
		t.Errorf("PrefixLengths should be ordered in reverse order (concrete prefixes should be tried first)")
	}

	for macPrefix, _ := range GetMACToManufacturer() {
		macPrefixLength := len(macPrefix)

		if !slices.Contains(PrefixLengths, macPrefixLength) {
			t.Errorf("PrefixLengths should contain %v, since %v has length %v", macPrefixLength, macPrefix, macPrefixLength)
		}
	}
}
