package help

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8a8a8a"))

	styleDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5e5e5e"))
)

func GetBucketsHelp(hasActiveSearch bool) string {
	var s strings.Builder

	s.WriteString(styleKey.Render("\u2191"))
	s.WriteString(styleDesc.Render(" up \u2022"))
	s.WriteString(styleKey.Render(" \u2193"))
	s.WriteString(styleDesc.Render(" down \u2022"))
	s.WriteString(styleKey.Render(" /"))
	s.WriteString(styleDesc.Render(" search \u2022"))
	if hasActiveSearch {
		s.WriteString(styleKey.Render("esc"))
		s.WriteString(styleDesc.Render(" clear search \u2022"))
	}
	s.WriteString(styleKey.Render(" ctrl + c"))
	s.WriteString(styleDesc.Render(" quit"))

	return s.String()
}

func GetFilesHelp(filterPromptVisible bool, currentFilter string) string {
	var s strings.Builder

	s.WriteString(styleKey.Render("\u2191"))
	s.WriteString(styleDesc.Render(" up \u2022"))
	s.WriteString(styleKey.Render(" \u2193"))
	s.WriteString(styleDesc.Render(" down \u2022"))
	if !filterPromptVisible {
		s.WriteString(styleKey.Render(" /"))
		s.WriteString(styleDesc.Render(" search \u2022"))
	}
	if filterPromptVisible {
		s.WriteString(styleKey.Render(" esc"))
		s.WriteString(styleDesc.Render(" exit search \u2022"))
	} else if currentFilter != "" {
		s.WriteString(styleKey.Render(" esc"))
		s.WriteString(styleDesc.Render(" clear filter \u2022"))
	}
	s.WriteString(styleKey.Render(" ctrl + c"))
	s.WriteString(styleDesc.Render(" quit"))

	return s.String()
}
