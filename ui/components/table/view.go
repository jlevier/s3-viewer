package table

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

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
