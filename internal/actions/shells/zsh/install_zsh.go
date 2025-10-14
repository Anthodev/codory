package shells

import "anthodev/codory/internal/actions"

func NewInstallZshAction() *actions.Action {
	return &actions.Action{
		ID:          "install_zsh",
		Name:        "Install Zsh",
		Description: "Install Zsh shell",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformDebian: {
				Command:       "sudo apt-get update && sudo apt-get install -y zsh",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "zsh",
				Interactive:   true,
			},
			actions.PlatformArch: {
				Command:       "sudo pacman -S --noconfirm zsh",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "zsh",
				Interactive:   true,
			},
			actions.PlatformMacOS: {
				Command:       "brew install zsh",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "zsh",
			},
			actions.PlatformLinux: {
				Command:       "brew install zsh",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "zsh",
			},
		},
	}
}
