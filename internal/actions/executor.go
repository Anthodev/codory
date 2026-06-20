package actions

import (
	"anthodev/codory/internal/platform"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type packageManagerInstaller interface {
	Install(context.Context) error
}

// Executor executes actions
type Executor struct {
	platformInfo  platform.Info
	yayInstaller  packageManagerInstaller
	brewInstaller packageManagerInstaller
}

// NewExecutor creates a new executor
func NewExecutor() *Executor {
	return &Executor{
		platformInfo:  platform.DetectInfo(),
		yayInstaller:  platform.NewYayInstaller(isTestEnvironment),
		brewInstaller: platform.NewBrewInstaller(isTestEnvironment),
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
	platformCmd, found := e.resolvePlatformCommand(action)
	if !found {
		return "", fmt.Errorf("no command defined for platform %s", e.platformInfo.OS)
	}

	// Check if the command exists already
	if platformCmd.CheckCommand != "" {
		exists := commandExists(platformCmd.CheckCommand)
		if exists {
			return "Command or files already exist, skipping installation", nil
		}
	}

	// Check dependencies
	if err := e.checkPackageManagerDependency(ctx, platformCmd.PackageSource); err != nil {
		return "", err
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
		cmd := NewShellCommand(ctx, cmdStr)
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
	case PackageSourceBrew:
		if e.platformInfo.OS == platform.Windows {
			return fmt.Errorf("brew packages are not available on Windows")
		}
		if !platform.IsBrewInstalled() {
			return fmt.Errorf("brew is not installed")
		}
	case PackageSourceWinget:
		if e.platformInfo.OS != platform.Windows {
			return fmt.Errorf("winget packages are only available on Windows")
		}
		if !platform.IsWingetInstalled() {
			return fmt.Errorf("winget is not installed")
		}
	}

	return nil
}

// InstallPackageManager installs the package manager needed by an action.
func (e *Executor) InstallPackageManager(ctx context.Context, source PackageSource) (string, error) {
	switch source {
	case PackageSourceAUR:
		if e.platformInfo.OS != platform.Arch {
			return "", fmt.Errorf("AUR packages are only available on Arch Linux")
		}
		if platform.IsYayInstalled() {
			return "Yay is already installed", nil
		}
		if e.yayInstaller == nil {
			e.yayInstaller = platform.NewYayInstaller(isTestEnvironment)
		}
		if err := e.yayInstaller.Install(ctx); err != nil {
			return "", fmt.Errorf("failed to install yay: %w", err)
		}
		return "Yay has been successfully installed", nil

	case PackageSourceBrew:
		if e.platformInfo.OS == platform.Windows {
			return "", fmt.Errorf("brew cannot be installed on Windows")
		}
		if platform.IsBrewInstalled() {
			return "Homebrew is already installed", nil
		}
		if e.brewInstaller == nil {
			e.brewInstaller = platform.NewBrewInstaller(isTestEnvironment)
		}
		if err := e.brewInstaller.Install(ctx); err != nil {
			return "", fmt.Errorf("failed to install brew: %w", err)
		}
		if err := prependToPathOnce(filepath.Dir(platform.GetBrewPath())); err != nil {
			return "", fmt.Errorf("failed to update PATH for brew: %w", err)
		}
		if !platform.IsBrewInstalled() {
			return "", fmt.Errorf("brew installed but not found in current PATH")
		}
		return "Homebrew has been successfully installed", nil

	case PackageSourceWinget:
		if e.platformInfo.OS != platform.Windows {
			return "", fmt.Errorf("winget packages are only available on Windows")
		}
		if platform.IsWingetInstalled() {
			return "Winget is already installed", nil
		}
		return "", errors.New(platform.NewWingetChecker().InstallInstructions())
	}

	return "", fmt.Errorf("unsupported package manager source: %s", source)
}

// NeedsPackageManagerInstallation checks if an action needs a package manager to be installed
func (e *Executor) NeedsPackageManagerInstallation(action *Action) (bool, PackageSource) {
	platformCmd, found := e.resolvePlatformCommand(action)
	if !found {
		return false, ""
	}
	if platformCmd.CheckCommand != "" && commandExists(platformCmd.CheckCommand) {
		return false, ""
	}

	switch platformCmd.PackageSource {
	case PackageSourceAUR:
		if e.platformInfo.OS == platform.Arch && !platform.IsYayInstalled() {
			return true, PackageSourceAUR
		}
	case PackageSourceBrew:
		if e.platformInfo.OS != platform.Windows && !platform.IsBrewInstalled() {
			return true, PackageSourceBrew
		}
	case PackageSourceWinget:
		if e.platformInfo.OS == platform.Windows && !platform.IsWingetInstalled() {
			return true, PackageSourceWinget
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
	platformCmd, found := e.resolvePlatformCommand(action)
	if !found {
		return "", false
	}
	return platformCmd.Command, true
}

// IsInteractiveCommand checks if an action requires interactive execution
func (e *Executor) IsInteractiveCommand(action *Action) bool {
	platformCmd, found := e.resolvePlatformCommand(action)
	if !found {
		return false
	}
	return platformCmd.Interactive
}

func (e *Executor) resolvePlatformCommand(action *Action) (PlatformCommand, bool) {
	return action.ResolvePlatformCommand(Platform(e.platformInfo.OS))
}

// NewShellCommand creates the platform shell command used for compound commands.
func NewShellCommand(ctx context.Context, cmdStr string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/C", cmdStr)
	}
	return exec.CommandContext(ctx, "sh", "-c", cmdStr)
}

func isTestEnvironment() bool {
	return os.Getenv("CODORY_TEST") == "1" || os.Getenv("GO_TEST") == "1"
}

func prependToPathOnce(dir string) error {
	current := os.Getenv("PATH")
	for _, entry := range filepath.SplitList(current) {
		if entry == dir {
			return nil
		}
	}
	if current == "" {
		return os.Setenv("PATH", dir)
	}
	return os.Setenv("PATH", dir+string(os.PathListSeparator)+current)
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
		command := NewShellCommand(ctx, cmd)
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
