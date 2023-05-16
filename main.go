package main

import (
	"fmt"
	"os"
	"s3-viewer/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There has been a problem: %s", err)
		os.Exit(1)
	}
}
