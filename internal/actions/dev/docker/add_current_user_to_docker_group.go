package dev

import "anthodev/codory/internal/actions"

func AddCurrentUserToDockerGroupAction() *actions.Action {
	return &actions.Action{
		ID:             "docker_add_user_to_group",
		Name:           "Add Current User to Docker Group",
		Description:    "Adds the current user to the docker group",
		Type:           actions.ActionTypeCommand,
		SuccessMessage: "Relaunch your shell to apply the changes",
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       "sudo usermod -aG docker $USER",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "(groups | grep docker) || which docker",
				Interactive:   true,
			},
			actions.PlatformMacOS: {
				Command:       "sudo dseditgroup -o edit -a $USER -t user docker",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "(groups | grep docker) || which docker",
				Interactive:   true,
			},
		},
	}
}
