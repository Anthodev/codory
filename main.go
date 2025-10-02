package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	_ "anthodev/codory/internal/actions/dev"
	"anthodev/codory/internal/ui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	p := tea.NewProgram(
		ui.NewModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
