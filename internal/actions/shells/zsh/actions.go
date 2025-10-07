package shells

import (
	"anthodev/codory/internal/actions"
)

func init() {
	registry := actions.GlobalRegistry()

	zshCategory := &actions.Category{
		ID:                "zsh",
		Name:              "Zsh",
		Description:       "Install zsh and most common plugins",
		Actions:           make([]*actions.Action, 0),
		HiddenOnPlatforms: []actions.Platform{actions.PlatformWindows},
	}

	registry.RegisterCategory("shells", zshCategory)

	// Register actions
	registry.RegisterAction("zsh", NewInstallZshAction())
	registry.RegisterAction("zsh", SetZshAsDefaultShell())
	registry.RegisterAction("zsh", NewInstallOmz())
}
