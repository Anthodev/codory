package dev

import "anthodev/codory/internal/actions"

func InstallJujutsuAction() *actions.Action {
	return &actions.Action{
		ID:          "install_jj",
		Name:        "Install Jujutsu (jj) vcs",
		Description: "Install Jujutsu (jj) version control system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S jujutsu",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "jj",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install jj",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "jj",
			},
			actions.PlatformMacOS: {
				Command:       "brew install jj",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "jj",
			},
			actions.PlatformWindows: {
				Command:       "winget install jj-vcs.jj",
				PackageSource: actions.PackageSourceWinget,
				CheckCommand:  "jj",
				Interactive:   true,
			},
		},
	}
}
