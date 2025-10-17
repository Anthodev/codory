package dev

import "anthodev/codory/internal/actions"

func init() {
	registry := actions.GlobalRegistry()

	devCategory := &actions.Category{
		ID:            "dev",
		Name:          "Development",
		Description:   "Development tools and utilities",
		SubCategories: make([]*actions.Category, 0),
		Actions:       make([]*actions.Action, 0),
	}

	registry.RegisterCategory("root", devCategory)

	// Register actions
	registry.RegisterAction("dev", NewUUIDv4Action())
	registry.RegisterAction("dev", NewUUIDv7Action())
	registry.RegisterAction("dev", DecodeUUIDv7Action())
	registry.RegisterAction("dev", NewSymfonySecretAction())
	registry.RegisterAction("dev", InstallLazygitAction())
	registry.RegisterAction("dev", InstallGitmojiAction())
	registry.RegisterAction("dev", InstallJujutsuAction())
}
