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

type helpItem struct {
	key, desc string
}

func GetBucketsHelp() string {
	items := []helpItem{
		{key: "\u2191", desc: "up"},
		{key: "\u2193", desc: "down"},
		{key: "enter", desc: "open folder"},
		{key: "ctrl + c", desc: "quit"},
	}

	return renderHelpItems(items)
}

func GetFilesHelp(filterPromptVisible bool, currentFilter string) string {
	items := make([]helpItem, 0)
	if !filterPromptVisible {
		items = append(items, helpItem{key: "\u2191", desc: "up"})
		items = append(items, helpItem{key: "\u2193", desc: "down"})
		items = append(items, helpItem{key: "enter", desc: "open folder"})
		items = append(items, helpItem{key: "/", desc: "filter"})
	}

	if filterPromptVisible {
		items = append(items, helpItem{key: "enter", desc: "apply filter"})
		items = append(items, helpItem{key: "esc", desc: "exit filter"})
	} else if currentFilter != "" {
		items = append(items, helpItem{key: "esc", desc: "clear filter"})
	}

	items = append(items, helpItem{key: "ctrl + c", desc: "quit"})

	return renderHelpItems(items)
}

func renderHelpItems(items []helpItem) string {
	var s strings.Builder

	for i, h := range items {
		s.WriteString(styleKey.Render(h.key))
		s.WriteString(styleDesc.Render(fmt.Sprintf(" %s", h.desc)))
		if i < len(items)-1 {
			s.WriteString(styleDesc.Render(" \u2022 "))
		}
	}

	return s.String()
}
