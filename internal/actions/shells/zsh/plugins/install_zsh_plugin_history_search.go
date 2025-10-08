package shells

import "anthodev/codory/internal/actions"

func InstallZshPluginHistorySearch() *actions.Action {
	return &actions.Action{
		ID:          "install_zsh_plugin_history_search",
		Name:        "Install Zsh Plugin zsh-history-substring-search",
		Description: "Installs the zsh-history-substring-search plugin",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       "git clone https://github.com/zsh-users/zsh-history-substring-search ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-history-substring-search",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-history-substring-search",
			},
			actions.PlatformMacOS: {
				Command:       "git clone https://github.com/zsh-users/zsh-history-substring-search ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-history-substring-search",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-history-substring-search",
			},
		},
	}
}
