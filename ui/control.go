package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	switch m.currentPage {
	case Main:
		return m.MainInit()
	default:
		return m.CredsInit()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.currentPage {
	case Main:
		return m.GetMainUpdate(msg)
	default:
		return m.GetCredsUpdate(msg)
	}
}

func (m Model) View() string {
	switch m.currentPage {
	case Main:
		return m.GetMainView()
	default:
		return m.GetCredsView()
	}
}
