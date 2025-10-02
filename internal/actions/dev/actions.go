package dev

import "anthodev/codory/internal/actions"

func init() {
	registry := actions.GlobalRegistry()

	// Create the Dev category (visible on all platforms)
	devCategory := &actions.Category{
		ID:          "dev",
		Name:        "Development",
		Description: "Development tools and utilities",
		Actions:     make([]*actions.Action, 0),
	}

	registry.RegisterCategory("root", devCategory)

	// Register actions
	registry.RegisterAction("dev", NewUUIDv4Action())
}
