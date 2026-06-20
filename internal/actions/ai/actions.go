package ai

import "anthodev/codory/internal/actions"

func init() {
	registry := actions.GlobalRegistry()

	aiCategory := &actions.Category{
		ID:          "ai",
		Name:        "AI",
		Description: "AI tools you can install on your system",
		Actions:     make([]*actions.Action, 0),
	}

	registry.RegisterCategory("root", aiCategory)
	registry.RegisterAction("ai", InstallHerdr())
}
