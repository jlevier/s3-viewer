package main

import (
	"fmt"
	"os"
	"s3-viewer/ui/control"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(control.Model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There has been a problem: %s", err)
		os.Exit(1)
	}
}
