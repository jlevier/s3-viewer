package ui

import "github.com/charmbracelet/lipgloss"

var (
	spinnerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	DialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
)
