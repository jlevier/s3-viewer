package table

import (
	"os"

	spin "s3-viewer/ui/components/spinner"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func New(c []Column, hasFiltering bool) *Model {
	m := Model{
		columns:             c,
		highlightedRowIndex: 0,
		firstVisibleRow:     0,
		hasFiltering:        hasFiltering,
		spinner:             spin.GetFooterSpinner(),
	}

	if hasFiltering {
		m.hasFiltering = hasFiltering
		m.filterInput = textinput.New()
		m.filterInput.Placeholder = "filter"
		m.filterInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		m.filterInput.CursorStyle = m.filterInput.PromptStyle.Copy()
		m.filterInput.CharLimit = 50
		m.filterInput.Width = 20
	}

	return &m
}

func (m *Model) SetData(r []Row) {
	m.data = r
	m.highlightedRowIndex = 0
	m.isLoading = false
}

func (m *Model) SetHasNextPage(hasNextPage bool) {
	m.hasNextPage = hasNextPage
}

func (m *Model) SetFooterInfo(f string) {
	m.footerInfo = f
}

func (m *Model) GetHighlightedRow() *Row {
	if len(m.data) > 0 {
		return &m.data[m.highlightedRowIndex]
	}

	return nil
}

func (m *Model) getVisibleRowCount() int {
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))
	calc := height - 6
	lastRow := len(m.data)

	if len(m.data) > calc {
		lastRow = calc - 1
	}

	return lastRow
}

func (m *Model) IsFilterVisible() bool {
	return m.isFilterVisible
}

func (m *Model) GetCurrentFilter() string {
	return m.currentFilter
}

func (m *Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, m.spinner.Tick)

	if m.hasFiltering {
		cmds = append(cmds, textinput.Blink)
	}

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.handleUpKey()

		case "down":
			m.handleDownKey()

		case "right":
			m.handleRightKey(&cmds)

		case "left":
			m.handleLeftKey(&cmds)

		case "esc":
			m.handleEscapeKey(&cmds)

		case "/":
			rm, rc := m.handleSlashKey()
			if rm != nil {
				return rm, rc
			}

		case "enter":
			m.handleEnterKey(&cmds)
		}
	}

	if m.isFilterVisible {
		var filterCmd tea.Cmd
		m.filterInput, filterCmd = m.filterInput.Update(msg)
		cmds = append(cmds, filterCmd)
	}

	if m.isLoading {
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	filter := m.renderFilter()
	h := m.renderHeader()
	r := m.renderRows()
	f := m.renderFooter()

	j := borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, h, r, f))
	jf := lipgloss.JoinVertical(lipgloss.Center, filter, j)

	return jf
}
