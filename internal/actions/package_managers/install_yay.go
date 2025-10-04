package package_managers

import (
	"context"
	"fmt"

	"anthodev/codory/internal/actions"
	"anthodev/codory/internal/platform"
)

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

	// Check for context cancellation before actual installation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	installer := platform.NewYayInstaller(isTestEnvironment)

	// Check context again before calling installer
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

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
