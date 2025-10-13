package shells

import "anthodev/codory/internal/actions"

func InstallAtuin() *actions.Action {
	return &actions.Action{
		ID:             "install_atuin",
		Name:           "Install Atuin",
		Description:    "Install Atuin shell history manager",
		SuccessMessage: "Run `atuin register -u <YOUR_USERNAME> -e <YOUR EMAIL>` to register to Atuin and `atuin login -u <USERNAME>` to login to Atuin",
		Type:           actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       "curl --proto '=https' --tlsv1.2 -LsSf https://setup.atuin.sh | sh",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "atuin",
			},
			actions.PlatformMacOS: {
				Command:       "curl --proto '=https' --tlsv1.2 -LsSf https://setup.atuin.sh | sh",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "atuin",
			},
		},
	}
}
