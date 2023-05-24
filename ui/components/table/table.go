package table

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	highlightedRowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#F25D94"))

	headerRowStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true).
			BorderForeground(lipgloss.Color("#383838"))
)

type Column struct {
	Width int
	Name  string
}

type Row []string

type Model struct {
	columns             []Column
	data                []Row
	highlightedRowIndex int
	firstVisibleRow     int
}

func New(c []Column) *Model {
	return &Model{
		columns:             c,
		highlightedRowIndex: 0,
		firstVisibleRow:     0,
	}
}

func (m *Model) SetData(r []Row) {
	m.data = r
}

func (m *Model) View() string {
	h := m.renderHeader()
	r := m.renderRows()

	return lipgloss.JoinVertical(lipgloss.Center, h, r)
}

func (m *Model) renderHeader() string {
	s := make([]string, len(m.columns))

	for _, c := range m.columns {
		style := headerRowStyle.Copy().
			Width(c.Width)
		s = append(s, style.Render(strings.ToUpper(c.Name)))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, s...)
}

func (m *Model) renderRows() string {
	lastRow := m.getVisibleRowCount()

	s := make([]string, lastRow)

	index := 0
	for i := m.firstVisibleRow; i < lastRow+m.firstVisibleRow; i++ {
		row := make([]string, len(m.columns))

		for j, c := range m.columns {
			style := lipgloss.NewStyle().Width(c.Width)
			if i == m.highlightedRowIndex {
				style = highlightedRowStyle.Copy().Width(c.Width)
			}
			row[j] = style.Render(m.data[i][j])
		}

		s[index] = lipgloss.JoinHorizontal(lipgloss.Center, row...)
		index++
	}

	return lipgloss.JoinVertical(lipgloss.Center, s...)
}

func (m *Model) GetHighlightedRow() *Row {
	if len(m.data) > 0 {
		return &m.data[m.highlightedRowIndex]
	}

	return nil
}

func (m *Model) getVisibleRowCount() int {
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))
	calc := height - 4
	lastRow := len(m.data) - 1
	if len(m.data) > calc {
		lastRow = calc - 1
	}

	return lastRow
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.highlightedRowIndex > 0 {
				m.highlightedRowIndex--
			}

			// See if you are now above the first visible row and need to shift displayed rows down by one
			if m.highlightedRowIndex < m.firstVisibleRow {
				m.firstVisibleRow--
			}
		case "down":
			if m.highlightedRowIndex < len(m.data)-1 {
				m.highlightedRowIndex++
			}

			// See if you're past the end and need to shift the displayed rows up by one
			if m.highlightedRowIndex > m.getVisibleRowCount()+m.firstVisibleRow-1 {
				m.firstVisibleRow++
			}
		}
	}

	return m, nil
}
