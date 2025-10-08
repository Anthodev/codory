package shells

import "anthodev/codory/internal/actions"

func InstallZshPluginSyntaxHighlighting() *actions.Action {
	return &actions.Action{
		ID:          "install_zsh_plugin_syntax_highlighting",
		Name:        "Install Zsh Plugin zsh-syntax-highlighting",
		Description: "Installs the zsh-syntax-highlighting plugin",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformLinux: {
				Command:       " git clone https://github.com/zsh-users/zsh-syntax-highlighting ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting",
			},
			actions.PlatformMacOS: {
				Command:       " git clone https://github.com/zsh-users/zsh-syntax-highlighting ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting",
				PackageSource: actions.PackageSourceAny,
				CheckCommand:  "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting",
			},
		},
	}
}
