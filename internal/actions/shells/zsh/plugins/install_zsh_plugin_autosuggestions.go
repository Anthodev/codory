package shells

import "anthodev/codory/internal/actions"

func InstallZshPluginAutosuggestions() *actions.Action {
	return &actions.Action{
		ID:          "install_zsh_plugin_autosuggestions",
		Name:        "Install Zsh Plugin zsh-autosuggestions",
		Description: "Install the zsh-autosuggestions plugin for oh-my-zsh",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       "git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
			},
			actions.PlatformMacOS: {
				Command:       "git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
			},
		},
	}
}
