package terminal

import "anthodev/codory/internal/actions"

func InstallZellij() *actions.Action {
	return &actions.Action{
		ID:          "install_zellij",
		Name:        "Install Zellij",
		Description: "Install Zellij, a terminal multiplexer, on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S zellij",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "zellij",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew installl zellij",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "zellij",
			},
			actions.PlatformMacOS: {
				Command:       "brew install zellij",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "zellij",
			},
		},
	}
}
