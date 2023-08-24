package main

import (
	"MITM/iface"
	"MITM/lanscanner"
	"MITM/persistance"
	"MITM/sniffer"
	"MITM/sniffer/modules"
	"MITM/ui"
	"fmt"
	"net"
)

func main() {
	// runPersistance()
	// runSniffer()
	runUI()
}

func runPersistance() {
	sqlite := persistance.CreateSQLite("test.db")
	sqlite.Save("test", net.IP{192, 168, 1, 1}, net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, "testing message")
}

func runSniffer() {
	validIface := iface.GetSingleIface()
	repository := persistance.CreateSQLite("test.db")
	snifferModules := []sniffer.Module{
		modules.CreateDNSModule(repository),
	}

	sniffer.Sniff(validIface, snifferModules)
}

func runUI() {
	validIface := iface.GetSingleIface()
	scanner := lanscanner.NewScanner(validIface)
	uiProgram := ui.GetProgram(validIface)

	repository := persistance.CreateSQLite("test.db")
	dnsModule := modules.CreateDNSModule(repository)

	if sniffer.IsRoot() {
		sniffer.EnableIpForwarding()
	} else {
		fmt.Println("Run as root to enable IP forwarding")
	}

	go func() {
		scanner.Scan(func(devices []lanscanner.Device) {
			uiProgram.Send(ui.DevicesMsg(devices))
		})
	}()

	go func() {
		sniffer.Sniff(validIface, []sniffer.Module{dnsModule})
	}()

	_, err := uiProgram.Run()
	if err != nil {
		panic(err)
	}
}
