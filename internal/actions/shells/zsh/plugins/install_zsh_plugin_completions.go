package shells

import "anthodev/codory/internal/actions"

func InstallZshPluginCompletions() *actions.Action {
	return &actions.Action{
		ID:          "install_zsh_plugin_completions",
		Name:        "Install Zsh Plugin zsh-completions",
		Description: "Installs the zsh-completions plugin for oh-my-zsh",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       "git clone https://github.com/zsh-users/zsh-completions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions",
			},
			actions.PlatformMacOS: {
				Command:       "git clone https://github.com/zsh-users/zsh-completions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions",
			},
		},
	}
}
