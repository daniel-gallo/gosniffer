package main

import (
	"flag"
	"fmt"
	"gosniffer/iface"
	"gosniffer/lanscanner"
	"gosniffer/persistance"
	"gosniffer/sniffer"
	"gosniffer/sniffer/modules"
	"gosniffer/ui"
	"time"
)

const dbFilename = "logs.db"

func main() {
	ifaceName := flag.String("iface", "en0", "Interface name")
	flag.Parse()

	validIface := iface.GetIface(*ifaceName)
	repository := persistance.NewSQLite(dbFilename)

	scanner := lanscanner.NewScanner(validIface)
	uiProgram := ui.NewProgram(validIface)

	dnsModule := modules.NewDNS(repository)
	sniffingModules := []sniffer.Module{dnsModule}

	if sniffer.IsRoot() {
		sniffer.EnableIpForwarding()
	} else {
		fmt.Println("Run as root to enable IP forwarding")
	}

	go func() {
		scanner.Scan(func(deviceList lanscanner.DeviceList) {
			uiProgram.Send(deviceList)
		})
	}()

	go func() {
		for {
			uiProgram.Send(ui.LogMessage(repository.Load(ui.NumRows)))
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		sniffer.Sniff(validIface, sniffingModules, true)
	}()

	_, err := uiProgram.Run()
	if err != nil {
		panic(err)
	}
}
