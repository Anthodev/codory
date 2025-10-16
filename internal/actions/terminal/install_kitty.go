package terminal

import "anthodev/codory/internal/actions"

func InstallKitty() *actions.Action {
	return &actions.Action{
		ID:          "install_kitty",
		Name:        "Install Kitty terminal",
		Description: "Install the Kitty terminal on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S kitty",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "kitty",
				Interactive:   true,
			},
			actions.PlatformDebian: {
				Command:       "sudo apt-get install kitty",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "kitty",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install --cask kitty",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "kitty",
			},
			actions.PlatformMacOS: {
				Command:       "brew install --cask kitty",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "kitty",
			},
		},
	}
}
