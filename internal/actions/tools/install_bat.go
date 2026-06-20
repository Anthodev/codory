package tools

import "anthodev/codory/internal/actions"

func InstallBat() *actions.Action {
	return &actions.Action{
		ID:          "install_bat",
		Name:        "Install bat",
		Description: "Install the bat content viewer tool",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S bat",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "bat",
				Interactive:   true,
			},
			actions.PlatformDebian: {
				Command:       "sudo apt install bat",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "bat",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install bat",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "bat",
			},
			actions.PlatformMacOS: {
				Command:       "brew install bat",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "bat",
			},
		},
	}
}
