package dev

import "anthodev/codory/internal/actions"

func InstallLazydocker() *actions.Action {
	return &actions.Action{
		ID:          "install_lazydocker",
		Name:        "Install lazydocker",
		Description: "Install lazydocker on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       "brew install lazydocker",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "lazydocker",
			},
			actions.PlatformMacOS: {
				Command:       "brew install lazydocker",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "lazydocker",
			},
			actions.PlatformWindows: {
				Command:       "winget install -e --id JesseDuffield.Lazydocker",
				PackageSource: actions.PackageSourceWinget,
				CheckCommand:  "lazydocker",
			},
		},
	}
}
