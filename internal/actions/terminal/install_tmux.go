package terminal

import "anthodev/codory/internal/actions"

func InstallTmux() *actions.Action {
	return &actions.Action{
		ID:          "install_tmux",
		Name:        "Install Tmux",
		Description: "Install Tmux, the terminal multiplexer, on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S tmux",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "tmux",
				Interactive:   true,
			},
			actions.PlatformDebian: {
				Command:       "sudo apt-get install tmux",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "tmux",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install tmux",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "tmux",
			},
			actions.PlatformMacOS: {
				Command:       "brew install tmux",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "tmux",
			},
		},
	}
}
