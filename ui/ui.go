package ui

import (
	"MITM/iface"
	"MITM/lanscanner"
	"MITM/poisoner"
	"MITM/sniffer"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"net"
)

type model struct {
	iface    iface.Iface
	poisoner *poisoner.Poisoner

	table       table.Model
	poisonedIps map[string]struct{}
}

type DevicesMsg []lanscanner.Device

const (
	numColumns     = 4
	numRows        = 10
	poisonedStatus = "POISONED"
	defaultStatus  = ""
	gatewayStatus  = "GATEWAY"
	usStatus       = "US"
)

func GetProgram(iface iface.Iface) *tea.Program {
	columns := [numColumns]table.Column{
		{Title: "Status", Width: 9},
		{Title: "IP", Width: 15},
		{Title: "MAC", Width: 17},
		{Title: "Manufacturer", Width: 40},
	}

	t := table.New(
		table.WithColumns(columns[:]),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(numRows),
	)

	m := model{iface, nil, t, map[string]struct{}{}}
	return tea.NewProgram(&m)
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.unPoisonEveryone()
			if sniffer.IsRoot() {
				sniffer.DisableIpForwarding()
			}
			return m, tea.Quit
		case "enter", " ":
			return m.togglePoisoning()
		}
	case DevicesMsg:
		return m.updateDeviceList(msg)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	s := fmt.Sprintf("# devices found: %d\n", len(m.table.Rows()))
	s += m.table.View()
	return s
}

func (m *model) unPoisonEveryone() {
	for _, row := range m.table.Rows() {
		ipAsString := row[1]
		macAsString := row[2]

		_, isPoisoned := m.poisonedIps[ipAsString]
		if isPoisoned {
			ip := net.ParseIP(ipAsString)[12:]
			mac, err := net.ParseMAC(macAsString)
			if err != nil {
				panic(err)
			}

			m.poisoner.RemoveVictim(ip, mac)
		}
	}
}

func (m *model) togglePoisoning() (tea.Model, tea.Cmd) {
	idx := m.table.Cursor()
	selectedRow := m.table.SelectedRow()
	status := selectedRow[0]
	selectedIPAsString := selectedRow[1]
	selectedIP := net.ParseIP(selectedIPAsString)[12:]
	selectedMACAsString := selectedRow[2]
	selectedMAC, err := net.ParseMAC(selectedMACAsString)
	if err != nil {
		panic(err)
	}

	if status == defaultStatus && m.poisoner != nil {
		m.poisonedIps[selectedIPAsString] = struct{}{}
		m.table.Rows()[idx][0] = poisonedStatus
		m.table.UpdateViewport()
		m.poisoner.AddVictim(selectedIP, selectedMAC)
	} else if status == poisonedStatus {
		delete(m.poisonedIps, selectedIPAsString)
		m.table.Rows()[idx][0] = defaultStatus
		m.table.UpdateViewport()
		m.poisoner.RemoveVictim(selectedIP, selectedMAC)
	}

	return m, nil
}

func (m *model) updateDeviceList(msg DevicesMsg) (tea.Model, tea.Cmd) {
	newRows := make([]table.Row, len(msg))

	selectedRow := m.table.SelectedRow()
	var selectedIP string
	if selectedRow == nil {
		selectedIP = ""
	} else {
		selectedIP = selectedRow[1]
	}

	cursor := m.table.Cursor()

	for idx, device := range msg {
		m.createPoisonerIfItDoesNotExist(device.Ip, device.Mac)
		if device.Ip.String() == selectedIP {
			cursor = idx
		}

		newRow := make(table.Row, numColumns)

		newRow[0] = m.getStatus(device.Ip, device.Mac)
		newRow[1] = device.Ip.String()
		newRow[2] = device.Mac.String()
		newRow[3] = device.Manufacturer
		newRows[idx] = newRow
	}

	m.table.SetRows(newRows)
	m.table.SetCursor(cursor)

	return m, nil
}

func (m *model) getStatus(ip net.IP, mac net.HardwareAddr) string {
	if m.iface.GatewayIP.String() == ip.String() {
		return gatewayStatus
	}

	if m.iface.HardwareAddr.String() == mac.String() {
		return usStatus
	}

	_, isPoisoned := m.poisonedIps[ip.String()]
	if isPoisoned {
		return poisonedStatus
	} else {
		return defaultStatus
	}
}

func (m *model) createPoisonerIfItDoesNotExist(ip net.IP, mac net.HardwareAddr) {
	if m.poisoner == nil && m.iface.GatewayIP.String() == ip.String() {
		m.poisoner = poisoner.NewPoisoner(m.iface, mac)
		go m.poisoner.Run()
	}
}
