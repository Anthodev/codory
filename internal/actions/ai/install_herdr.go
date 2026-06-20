package ai

import "anthodev/codory/internal/actions"

func InstallHerdr() *actions.Action {
	return &actions.Action{
		ID:          "install_herdr",
		Name:        "Install Herdr",
		Description: "Install Herdr agent multiplexer",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "yay -S herdr-bin",
				PackageSource: actions.PackageSourceAUR,
				CheckCommand:  "herdr",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "brew install herdr",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "herdr",
			},
			actions.PlatformMacOS: {
				Command:       "brew install herdr",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "herdr",
			},
			actions.PlatformWindows: {
				Command:       "powershell -ExecutionPolicy Bypass -c \"irm https://herdr.dev/install.ps1 | iex\"",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "herdr",
				Interactive:   true,
			},
		},
	}
}
