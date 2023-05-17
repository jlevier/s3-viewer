package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	switch m.currentPage {
	case Main:
		return m.MainInit()
	case Buckets:
		return m.BucketsInit()
	default:
		return m.CredsInit()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.currentPage {
	case Main:
		return m.GetMainUpdate(msg)
	case Buckets:
		return m.GetBucketsUpdate(msg)
	default:
		return m.GetCredsUpdate(msg)
	}
}

func (m Model) View() string {
	switch m.currentPage {
	case Main:
		return m.GetMainView()
	case Buckets:
		return m.GetBucketsView()
	default:
		return m.GetCredsView()
	}
}
