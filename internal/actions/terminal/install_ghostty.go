package terminal

import "anthodev/codory/internal/actions"

func InstallGhostty() *actions.Action {
	return &actions.Action{
		ID:          "install_ghostty",
		Name:        "Install ghostty terminal",
		Description: "Install ghostty terminal on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S ghostty",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "ghostty",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install --cask ghostty",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "ghostty",
			},
			actions.PlatformMacOS: {
				Command:       "brew install --cask ghostty",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "ghostty",
			},
		},
	}
}
