package shells

import "anthodev/codory/internal/actions"

func NewInstallOmz() *actions.Action {
	return &actions.Action{
		ID:          "install_omz",
		Name:        "Install Oh My Zsh",
		Description: "Install Oh My Zsh on the system",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`,
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which zsh && (test -d $HOME/.oh-my-zsh || which omz)",
			},
			actions.PlatformMacOS: {
				Command:       `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`,
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which zsh && (test -d $HOME/.oh-my-zsh || which omz)",
			},
		},
	}
}
