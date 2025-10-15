package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	_ "anthodev/codory/internal/actions/dev"
	_ "anthodev/codory/internal/actions/dev/code_editors"
	_ "anthodev/codory/internal/actions/dev/docker"
	_ "anthodev/codory/internal/actions/package_managers"
	_ "anthodev/codory/internal/actions/shells"
	_ "anthodev/codory/internal/actions/shells/zsh"
	_ "anthodev/codory/internal/actions/shells/zsh/plugins"
	_ "anthodev/codory/internal/actions/terminal"
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
