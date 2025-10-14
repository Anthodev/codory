package dev

import "anthodev/codory/internal/actions"

func InstallLazygit() *actions.Action {
	return &actions.Action{
		ID:          "install_lazygit",
		Name:        "Install Lazygit",
		Description: "Installs Lazygit",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       "brew install lazygit",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "lazygit",
			},
			actions.PlatformMacOS: {
				Command:       "brew install lazygit",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "lazygit",
			},
			actions.PlatformWindows: {
				Command:       "winget install -e --id JesseDuffield.lazygit",
				PackageSource: actions.PackageSourceWinget,
				CheckCommand:  "lazygit",
			},
		},
	}
}
