package tools

import "anthodev/codory/internal/actions"

func InstallBtop() *actions.Action {
	return &actions.Action{
		ID:          "install_btop",
		Name:        "Install btop",
		Description: "Install btop system monitor",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S btop",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "btop",
				Interactive:   true,
			},
			actions.PlatformDebian: {
				Command:       "sudo apt install btop",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "btop",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install btop",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "btop",
			},
			actions.PlatformMacOS: {
				Command:       "brew install btop",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "btop",
			},
		},
	}
}
