package platform

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

type BrewInstaller struct {
	isTestMode func() bool
}

func NewBrewInstaller(isTestMode func() bool) *BrewInstaller {
	return &BrewInstaller{
		isTestMode: isTestMode,
	}
}

func (b *BrewInstaller) validateRequirements() error {
	if IsBrewInstalled() {
		return fmt.Errorf("brew is already installed")
	}

	if Detect() == Windows {
		return fmt.Errorf("brew cannot be installed on Windows")
	}

	if !commandExists("bash") {
		return fmt.Errorf("bash is required to install brew")
	}

	if !commandExists("curl") {
		return fmt.Errorf("curl is required to install brew")
	}

	if !commandExists("git") {
		return fmt.Errorf("git is required to install brew")
	}

	return nil
}

func (b *BrewInstaller) Install(ctx context.Context) error {
	if err := b.validateRequirements(); err != nil {
		return err
	}

	// Check test mode (only block actual installation)
	if b.isTestMode != nil && b.isTestMode() {
		return fmt.Errorf("installation blocked in test mode")
	}

	// Check for context cancellation before actual installation steps
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Check context before executing installation script
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	installScript := "/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""

	cmd := exec.CommandContext(ctx, "bash", "-c", installScript)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to install Homebrew: %w\n%s", err, output)
	}

	// Sur Linux, ajouter brew au PATH
	if runtime.GOOS == "linux" {
		return b.addBrewToPath(ctx)
	}

	return nil
}

func (b *BrewInstaller) addBrewToPath(ctx context.Context) error {
	evalCmd := `echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> ~/.bashrc`
	cmd := exec.CommandContext(ctx, "bash", "-c", evalCmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add brew to PATH: %w\n%s", err, output)
	}

	// Si zsh est installé, l'ajouter aussi
	if commandExists("zsh") {
		evalCmd = `echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> ~/.zshrc`
		cmd = exec.CommandContext(ctx, "bash", "-c", evalCmd)
		cmd.CombinedOutput()
	}

	if commandExists("fish") {
		evalCmd = `echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> ~/.config/fish/config.fish`
		cmd = exec.CommandContext(ctx, "bash", "-c", evalCmd)
		cmd.CombinedOutput()
	}

	return nil
}

func GetBrewPath() string {
	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			return "/opt/homebrew/bin/brew"
		}
		return "/usr/local/bin/brew"
	}

	return "/home/linuxbrew/.linuxbrew/bin/brew"
}
