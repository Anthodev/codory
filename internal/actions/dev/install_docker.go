package dev

import "anthodev/codory/internal/actions"

func NewInstallDockerAction() *actions.Action {
	return &actions.Action{
		ID:          "install_docker",
		Name:        "Install Docker engine",
		Description: "Install Docker engine",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S docker",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "docker",
			},
			actions.PlatformDebian: {
				Command:       "sudo apt get install docker",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "docker",
			},
			actions.PlatformLinux: {
				Command:       "curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "docker",
			},
			actions.PlatformMacOS: {
				Command:       "brew install docker",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "docker",
			},
			actions.PlatformWindows: {
				Command:       "winget install Docker.DockerDesktop",
				PackageSource: actions.PackageSourceWinget,
				CheckCommand:  "docker",
			},
		},
	}
}
