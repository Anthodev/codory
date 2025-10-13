package shells

import (
	"anthodev/codory/internal/actions"
)

func init() {
	registry := actions.GlobalRegistry()

	shellsCategory := &actions.Category{
		ID:            "shells",
		Name:          "Shells",
		Description:   "Install and manage different shells",
		SubCategories: make([]*actions.Category, 0),
		Actions:       make([]*actions.Action, 0),
	}

	registry.RegisterCategory("root", shellsCategory)

	// Register actions
	registry.RegisterAction("shells", InstallAtuin())
}
