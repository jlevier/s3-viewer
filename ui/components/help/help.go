package help

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8a8a8a"))

	styleDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5e5e5e"))
)

func GetBucketsHelp() string {
	var s strings.Builder

	s.WriteString(renderHelpItem("\u2191", "up", true))
	s.WriteString(renderHelpItem("\u2193", "down", true))
	s.WriteString(renderHelpItem("enter", "open folder", true))
	s.WriteString(renderHelpItem("ctrl + c", "quit", false))

	return s.String()
}

func GetFilesHelp(filterPromptVisible bool, currentFilter string) string {
	var s strings.Builder

	if !filterPromptVisible {
		s.WriteString(renderHelpItem("\u2191", "up", true))
		s.WriteString(renderHelpItem("\u2193", "down", true))
		s.WriteString(renderHelpItem("enter", "open folder", true))
		s.WriteString(renderHelpItem("/", "filter", true))
	}

	if filterPromptVisible {
		s.WriteString(renderHelpItem("enter", "apply filter", true))
		s.WriteString(renderHelpItem("esc", "exit filter", true))
	} else if currentFilter != "" {
		s.WriteString(renderHelpItem("esc", "clear filter", true))
	}

	s.WriteString(renderHelpItem("ctrl + c", "quit", false))

	return s.String()
}

func renderHelpItem(key, desc string, separator bool) string {
	var s strings.Builder
	s.WriteString(styleKey.Render(key))
	s.WriteString(styleDesc.Render(fmt.Sprintf(" %s", desc)))
	if separator {
		s.WriteString(styleDesc.Render(" \u2022 "))
	}

	return s.String()
}
