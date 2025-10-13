package dev

import "anthodev/codory/internal/actions"

func InstallLazyVim() *actions.Action {
	return &actions.Action{
		ID:          "install_lazyvim",
		Name:        "Install LazyVim",
		Description: "Install LazyVim on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       "git clone https://github.com/LazyVim/starter ~/.config/nvim && rm -rf ~/.config/nvim/.git",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which git && which nvim && test -d ~/.config/nvim && test -f ~/.config/nvim/init.lua",
			},
			actions.PlatformMacOS: {
				Command:       "git clone https://github.com/LazyVim/starter ~/.config/nvim && rm -rf ~/.config/nvim/.git",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which git && which nvim && test -d ~/.config/nvim && test -f ~/.config/nvim/init.lua",
			},
		},
	}
}
