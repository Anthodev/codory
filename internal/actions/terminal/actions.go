package terminal

import (
	"anthodev/codory/internal/actions"
)

func init() {
	registry := actions.GlobalRegistry()

	terminalCategory := &actions.Category{
		ID:          "terminal",
		Name:        "Terminal",
		Description: "Terminals you can install on your system",
		Actions:     make([]*actions.Action, 0),
	}

	registry.RegisterCategory("root", terminalCategory)

	// Register actions
	registry.RegisterAction("terminal", InstallGhostty())
	registry.RegisterAction("terminal", InstallKitty())
	registry.RegisterAction("terminal", InstallRio())
	registry.RegisterAction("terminal", InstallWarpTerminal())
	registry.RegisterAction("terminal", InstallTmux())
	registry.RegisterAction("terminal", InstallZellij())
}
