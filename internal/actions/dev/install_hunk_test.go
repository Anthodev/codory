package dev

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestInstallHunkAction(t *testing.T) {
	action := InstallHunkAction()

	if action.ID != "install_hunk" {
		t.Errorf("Expected action ID to be 'install_hunk', got '%s'", action.ID)
	}
	if action.Name != "Install hunk" {
		t.Errorf("Expected action name to be 'Install hunk', got '%s'", action.Name)
	}
	if action.Description != "Install hunk diff viewer" {
		t.Errorf("Expected action description to be 'Install hunk diff viewer', got '%s'", action.Description)
	}
	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be command, got %s", action.Type)
	}

	for _, platform := range []actions.Platform{actions.PlatformLinux, actions.PlatformMacOS, actions.PlatformWindows} {
		cmd, ok := action.GetPlatformCommand(platform)
		if !ok {
			t.Fatalf("Expected command for %s", platform)
		}
		if cmd.CheckCommand != "hunk" {
			t.Errorf("Expected check command hunk for %s, got %s", platform, cmd.CheckCommand)
		}
	}

	for _, platform := range []actions.Platform{actions.PlatformLinux, actions.PlatformMacOS} {
		cmd, _ := action.GetPlatformCommand(platform)
		if cmd.Command != "brew install modem-dev/tap/hunk" {
			t.Errorf("Expected brew install for %s, got %s", platform, cmd.Command)
		}
		if cmd.PackageSource != actions.PackageSourceBrew {
			t.Errorf("Expected brew package source for %s, got %s", platform, cmd.PackageSource)
		}
	}

	cmd, _ := action.GetPlatformCommand(actions.PlatformWindows)
	if cmd.Command != "npm i -g hunkdiff" {
		t.Errorf("Expected npm install for windows, got %s", cmd.Command)
	}
	if cmd.PackageSource != actions.PackageSourceAny {
		t.Errorf("Expected any package source for windows, got %s", cmd.PackageSource)
	}
}
