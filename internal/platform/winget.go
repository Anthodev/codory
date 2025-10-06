package platform

import (
	"context"
	"fmt"
	"os/exec"
)

// WingetChecker handles winget checks
type WingetChecker struct{}

// NewWingetChecker creates a new winget checker
func NewWingetChecker() *WingetChecker {
	return &WingetChecker{}
}

// Check verifies if winget is installed and available
func (w *WingetChecker) Check(ctx context.Context) error {
	if !IsWingetInstalled() {
		return fmt.Errorf("winget is not installed")
	}

	// Check version
	cmd := exec.CommandContext(ctx, "winget", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("winget is installed but not working properly: %w\n%s", err, output)
	}

	return nil
}

// InstallInstructions returns instructions for installing winget
func (w *WingetChecker) InstallInstructions() string {
	return `Winget is not installed on your system.

To install winget:
1. Open Microsoft Store
2. Search for "App Installer"
3. Install or Update "App Installer"

Or download from: https://aka.ms/getwinget

After installation, restart your terminal and try again.`
}
