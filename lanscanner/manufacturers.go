package lanscanner

import (
	"bufio"
	"net"
	"os"
	"path"
	"runtime"
	"strings"
)

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

	manufacturer, ok := macToManufacturer[macAsString[:13]]
	if ok {
		return manufacturer
	}

	manufacturer, ok = macToManufacturer[macAsString[:10]]
	if ok {
		return manufacturer
	}

	manufacturer, ok = macToManufacturer[macAsString[:8]]
	if ok {
		return manufacturer
	}

	return ""
}
