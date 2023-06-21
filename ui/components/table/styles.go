package table

import (
	"github.com/charmbracelet/lipgloss"
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
