package shells

import (
	"anthodev/codory/internal/actions"
	"anthodev/codory/internal/platform"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func SetZshAsDefaultShell() *actions.Action {
	return &actions.Action{
		ID:                "set_zsh_default_shell",
		Name:              "Set Zsh as default shell",
		Description:       "Set Zsh as the default shell in your system",
		Type:              actions.ActionTypeFunction,
		Handler:           setZshAsDefaultShell,
		HiddenOnPlatforms: []actions.Platform{actions.PlatformWindows},
	}
}

func setZshAsDefaultShell(ctx context.Context) (string, error) {
	if err := validateZshInstallation(); err != nil {
		return "", err
	}
	zshPath, err := exec.LookPath("zsh")
	if err != nil {
		return "", fmt.Errorf("Zsh is not installed")
	}

	// Check if Zsh is already the default shell
	if shell, err := exec.Command("sh", "-c", "echo $SHELL").Output(); err != nil || strings.Contains(string(shell), "zsh") {
		return "Zsh is already the default shell", nil
	}

	// Try to set Zsh as the default shell
	if err := exec.Command("chsh", "-s", zshPath).Run(); err != nil {
		return "", fmt.Errorf("failed to set Zsh as default shell: %w", err)
	}

	return "Zsh has been set as the default shell!", nil
}

func validateZshInstallation() error {
	if platform.Detect() == platform.Windows {
		return fmt.Errorf("Zsh is not supported on Windows")
	}

	_, err := exec.LookPath("zsh")
	if err != nil {
		return fmt.Errorf("Zsh is not installed")
	}

	return nil
}
