package sniffer

import (
	"os/exec"
	"os/user"
)

func IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	return currentUser.Username == "root"
}

func EnableIpForwarding() {
	cmd := exec.Command("sysctl", "-w", "net.inet.ip.forwarding=1")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func DisableIpForwarding() {
	cmd := exec.Command("sysctl", "-w", "net.inet.ip.forwarding=0")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
