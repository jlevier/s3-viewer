package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// if m.session == nil {
	// 	ch := make(chan *s3.SessionResponse)
	// 	go s3.GetSession(ch)
	// 	resp := <-ch
	// 	if resp.Err != nil {
	// 		return Model{currentPage: Creds, session: nil}, nil
	// 	}
	// }

	return m, nil
}

func (m Model) View() string {
	return "Cached credentials not found.  Please enter"
}
