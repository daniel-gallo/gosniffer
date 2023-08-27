package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gosniffer/persistance"
	"time"
)

type logsModel struct {
	table table.Model
}

type LogMessage []persistance.Log

func newLogsModel() tea.Model {
	columns := []table.Column{
		{Title: "Module", Width: 6},
		{Title: "When", Width: 14},
		{Title: "IP", Width: 30},
		{Title: "Message", Width: 50},
	}

	// Don't highlight selected row
	s := table.DefaultStyles()
	s.Selected = lipgloss.NewStyle()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithHeight(NumRows),
		// Don't highlight selected row
		table.WithStyles(s),
	)

	return logsModel{
		table: t,
	}
}

func (m logsModel) Init() tea.Cmd {
	return nil
}

func (m logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LogMessage:
		newRows := make([]table.Row, len(msg))

		for idx, log := range msg {
			timeElapsed := time.Since(log.Timestamp).Truncate(time.Second)

			row := table.Row{
				log.Module,
				timeElapsed.String() + " ago",
				log.Ip.String(),
				log.Message,
			}

			newRows[idx] = row
		}

		m.table.SetRows(newRows)
	}
	return m, nil
}

func (m logsModel) View() string {
	return m.table.View() + "\n"
}
