package dev

import (
	"anthodev/codory/internal/actions"
)

func init() {
	registry := actions.GlobalRegistry()

	// Ensure dev category exists for proper nesting
	// This handles initialization order issues gracefully
	if _, exists := registry.GetCategory("dev"); !exists {
		// Create dev category if not already registered
		// This is safe because init() functions are called sequentially
		devCategory := &actions.Category{
			ID:            "dev",
			Name:          "Development",
			Description:   "Development tools and utilities",
			SubCategories: make([]*actions.Category, 0),
			Actions:       make([]*actions.Action, 0),
		}
		registry.RegisterCategory("root", devCategory)
	}

	codeEditorCategory := &actions.Category{
		ID:          "code_editor",
		Name:        "Code Editors",
		Description: "Listing all code editors",
		Actions:     make([]*actions.Action, 0),
	}

	registry.RegisterCategory("dev", codeEditorCategory)

	// Register actions
	registry.RegisterAction("code_editor", InstallNeovim())
}
