package dev

import "anthodev/codory/internal/actions"

func InstallNeovim() *actions.Action {
	return &actions.Action{
		ID:          "install_neovim",
		Name:        "Install Neovim",
		Description: "Install neovim on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S neovim",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "nvim",
				Interactive:   true,
			},
			actions.PlatformDebian: {
				Command:       "sudo apt install neovim",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "nvim",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install neovim",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "nvim",
			},
			actions.PlatformMacOS: {
				Command:       "brew install neovim",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "nvim",
			},
		},
	}
}
