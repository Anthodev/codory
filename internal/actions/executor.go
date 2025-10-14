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
		// Try with generic platform (Linux for all Linux distros)
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
		exists := commandExists(platformCmd.CheckCommand)
		if exists {
			return "Command or files already exist, skipping installation", nil
		}
	}

	// Execute the command
	cmdStr := platformCmd.Command

	// For interactive commands, return special marker - the TUI will handle suspension
	if platformCmd.Interactive {
		return "", fmt.Errorf("INTERACTIVE_COMMAND:%s", cmdStr)
	}

	// Check if the command contains shell special characters that require shell execution
	if strings.Contains(cmdStr, "$") || strings.Contains(cmdStr, "|") ||
		strings.Contains(cmdStr, "&&") || strings.Contains(cmdStr, "||") ||
		strings.Contains(cmdStr, ";") {
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

	// For simple commands, use direct execution
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

// checkPackageManagerDependency checks if required package manager is available
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

// NeedsPackageManagerInstallation checks if an action needs a package manager to be installed
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

// GetCommandString returns the command string for an action on the current platform
func (e *Executor) GetCommandString(action *Action) (string, bool) {
	platformCmd, found := action.GetPlatformCommand(Platform(e.platformInfo.OS))
	if !found {
		if e.platformInfo.OS == platform.Debian || e.platformInfo.OS == platform.Arch {
			platformCmd, found = action.GetPlatformCommand(PlatformLinux)
		}
		if !found {
			return "", false
		}
	}
	return platformCmd.Command, true
}

// IsInteractiveCommand checks if an action requires interactive execution
func (e *Executor) IsInteractiveCommand(action *Action) bool {
	platformCmd, found := action.GetPlatformCommand(Platform(e.platformInfo.OS))
	if !found {
		if e.platformInfo.OS == platform.Debian || e.platformInfo.OS == platform.Arch {
			platformCmd, found = action.GetPlatformCommand(PlatformLinux)
		}
		if !found {
			return false
		}
	}
	return platformCmd.Interactive
}

// commandExists checks if a command or file exists
func commandExists(cmd string) bool {
	// In test mode, don't execute actual commands - just return true
	// This prevents CI failures when commands are not available
	if os.Getenv("CODORY_TEST") == "1" {
		return true
	}

	// Parse the command to extract the binary name
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return false
	}

	// First try to find the main binary in PATH
	if _, err := exec.LookPath(parts[0]); err != nil {
		// If the main binary doesn't exist, the command definitely doesn't exist
		return false
	}

	// If the command has arguments or is complex, try to execute it as a shell command
	// This handles cases like "docker compose version", "test -d /path", etc.
	if len(parts) > 1 || strings.Contains(cmd, "&&") || strings.Contains(cmd, "||") || strings.Contains(cmd, ";") {
		ctx := context.Background()
		command := exec.CommandContext(ctx, "sh", "-c", cmd)
		output, err := command.CombinedOutput()

		// Special handling for Docker in WSL environments
		if err != nil && strings.Contains(cmd, "docker") && strings.Contains(string(output), "WSL") {
			// In WSL, if docker command fails with WSL-related message,
			// we should consider it as "not properly available"
			return false
		}

		return err == nil
	}

	// For simple single commands that passed LookPath, they exist
	return true
}

// CommandExists checks if a command or file exists (public wrapper)
func (e *Executor) CommandExists(cmd string) bool {
	return commandExists(cmd)
}
