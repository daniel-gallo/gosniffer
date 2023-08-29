package lanscanner

import (
	"bufio"
	"net"
	"os"
	"path"
	"runtime"
	"strings"
)

// PrefixLengths GetManufacturer will try more concrete prefixes (i.e. longer) first
var PrefixLengths = []int{17, 14, 11, 8, 5}

func GetMACToManufacturer() map[string]string {
	macToManufacturer := make(map[string]string)

	_, filename, _, _ := runtime.Caller(1)
	directory := path.Dir(filename)
	f, err := os.Open(directory + "/manufacturers.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, "\t")
		if len(splitLine) != 2 {
			panic("Error while parsing file")
		}

		mac := splitLine[0]
		manufacturer := splitLine[1]
		macToManufacturer[mac] = manufacturer
	}

	return macToManufacturer
}

func GetManufacturer(macToManufacturer map[string]string, mac net.HardwareAddr) string {
	macAsString := mac.String()

	for _, prefixLength := range PrefixLengths {
		manufacturer, ok := macToManufacturer[macAsString[:prefixLength]]
		if ok {
			return manufacturer
		}
	}

	return ""
}
