package actions

import (
	"anthodev/codory/internal/platform"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Executor executes actions
type Executor struct {
	platformInfo  platform.Info
	yayInstaller  *platform.YayInstaller
	brewInstaller *platform.BrewInstaller
}

// NewExecutor creates a new executor
func NewExecutor() *Executor {
	return &Executor{
		platformInfo:  platform.DetectInfo(),
		yayInstaller:  platform.NewYayInstaller(nil),
		brewInstaller: platform.NewBrewInstaller(nil),
	}
}

// Execute executes an action
func (e *Executor) Execute(ctx context.Context, action *Action) (string, error) {
	switch action.Type {
	case ActionTypeFunction:
		if action.Handler == nil {
			return "", fmt.Errorf("no handler defined for action %s", action.ID)
		}
		return action.Handler(ctx)

	case ActionTypeCommand:
		return e.executeCommand(ctx, action)

	default:
		return "", fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// executeCommand executes a system command
func (e *Executor) executeCommand(ctx context.Context, action *Action) (string, error) {
	platformCmd, found := action.GetPlatformCommand(Platform(e.platformInfo.OS))
	if !found {
		// Essayer avec le platform générique (Linux pour toutes les distros Linux)
		if e.platformInfo.OS == platform.Debian || e.platformInfo.OS == platform.Arch {
			platformCmd, found = action.GetPlatformCommand(PlatformLinux)
		}

		if !found {
			platformCmd, found = action.GetPlatformCommand(PlatformAny)
			if !found {
				return "", fmt.Errorf("no command defined for platform %s", e.platformInfo.OS)
			}
		}
	}

	// Check dependencies
	if err := e.checkPackageManagerDependency(ctx, platformCmd.PackageSource); err != nil {
		return "", err
	}

	// Check if the command exists already
	if platformCmd.CheckCommand != "" {
		if commandExists(platformCmd.CheckCommand) {
			return "Command or files already exist, skipping installation", nil
		}
	}

	// Execute the command
	cmdStr := platformCmd.Command

	// Check if the command contains shell special characters that require shell execution
	if strings.Contains(cmdStr, "$") || strings.Contains(cmdStr, "|") || strings.Contains(cmdStr, "&&") || strings.Contains(cmdStr, "||") || strings.Contains(cmdStr, ";") {
		// Use shell to execute complex commands
		cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return string(output), fmt.Errorf("command failed: %w\n%s", err, output)
		}
		// Use custom success message if provided, otherwise use command output
		if action.SuccessMessage != "" {
			return action.SuccessMessage, nil
		}
		return string(output), nil
	}

	// For simple commands, use the original logic
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf("command failed: %w\n%s", err, output)
	}

	// Use custom success message if provided, otherwise use command output
	if action.SuccessMessage != "" {
		return action.SuccessMessage, nil
	}
	return string(output), nil
}

func (e *Executor) checkPackageManagerDependency(ctx context.Context, source PackageSource) error {
	switch source {
	case PackageSourceAUR:
		if e.platformInfo.OS != platform.Arch {
			return fmt.Errorf("AUR packages are only available on Arch Linux")
		}
		if !platform.IsYayInstalled() {
			return fmt.Errorf("yay is not installed")
		}
	}

	return nil
}

func (e *Executor) NeedsPackageManagerInstallation(action *Action) (bool, PackageSource) {
	platformCmd, found := action.GetPlatformCommand(Platform(e.platformInfo.OS))
	if !found {
		if e.platformInfo.OS == platform.Debian || e.platformInfo.OS == platform.Arch {
			platformCmd, found = action.GetPlatformCommand(PlatformLinux)
		}
		if !found {
			return false, ""
		}
	}

	switch platformCmd.PackageSource {
	case PackageSourceAUR:
		if e.platformInfo.OS == platform.Arch && !platform.IsYayInstalled() {
			return true, PackageSourceAUR
		}
	}

	return false, ""
}

// GetPlatformInfo returns the platform information
func (e *Executor) GetPlatformInfo() platform.Info {
	return e.platformInfo
}

func commandExists(cmd string) bool {
	// In test mode, don't execute actual commands - just return true
	// This prevents CI failures when commands like 'which zsh' are not available
	if os.Getenv("CODORY_TEST") == "1" {
		return true
	}

	// First try to find it as a binary in PATH
	if _, err := exec.LookPath(cmd); err == nil {
		return true
	}

	// If that fails, try to execute it as a shell command
	// This handles cases like "test -d /path" or complex checks
	ctx := context.Background()
	command := exec.CommandContext(ctx, "sh", "-c", cmd)
	err := command.Run()
	return err == nil
}
