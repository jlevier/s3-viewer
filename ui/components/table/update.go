package table

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) handleUpKey() {
	if m.highlightedRowIndex > 0 && !m.isFilterVisible {
		m.highlightedRowIndex--
	}

	// See if you are now above the first visible row and need to shift displayed rows down by one
	if m.highlightedRowIndex < m.firstVisibleRow {
		m.firstVisibleRow--
	}
}

func (m *Model) handleDownKey() {
	if m.highlightedRowIndex < len(m.data)-1 && !m.isFilterVisible {
		m.highlightedRowIndex++
	}

	// See if you're past the end and need to shift the displayed rows up by one
	if m.highlightedRowIndex > m.getVisibleRowCount()+m.firstVisibleRow-1 {
		m.firstVisibleRow++
	}
}

func (m *Model) handleRightKey(cmds *[]tea.Cmd) {
	if m.hasNextPage {
		m.isLoading = true
		*cmds = append(
			make([]tea.Cmd, 1),
			func() tea.Msg {
				m.currentPageIndex++
				return NextPageMsg{
					CurrentPageIndex: m.currentPageIndex - 1,
				}
			})
	}
}

func (m *Model) handleLeftKey(cmds *[]tea.Cmd) {
	if m.currentPageIndex > 0 {
		m.isLoading = true
		*cmds = append(
			make([]tea.Cmd, 1),
			func() tea.Msg {
				m.currentPageIndex--
				return PrevPageMsg{
					CurrentPageIndex: m.currentPageIndex,
				}
			})
	}
}

func (m *Model) handleEscapeKey(cmds *[]tea.Cmd) {
	if m.isFilterVisible {
		m.isFilterVisible = false
	} else if m.currentFilter != "" {
		m.currentFilter = ""
		m.filterInput.SetValue("")

		*cmds = append(
			make([]tea.Cmd, 1),
			func() tea.Msg {
				return FilterAppliedMsg{
					Filter: m.currentFilter,
				}
			})
	}
}

func (m *Model) handleSlashKey() (*Model, tea.Cmd) {
	if m.hasFiltering {
		if !m.isFilterVisible {
			m.isFilterVisible = true
			return m, m.filterInput.Focus()
		}
	}

	return nil, nil
}

func (m *Model) handleEnterKey(cmds *[]tea.Cmd) {
	if m.isFilterVisible {
		m.isLoading = true
		*cmds = append(make([]tea.Cmd, 1), m.spinner.Tick)
		m.currentFilter = m.filterInput.Value()
		*cmds = append(
			make([]tea.Cmd, 1),
			func() tea.Msg {
				return FilterAppliedMsg{
					Filter: m.currentFilter,
				}
			})
		m.isFilterVisible = false
	}
}
