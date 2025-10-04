package platform

import (
	"context"
	"fmt"
	"os/exec"
)

type YayInstaller struct {
	isTestMode func() bool
}

func NewYayInstaller(isTestMode func() bool) *YayInstaller {
	return &YayInstaller{
		isTestMode: isTestMode,
	}
}

func (y *YayInstaller) validateRequirements() error {
	// Vérifier si on est sur Arch Linux
	if Detect() != Arch {
		return fmt.Errorf("yay can only be installed on Arch Linux")
	}

	if !IsPacmanInstalled() {
		return fmt.Errorf("pacman is not installed")
	}

	if !commandExists("git") {
		return fmt.Errorf("git is required to install yay")
	}

	return nil
}

func (y *YayInstaller) Install(ctx context.Context) error {
	if err := y.validateRequirements(); err != nil {
		return err
	}

	// Check test mode (only block actual installation)
	if y.isTestMode != nil && y.isTestMode() {
		return fmt.Errorf("installation blocked in test mode")
	}

	// Check for context cancellation before actual installation steps
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := y.ensureBaseDevel(ctx); err != nil {
		return fmt.Errorf("failed to ensure base-devel: %w", err)
	}

	tmpDir := "/tmp/yay-install"

	exec.CommandContext(ctx, "rm", "-rf", tmpDir).Run()

	// Check context before executing commands
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	cloneCmd := exec.CommandContext(ctx, "git", "clone",
		"https://aur.archlinux.org/yay.git", tmpDir)
	if output, err := cloneCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone yay: %w\n%s", err, output)
	}

	// Check context before build step
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	makepkgCmd := exec.CommandContext(ctx, "makepkg", "-si", "--noconfirm")
	makepkgCmd.Dir = tmpDir
	if output, err := makepkgCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to build yay: %w\n%s", err, output)
	}

	// Nettoyer
	exec.CommandContext(ctx, "rm", "-rf", tmpDir).Run()

	return nil
}

func (y *YayInstaller) ensureBaseDevel(ctx context.Context) error {
	// Check test mode (only block actual installation)
	if y.isTestMode != nil && y.isTestMode() {
		return fmt.Errorf("installation blocked in test mode")
	}

	checkCmd := exec.CommandContext(ctx, "pacman", "-Qg", "base-devel")
	if err := checkCmd.Run(); err == nil {
		return nil
	}

	installCmd := exec.CommandContext(ctx, "sudo", "pacman", "-S",
		"--needed", "--noconfirm", "base-devel", "git")
	if output, err := installCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install base-devel: %w\n%s", err, output)
	}

	return nil
}
