package ui

import (
	"MITM/iface"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tabModel struct {
	tabNames  [2]string
	tabModels [2]tea.Model
	activeTab int
}

const (
	NumRows = 10
	UIWidth = 90
)

func (m tabModel) Init() tea.Cmd {
	return nil
}

func (m tabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			// Propagate the quit message to the different tabs
			for idx, tabModel := range m.tabModels {
				m.tabModels[idx], _ = tabModel.Update(msg)
			}

			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabNames)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		default:
			// Pass keypress to the current tab
			m.tabModels[m.activeTab], cmd = m.tabModels[m.activeTab].Update(msg)
		}
	default:
		// Propagate the (info) message to the different tabs
		for idx, tabModel := range m.tabModels {
			m.tabModels[idx], _ = tabModel.Update(msg)
		}
	}

	return m, cmd
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m tabModel) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabNames {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabNames)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		namePadding := UIWidth/len(m.tabNames) - lipgloss.Width(t)
		if namePadding < 0 {
			namePadding = 0
		}
		renderedTabs = append(renderedTabs, style.Render(t+strings.Repeat(" ", namePadding)))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width(lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize()).Render(m.tabModels[m.activeTab].View()))
	return docStyle.Render(doc.String())
}

func GetProgram(iface iface.Iface) *tea.Program {
	tabs := [2]string{"LAN", "Sniffer logs"}
	tabModules := [2]tea.Model{
		GetLANModel(iface),
		GetLogsModel(),
	}
	m := tabModel{tabNames: tabs, tabModels: tabModules}

	return tea.NewProgram(m)

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
