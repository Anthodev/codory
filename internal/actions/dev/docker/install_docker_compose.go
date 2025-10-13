package dev

import "anthodev/codory/internal/actions"

func NewInstallDockerComposeAction() *actions.Action {
	return &actions.Action{
		ID:          "install_docker_compose",
		Name:        "Install Docker Compose",
		Description: "Install docker compose on the system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S docker-compose",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "docker compose version",
			},
			actions.PlatformDebian: {
				Command:       "sudo apt install docker-compose",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "docker compose version",
			},
			actions.PlatformLinux: {
				Command:       "brew install docker-compose",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "docker compose version",
			},
			actions.PlatformMacOS: {
				Command:       "brew install docker-compose",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "docker compose version",
			},
		},
	}
}
