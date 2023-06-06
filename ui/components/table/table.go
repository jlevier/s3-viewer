package table

import (
	"fmt"
	"os"
	"strings"

	spin "s3-viewer/ui/components/spinner"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	highlightedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#9a87a1"))

	headerRowStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true).
			BorderForeground(lipgloss.Color("#383838"))

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3C3836"))

	footerPrefixStyle = lipgloss.NewStyle().
				Width(4).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#F25D93"))

	footerNavStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#A550DF")).
			Padding(0, 1)

	footerPagingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#5CC1F7")).
				Padding(0, 1)

	footerPosStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#6124DF")).
			Padding(0, 1)

	footerFilterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#FCA17D")).
				Padding(0, 1)

	footerPathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3C3836")).
			Padding(0, 0, 0, 1)

	footerLoadingIconStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#F25D93")).
				Padding(0, 0, 0, 1)

	footerLoadingTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#F25D93")).
				Padding(0, 1, 0, 0)

	filterWrapperStyle = lipgloss.NewStyle().
				Width(30).
				AlignHorizontal(lipgloss.Left)
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
	footerInfo          string

	// Loading
	spinner   spinner.Model
	isLoading bool

	// Filter
	hasFiltering    bool
	currentFilter   string
	isFilterVisible bool
	filterInput     textinput.Model

	// Paging
	hasNextPage      bool
	currentPageIndex int
}

type FilterAppliedMsg struct {
	Filter string
}

type NextPageMsg struct {
	CurrentPageIndex int
}

type PrevPageMsg struct {
	CurrentPageIndex int
}

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

func (m *Model) renderFilter() string {
	if m.isFilterVisible {
		return filterWrapperStyle.Render(m.filterInput.View())
	}

	return ""
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

	if lastRow < 0 {
		return ""
	}

	s := make([]string, lastRow)

	index := 0
	for i := m.firstVisibleRow; i < lastRow+m.firstVisibleRow; i++ {
		row := make([]string, len(m.columns))

		for j, c := range m.columns {
			row[j] = m.renderColumn(m.data[i][j], c, i, j)
		}

		s[index] = lipgloss.JoinHorizontal(lipgloss.Center, row...)
		index++
	}

	_, height, _ := term.GetSize(int(os.Stdout.Fd()))
	h := lipgloss.NewStyle().Height(height - 10)

	return h.Render(lipgloss.JoinVertical(lipgloss.Center, s...))
}

func (m *Model) renderColumn(data string, c Column, currentRow, currentCol int) string {
	style := lipgloss.NewStyle().Width(c.Width)
	if currentRow == m.highlightedRowIndex {
		style = highlightedRowStyle.Copy().Width(c.Width)
	}
	if currentCol == 0 {
		style = style.Copy().Padding(0, 0, 0, 1)
	} else if currentCol == len(m.columns)-1 {
		style = style.Copy().Padding(0, 1, 0, 0)
	}

	dataFinal := data

	// if data has folder path, strip off all folders but the last
	folders := strings.Split(dataFinal, "/")
	if len(folders) > 1 {
		// path ended with / so last item is empty
		if folders[len(folders)-1] == "" {
			dataFinal = folders[len(folders)-2]
		} else {
			dataFinal = folders[len(folders)-1]
		}
	}

	// If data is too large for the column, to prevent wrapping, truncate and add ellipses
	calc := c.Width - 5
	if len(dataFinal) > calc && calc > 0 {
		dataFinal = fmt.Sprintf("%s...", dataFinal[:calc])
	}

	return style.Render(dataFinal)
}

