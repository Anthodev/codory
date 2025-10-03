package package_managers

import (
	"anthodev/codory/internal/actions"
)

func init() {
	registry := actions.GlobalRegistry()

	packageManagersCategory := &actions.Category{
		ID:          "package_managers",
		Name:        "Package Managers",
		Description: "Tools and utilities specific to Linux",
		Actions:     make([]*actions.Action, 0),
	}

	registry.RegisterCategory("root", packageManagersCategory)

	registry.RegisterAction("package_managers", NewInstallYayAction())
}
