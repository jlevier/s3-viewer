package ui

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) GetMainUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) GetMainView() string {
	return "YOU ARE NOW IN THE MAIN VIEW"
}
