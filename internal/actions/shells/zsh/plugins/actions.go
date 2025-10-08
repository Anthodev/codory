package shells

import (
	"anthodev/codory/internal/actions"
)

func init() {
	registry := actions.GlobalRegistry()

	zshPluginsCategory := &actions.Category{
		ID:                "zsh_plugins",
		Name:              "Zsh Plugins",
		Description:       "Install plugins for oh-my-zsh",
		Actions:           make([]*actions.Action, 0),
		HiddenOnPlatforms: []actions.Platform{actions.PlatformWindows},
	}

	registry.RegisterCategory("zsh", zshPluginsCategory)

	// Register actions
	registry.RegisterAction("zsh_plugins", InstallZshPluginHistorySearch())
	registry.RegisterAction("zsh_plugins", InstallZshPluginSyntaxHighlighting())
	registry.RegisterAction("zsh_plugins", InstallZshPluginAutosuggestions())
	registry.RegisterAction("zsh_plugins", InstallZshPluginCompletions())
}
