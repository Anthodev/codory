package package_managers

import (
	"anthodev/codory/internal/actions"
	"anthodev/codory/internal/platform"
	"context"
	"fmt"
)

func NewInstallBrewAction() *actions.Action {
	return &actions.Action{
		ID:                "install_brew",
		Name:              "Install Homebrew",
		Description:       "Install Homebrew on the system",
		Type:              actions.ActionTypeFunction,
		Handler:           installBrew,
		HiddenOnPlatforms: []actions.Platform{actions.PlatformWindows},
	}
}

// InstallBrew installs Homebrew on the system.
func installBrew(ctx context.Context) (string, error) {

	if err := validateInstallBrewRequirements(); err != nil {
		return "", err
	}

	if platform.IsBrewInstalled() {
		return "Homebrew is already installed", nil
	}

	if isTestEnvironment() {
		return "", fmt.Errorf("skipping brew installation in test environment")
	}

	// Check for context cancellation before actual installation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	installer := platform.NewBrewInstaller(isTestEnvironment)

	// Check context again before calling installer
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	if err := installer.Install(ctx); err != nil {
		return "", fmt.Errorf("failed to install brew: %w", err)
	}

	return "Homebrew has been successfully installed!\n\nYou can now install packages from Homebrew using the 'brew install' command.", nil
}

func validateInstallBrewRequirements() error {
	if platform.Detect() == platform.Windows {
		return fmt.Errorf("brew can not be installed on Windows")
	}

	return nil
}
