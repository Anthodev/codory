package terminal

import "anthodev/codory/internal/actions"

func InstallWarpTerminal() *actions.Action {
	return &actions.Action{
		ID:          "install_warp_terminal",
		Name:        "Install Warp Terminal",
		Description: "Install Warp Terminal in your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "yay -S warp-terminal-bin",
				PackageSource: actions.PackageSourceAUR,
				CheckCommand:  "warp",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install --cask warp",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "warp",
			},
			actions.PlatformMacOS: {
				Command:       "brew install --cask warp",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "warp",
			},
		},
	}
}
