package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	categoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	actionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C91E6"))

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EE6FF8"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF0000"))

	resultStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

	warningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFA500"))

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C91E6")).
			MarginTop(1)

	platformStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)
)
