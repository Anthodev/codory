package package_managers

import (
	"context"
	"fmt"

	"anthodev/codory/internal/actions"
	"anthodev/codory/internal/platform"
)

func NewCheckWingetAction() *actions.Action {
	return &actions.Action{
		ID:          "check_winget",
		Name:        "Check Winget",
		Description: "Check if Winget is installed and show installation instructions if not installed",
		Type:        actions.ActionTypeFunction,
		Handler:     checkWinget,
		// Visible only on Windows
		VisibleOnPlatforms: []actions.Platform{actions.PlatformWindows},
	}
}

func checkWinget(ctx context.Context) (string, error) {
	if err := validateCheckWinget(ctx); err != nil {
		return "", err
	}

	checker := platform.NewWingetChecker()

	if err := checker.Check(ctx); err != nil {
		return checker.InstallInstructions(), err
	}

	return "Winget is installed and working correctly!", nil
}

func validateCheckWinget(ctx context.Context) error {
	if platform.Detect() != platform.Windows {
		return fmt.Errorf("winget is not supported on this platform")
	}

	return nil
}
