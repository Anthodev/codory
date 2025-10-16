package terminal

import "anthodev/codory/internal/actions"

func InstallWezterm() *actions.Action {
	return &actions.Action{
		ID:          "install_wezterm",
		Name:        "Install Wezterm",
		Description: "Install Wezterm, a terminal emulator, on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S wezterm",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "wezterm",
				Interactive:   true,
			},
			actions.PlatformDebian: {
				Command:       "curl -fsSL https://apt.fury.io/wez/gpg.key | sudo gpg --yes --dearmor -o /usr/share/keyrings/wezterm-fury.gpg && echo 'deb [signed-by=/usr/share/keyrings/wezterm-fury.gpg] https://apt.fury.io/wez/ * *' | sudo tee /etc/apt/sources.list.d/wezterm.list && sudo chmod 644 /usr/share/keyrings/wezterm-fury.gpg && sudo apt update && sudo apt install wezterm",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "wezterm",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install --cask wezterm",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "wezterm",
			},
			actions.PlatformMacOS: {
				Command:       "brew install --cask wezterm",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "wezterm",
			},
			actions.PlatformWindows: {
				Command:       "winget install -e --id wez.wezterm",
				PackageSource: actions.PackageSourceWinget,
				CheckCommand:  "wezterm",
				Interactive:   true,
			},
		},
	}
}
