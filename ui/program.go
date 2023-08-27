package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"gosniffer/iface"
)

const (
	NumRows = 10
	Width   = 110
)

func NewProgram(iface iface.Iface) *tea.Program {
	return tea.NewProgram(newTabModel(iface))
}
