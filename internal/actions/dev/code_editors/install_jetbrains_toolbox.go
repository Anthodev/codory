package dev

import "anthodev/codory/internal/actions"

func InstallJetBrainsToolbox() *actions.Action {
	return &actions.Action{
		ID:          "install_jetbrains_toolbox",
		Name:        "Install Jetbrains Toolbox",
		Description: "Install Jetbrains Toolbox on your system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "yay -S jetbrains-toolbox",
				PackageSource: actions.PackageSourceAUR,
				CheckCommand:  "jetbrains-toolbox",
				Interactive:   true,
			},
			actions.PlatformLinux: {
				Command:       "echo Follow the instructions on the following page: https://www.jetbrains.com/help/toolbox-app/toolbox-app-silent-installation.html#tba_installation",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "jetbrains-toolbox",
			},
			actions.PlatformMacOS: {
				Command:       "echo Follow the instructions on the following page: https://www.jetbrains.com/help/toolbox-app/toolbox-app-silent-installation.html#toolbox_macOS",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "jetbrains-toolbox",
			},
			actions.PlatformWindows: {
				Command:       "winget install -e --id JetBrains.Toolbox",
				PackageSource: actions.PackageSourceWinget,
				CheckCommand:  "jetbrains-toolbox",
				Interactive:   true,
			},
		},
	}
}
