package persistance

import (
	"net"
	"testing"
	"time"
)

var (
	module  = "module"
	ip      = net.IP{192, 168, 1, 2}
	mac     = net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	message = "message"
)

func TestSQLitePersistence(t *testing.T) {
	tmpFolder := t.TempDir()
	sqlite := CreateSQLite(tmpFolder + "/test.db")
	sqlite.Save(module, ip, mac, message)

	logs := sqlite.Load(1)
	if len(logs) != 1 {
		t.Errorf("Expected one log, but got %v", len(logs))
	}
	log := logs[0]

	now := time.Now()
	if now.Sub(log.Timestamp) > time.Second {
		t.Errorf("The log timestamp (%v) should be closer to the current timestamp (%v)", log.Timestamp, now)
	}
	if log.Module != module {
		t.Errorf("Expected %v but got %v", module, log.Module)
	}
	if log.Ip.String() != ip.String() {
		t.Errorf("Expected %v but got %v", ip, log.Ip)
	}
	if log.Mac.String() != mac.String() {
		t.Errorf("Expected %v but got %v", mac, log.Mac)
	}
	if log.Message != message {
		t.Errorf("Expected %v but got %v", message, log.Message)
	}
}

func TestSQLiteLoadLimit(t *testing.T) {
	tmpFolder := t.TempDir()
	sqlite := CreateSQLite(tmpFolder + "/test.db")
	numLogsInDb := 10
	numLogsToRetrieve := 5

	for i := 0; i < numLogsInDb; i++ {
		sqlite.Save(module, ip, mac, message)
	}

	logs := sqlite.Load(numLogsToRetrieve)
	if len(logs) != numLogsToRetrieve {
		t.Errorf("Expected %v logs but got %v", numLogsToRetrieve, len(logs))
	}
}
