package package_managers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"anthodev/codory/internal/actions"
	"anthodev/codory/internal/platform"
)

// testMode is a package-level variable that can be set by tests to prevent actual installation
var testMode = false

// SetTestMode enables or disables test mode (should only be called by tests)
func SetTestMode(enabled bool) {
	testMode = enabled
}

func NewInstallYayAction() *actions.Action {
	return &actions.Action{
		ID:                 "install_yay",
		Name:               "Install Yay (AUR Helper)",
		Description:        "Install Yay AUR helper for Arch Linux",
		Type:               actions.ActionTypeFunction,
		Handler:            installYay,
		VisibleOnPlatforms: []actions.Platform{actions.PlatformArch},
	}
}

func installYay(ctx context.Context) (string, error) {
	// Validate platform requirements first
	if err := validateInstallYayRequirements(); err != nil {
		return "", err
	}

	// Check if yay is already installed
	if platform.IsYayInstalled() {
		return "Yay is already installed!", nil
	}

	// Check if we're in test mode (multiple detection methods)
	if isTestEnvironment() {
		return "", fmt.Errorf("skipping yay installation in test environment")
	}

	installer := platform.NewYayInstaller(isTestEnvironment)

	if err := installer.Install(ctx); err != nil {
		return "", fmt.Errorf("failed to install yay: %w", err)
	}

	return "Yay has been successfully installed!\n\nYou can now install packages from AUR using yay.", nil
}

// validateInstallYayRequirements validates that yay can be installed on the current system
func validateInstallYayRequirements() error {
	if platform.Detect() != platform.Arch {
		return fmt.Errorf("yay can only be installed on Arch Linux")
	}
	return nil
}

// isTestEnvironment checks if we're running in a test environment
func isTestEnvironment() bool {
	// 1. Check package-level test mode flag (set by tests)
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

	// 4. Check if test binary is running (test binaries have .test suffix or contain .test. in name)
	if exePath, err := os.Executable(); err == nil {
		if len(exePath) > 0 {
			// Check if it's a test binary
			return strings.Contains(exePath, ".test") || strings.Contains(exePath, "_test")
		}
	}

	return false
}
