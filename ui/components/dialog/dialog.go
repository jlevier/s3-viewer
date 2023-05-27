package dialog

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	DialogBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
)

func GetLoadingDialog(msg string, s spinner.Model) string {
	// Get terminal size and place dialog in the center
	docStyle := lipgloss.NewStyle()
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	if width > 0 {
		docStyle = docStyle.MaxWidth(width)
	}
	if height > 0 {
		docStyle = docStyle.MaxHeight(height)
	}

	p := lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		DialogBoxStyle.Render(fmt.Sprintf("%s%s", s.View(), msg)),
		lipgloss.WithWhitespaceChars("Ш#"),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}))

	return docStyle.Render(p)
}
