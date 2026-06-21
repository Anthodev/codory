package dev

import "anthodev/codory/internal/actions"

func InstallHunkAction() *actions.Action {
	return &actions.Action{
		ID:          "install_hunk",
		Name:        "Install hunk",
		Description: "Install hunk diff viewer",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       "brew install modem-dev/tap/hunk",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "hunk",
			},
			actions.PlatformMacOS: {
				Command:       "brew install modem-dev/tap/hunk",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "hunk",
			},
			actions.PlatformWindows: {
				Command:       "npm i -g hunkdiff",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "hunk",
			},
		},
	}
}
