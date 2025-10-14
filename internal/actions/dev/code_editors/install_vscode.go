package dev

import "anthodev/codory/internal/actions"

func InstallVsCode() *actions.Action {
	return &actions.Action{
		ID:          "install_vscode",
		Name:        "Install VSCode",
		Description: "Install Visual Studio Code on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S vscode",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "vscode",
				Interactive:   true,
			},
			actions.PlatformDebian: {
				Command:       "sudo apt install code",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "vscode",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install --cask vscode",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "vscode",
			},
			actions.PlatformMacOS: {
				Command:       "brew install --cask vscode",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "vscode",
				Interactive:   true,
			},
			actions.PlatformWindows: {
				Command:       "winget install -e --id Microsoft.VisualStudioCode",
				PackageSource: actions.PackageSourceWinget,
				CheckCommand:  "vscode",
				Interactive:   true,
			},
		},
	}
}
