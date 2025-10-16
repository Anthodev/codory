package terminal

import "anthodev/codory/internal/actions"

func InstallRio() *actions.Action {
	return &actions.Action{
		ID:          "install_rio",
		Name:        "Install Rio Terminal",
		Description: "Install the Rio Terminal on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S rio",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "rio",
				Interactive:   true,
			},
			actions.PlatformDebian: {
				Command:       "flatpak install flathub com.rioterm.Rio",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "rio",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install rio-terminal",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "rio",
				Interactive:   true,
			},
			actions.PlatformMacOS: {
				Command:       "brew install rio-terminal",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "rio",
				Interactive:   true,
			},
			actions.PlatformWindows: {
				Command:       "winget.exe install --id \"raphamorim.rio\"",
				PackageSource: actions.PackageSourceWinget,
				CheckCommand:  "rio",
				Interactive:   true,
			},
		},
	}
}
