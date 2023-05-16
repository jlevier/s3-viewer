package ui

import tea "github.com/charmbracelet/bubbletea"

type CredsModel struct {
	invalid bool
}

func (m *Model) GetCredsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()

		switch k {
		case "m":
			return Model{currentPage: Main, session: nil}, nil
		case "enter":
			return Model{currentPage: Creds, session: nil, creds: CredsModel{invalid: true}}, nil
		}
	}

	return m, nil
}

func (m *Model) GetCredsView() string {
	s := "Cached credentials not found.  Please enter"

	if m.creds.invalid {
		s += "\n You have entered invalid credentials!"
	}

	return s
}
