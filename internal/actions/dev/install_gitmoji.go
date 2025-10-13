package dev

import "anthodev/codory/internal/actions"

func InstallGitmojiAction() *actions.Action {
	return &actions.Action{
		ID:             "install_gitmoji",
		Name:           "Install Gitmoji",
		Description:    "Install Gitmoji",
		Type:           actions.ActionTypeCommand,
		SuccessMessage: "Gitmoji installed successfully! Run `gitmoji -g` to configure it.",
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       "brew install gitmoji",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "gitmoji",
			},
			actions.PlatformMacOS: {
				Command:       "brew install gitmoji",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "gitmoji",
			},
		},
	}
}