func (m *Model) renderFooter() string {
	width := 0
	for _, w := range m.columns {
		width += w.Width
	}

	var left strings.Builder
	var right strings.Builder

	left.WriteString(footerPrefixStyle.Render(" .. "))
	left.WriteString(footerPathStyle.Render(m.footerInfo))

	if m.isLoading {
		right.WriteString(footerLoadingIconStyle.Render(m.spinner.View()))
		right.WriteString(footerLoadingTextStyle.Render("loading"))
	}

	if m.currentFilter != "" {
		right.WriteString(footerFilterStyle.Render(fmt.Sprintf("\uf002 %s", m.currentFilter)))
	}

	right.WriteString(m.renderPagingFooter())

	right.WriteString(footerNavStyle.Render(m.renderNavFooter()))

	right.WriteString(footerPosStyle.Render(fmt.Sprintf("%v/%v", m.highlightedRowIndex+1, len(m.data))))

	rightStyle := footerStyle.Copy().
		AlignHorizontal(lipgloss.Right).
		Width(width / 2)

	leftStyle := footerStyle.Copy().
		Width(width / 2)

	leftFinal := leftStyle.Render(left.String())
	rightFinal := rightStyle.Render(right.String())

	return lipgloss.JoinHorizontal(lipgloss.Bottom, leftFinal, rightFinal)
}

func (m *Model) renderNavFooter() string {
	var nav strings.Builder

	if m.highlightedRowIndex > 0 {
		nav.WriteString("\uf062") // up arrow
	} else {
		nav.WriteString(" ")
	}

	nav.WriteString(" ")

	if m.highlightedRowIndex < len(m.data)-1 {
		nav.WriteString("\uf063") // down arrow
	} else {
		nav.WriteString(" ")
	}

	return nav.String()
}

func (m *Model) renderPagingFooter() string {
	chars := make([]string, 0)

	chars = append(chars, "\uf405")
	chars = append(chars, fmt.Sprintf("%v", m.currentPageIndex+1))

	if m.currentPageIndex > 0 {
		chars = append(chars, "\uf04a") // left arrow
	}

	if m.hasNextPage {
		chars = append(chars, "\uf04e") // right arrow
	}

	return footerPagingStyle.Render(strings.Join(chars, " "))
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
			if m.highlightedRowIndex > 0 && !m.isFilterVisible {
				m.highlightedRowIndex--
			}

			// See if you are now above the first visible row and need to shift displayed rows down by one
			if m.highlightedRowIndex < m.firstVisibleRow {
				m.firstVisibleRow--
			}

		case "down":
			if m.highlightedRowIndex < len(m.data)-1 && !m.isFilterVisible {
				m.highlightedRowIndex++
			}

			// See if you're past the end and need to shift the displayed rows up by one
			if m.highlightedRowIndex > m.getVisibleRowCount()+m.firstVisibleRow-1 {
				m.firstVisibleRow++
			}

		case "right":
			if m.hasNextPage {
				cmds = append(
					cmds,
					func() tea.Msg {
						m.currentPageIndex++
						return NextPageMsg{
							CurrentPageIndex: m.currentPageIndex - 1,
						}
					})
			}

		case "left":
			if m.currentPageIndex > 0 {
				cmds = append(
					cmds,
					func() tea.Msg {
						m.currentPageIndex--
						return PrevPageMsg{
							CurrentPageIndex: m.currentPageIndex + 1,
						}
					})
			}

		case "esc":
			if m.isFilterVisible {
				m.isFilterVisible = false
			} else if m.currentFilter != "" {
				m.currentFilter = ""
				m.filterInput.SetValue("")

				cmds = append(
					cmds,
					func() tea.Msg {
						return FilterAppliedMsg{
							Filter: m.currentFilter,
						}
					})
			}

		case "/":
			if m.hasFiltering {
				if !m.isFilterVisible {
					m.isFilterVisible = true
					return m, m.filterInput.Focus()
				}
			}

		case "enter":
			if m.isFilterVisible {
				m.isLoading = true
				cmds = append(cmds, m.spinner.Tick)
				m.currentFilter = m.filterInput.Value()
				cmds = append(
					cmds,
					func() tea.Msg {
						return FilterAppliedMsg{
							Filter: m.currentFilter,
						}
					})
				m.isFilterVisible = false
			}
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
