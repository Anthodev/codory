package platform

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var testMode = false

func SetTestMode(enabled bool) {
	testMode = enabled
}

type YayInstaller struct {
	isTestMode func() bool
}

func NewYayInstaller(isTestMode func() bool) *YayInstaller {
	return &YayInstaller{
		isTestMode: isTestMode,
	}
}

func isTestEnvironment() bool {
	if testMode {
		return true
	}

	// 2. Check CODORY_TEST environment variable
	if os.Getenv("CODORY_TEST") == "1" {
		return true
	}

	// 3. Check GO_TEST environment variable
	if os.Getenv("GO_TEST") == "1" {
		return true
	}

	// 4. Check if test binary is running
	if exePath, err := os.Executable(); err == nil {
		if len(exePath) > 0 {
			return strings.Contains(exePath, ".test") || strings.Contains(exePath, "_test")
		}
	}

	return false
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

	// Prevent installation in test mode (check multiple indicators)
	if isTestEnvironment() {
		return fmt.Errorf("installation blocked in test mode")
	}

	if y.isTestMode != nil && y.isTestMode() {
		return fmt.Errorf("installation blocked in test mode")
	}

	if err := y.ensureBaseDevel(ctx); err != nil {
		return fmt.Errorf("failed to ensure base-devel: %w", err)
	}

	tmpDir := "/tmp/yay-install"

	exec.CommandContext(ctx, "rm", "-rf", tmpDir).Run()

	cloneCmd := exec.CommandContext(ctx, "git", "clone",
		"https://aur.archlinux.org/yay.git", tmpDir)
	if output, err := cloneCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone yay: %w\n%s", err, output)
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
	// Block in test mode before any sudo commands
	if isTestEnvironment() {
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
