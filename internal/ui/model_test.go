package ui

import (
	"strings"
	"testing"

	"anthodev/codory/internal/actions"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestMenuViewZeroSizeRendersDashboard(t *testing.T) {
	m := NewModel()
	view := m.View()

	for _, want := range []string{"Codory", "Path", "Actions", "Details"} {
		if !strings.Contains(view, want) {
			t.Fatalf("View() missing %q in %q", want, view)
		}
	}
}

func TestDetailPaneActionShowsDescriptionAndSourceOnly(t *testing.T) {
	action := &actions.Action{
		Name:        "Install Foo",
		Description: "Install foo tool",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformAny: {Command: "echo install", PackageSource: actions.PackageSourceBrew},
		},
	}

	view := NewModel().viewDetailPane(80, 20, nil, []*actions.Action{action})

	for _, want := range []string{"Install foo tool", "Source: brew"} {
		if !strings.Contains(view, want) {
			t.Fatalf("detail pane missing %q in %q", want, view)
		}
	}
	for _, unwanted := range []string{"Install Foo", "Type:"} {
		if strings.Contains(view, unwanted) {
			t.Fatalf("detail pane should not include %q in %q", unwanted, view)
		}
	}
}

func TestMenuStylesAvoidBackgroundFills(t *testing.T) {
	for name, style := range map[string]lipgloss.Style{
		"app":      appStyle,
		"panel":    panelStyle,
		"active":   activePanelStyle,
		"selected": selectedStyle,
	} {
		if _, ok := style.GetBackground().(lipgloss.NoColor); !ok {
			t.Fatalf("%s style has background fill %T", name, style.GetBackground())
		}
		if style.GetHeight() != 0 {
			t.Fatalf("%s style has fixed height %d", name, style.GetHeight())
		}
	}
}

func TestMenuViewCapsLinesToTerminalHeight(t *testing.T) {
	for _, size := range []struct {
		width  int
		height int
	}{
		{width: 100, height: 12},
		{width: 80, height: 10},
		{width: 30, height: 6},
	} {
		m := NewModel()
		updated, _ := m.Update(tea.WindowSizeMsg{Width: size.width, Height: size.height})
		m = updated.(Model)

		if got := lipgloss.Height(m.View()); got > size.height {
			t.Fatalf("View() height = %d, want <= %d for %dx%d", got, size.height, size.width, size.height)
		}
	}
}
