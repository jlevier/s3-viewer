package table

import (
	"github.com/charmbracelet/lipgloss"
)

type Column struct {
	Width int
	Name  string
}

type Row []string

type Model struct {
	columns []Column
	data    []Row
}

func New(c []Column) *Model {
	return &Model{
		columns: c,
	}
}

func (m *Model) SetData(r []Row) {
	m.data = r
}

func (m *Model) View() string {
	h := m.getHeaderView()
	r := m.getRows()

	return lipgloss.JoinVertical(lipgloss.Center, h, r)
}

func (m *Model) getHeaderView() string {
	s := make([]string, len(m.columns))

	for _, c := range m.columns {
		style := lipgloss.NewStyle().Width(c.Width)
		s = append(s, style.Render(c.Name))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, s...)
}

func (m *Model) getRows() string {
	s := make([]string, len(m.data))

	for i, r := range m.data {
		row := make([]string, len(m.columns))

		for j, c := range m.columns {
			style := lipgloss.NewStyle().Width(c.Width)
			row[j] = style.Render(r[j])
		}

		s[i] = lipgloss.JoinHorizontal(lipgloss.Center, row...)
	}

	return lipgloss.JoinVertical(lipgloss.Center, s...)
}
