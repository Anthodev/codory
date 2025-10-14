package dev

import "anthodev/codory/internal/actions"

func InstallHelix() *actions.Action {
	return &actions.Action{
		ID:          "install_helix",
		Name:        "Install Helix Editor",
		Description: "Install Helix Editor on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S helix",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "hx",
				Interactive:   true,
			},
			actions.PlatformDebian: {
				Command:       "sudo add-apt-repository ppa:maveonair/helix-editor && sudo apt-get update && sudo apt install helix",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "hx",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install helix",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "hx",
			},
			actions.PlatformMacOS: {
				Command:       "brew install helix",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "hx",
			},
		},
	}
}
