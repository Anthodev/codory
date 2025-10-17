package tools

import "anthodev/codory/internal/actions"

func init() {
	registry := actions.GlobalRegistry()

	toolsCategory := &actions.Category{
		ID:          "tools",
		Name:        "Tools",
		Description: "Useful tools that you can add to your system",
		Actions:     make([]*actions.Action, 0),
	}

	registry.RegisterCategory("root", toolsCategory)

	// Register actions
	registry.RegisterAction("tools", InstallBat())
}
