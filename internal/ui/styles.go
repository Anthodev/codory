package ui

import "github.com/charmbracelet/lipgloss"

const (
	macchiatoRosewater = "#F4DBD6"
	macchiatoFlamingo  = "#F0C6C6"
	macchiatoPink      = "#F5BDE6"
	macchiatoMauve     = "#C6A0F6"
	macchiatoRed       = "#ED8796"
	macchiatoPeach     = "#F5A97F"
	macchiatoYellow    = "#EED49F"
	macchiatoGreen     = "#A6DA95"
	macchiatoTeal      = "#8BD5CA"
	macchiatoSky       = "#91D7E3"
	macchiatoBlue      = "#8AADF4"
	macchiatoLavender  = "#B7BDF8"
	macchiatoText      = "#CAD3F5"
	macchiatoSubtext   = "#A5ADCB"
	macchiatoOverlay   = "#6E738D"
	macchiatoSurface   = "#363A4F"
	macchiatoBase      = "#24273A"
	macchiatoMantle    = "#1E2030"
)

var (
	appStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(macchiatoText)).
			Padding(1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(macchiatoMauve))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(macchiatoSubtext))

	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(macchiatoOverlay))

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(macchiatoSurface)).
			Padding(1)

	activePanelStyle = panelStyle.Copy().
				BorderForeground(lipgloss.Color(macchiatoMauve))

	categoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(macchiatoGreen))

	actionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(macchiatoBlue))

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(macchiatoMauve))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(macchiatoOverlay))

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(macchiatoGreen))

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(macchiatoRed))

	resultStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(macchiatoLavender))

	warningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(macchiatoPeach))

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(macchiatoSky))

	platformStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(macchiatoSubtext)).
			Italic(true)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(macchiatoBlue)).
			Padding(1)
)
