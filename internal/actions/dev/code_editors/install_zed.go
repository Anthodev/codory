package dev

import "anthodev/codory/internal/actions"

func InstallZed() *actions.Action {
	return &actions.Action{
		ID:          "install_zed",
		Name:        "Install Zed Editor",
		Description: "Install Zed Editor on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S zed",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "zed",
				Interactive:   true,
			},
			actions.PlatformDebian: {
				Command:       "flatpak install flathub dev.zed.Zed",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "zed",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "curl -f https://zed.dev/install.sh | sh",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "zed",
				Interactive:   true,
			},
			actions.PlatformMacOS: {
				Command:       "brew install --cask zed",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "zed",
			},
		},
	}
}
