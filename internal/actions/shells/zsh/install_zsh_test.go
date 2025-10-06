package shells

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestNewInstallZshAction(t *testing.T) {
	action := NewInstallZshAction()

	if action == nil {
		t.Fatal("NewInstallZshAction() returned nil")
	}

	if action.ID != "install_zsh" {
		t.Errorf("Expected action ID to be 'install_zsh', got '%s'", action.ID)
	}

	if action.Name != "Install Zsh" {
		t.Errorf("Expected action name to be 'Install Zsh', got '%s'", action.Name)
	}

	if action.Description != "Install Zsh shell" {
		t.Errorf("Expected action description to be 'Install Zsh shell', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

func TestInstallZshAction_PlatformCommands(t *testing.T) {
	action := NewInstallZshAction()

	tests := []struct {
		name            string
		platform        actions.Platform
		expectedCommand string
		expectedSource  actions.PackageSource
		expectedCheck   string
	}{
		{
			name:            "Debian platform",
			platform:        actions.PlatformDebian,
			expectedCommand: "sudo apt-get update && sudo apt-get install -y zsh",
			expectedSource:  actions.PackageSourceOfficial,
			expectedCheck:   "zsh",
		},
		{
			name:            "Arch platform",
			platform:        actions.PlatformArch,
			expectedCommand: "sudo pacman -S --noconfirm zsh",
			expectedSource:  actions.PackageSourceOfficial,
			expectedCheck:   "zsh",
		},
		{
			name:            "macOS platform",
			platform:        actions.PlatformMacOS,
			expectedCommand: "brew install zsh",
			expectedSource:  actions.PackageSourceBrew,
			expectedCheck:   "zsh",
		},
		{
			name:            "Linux platform",
			platform:        actions.PlatformLinux,
			expectedCommand: "brew install zsh",
			expectedSource:  actions.PackageSourceBrew,
			expectedCheck:   "zsh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, exists := action.PlatformCommands[tt.platform]
			if !exists {
				t.Fatalf("Expected platform command for %s to exist", tt.platform)
			}

			if cmd.Command != tt.expectedCommand {
				t.Errorf("Expected command for %s to be '%s', got '%s'", tt.platform, tt.expectedCommand, cmd.Command)
			}

			if cmd.PackageSource != tt.expectedSource {
				t.Errorf("Expected package source for %s to be '%s', got '%s'", tt.platform, tt.expectedSource, cmd.PackageSource)
			}

			if cmd.CheckCommand != tt.expectedCheck {
				t.Errorf("Expected check command for %s to be '%s', got '%s'", tt.platform, tt.expectedCheck, cmd.CheckCommand)
			}
		})
	}
}

func TestInstallZshAction_UnsupportedPlatforms(t *testing.T) {
	action := NewInstallZshAction()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.Platform("unsupported"),
		actions.Platform("windows"),
		actions.Platform("fedora"),
	}

	for _, platform := range unsupportedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if _, exists := action.PlatformCommands[platform]; exists {
				t.Errorf("Expected platform %s to not have a command configured", platform)
			}
		})
	}
}

func TestInstallZshAction_ActionConsistency(t *testing.T) {
	action := NewInstallZshAction()

	// Verify that all platform commands have consistent structure
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Platform %s has empty command", platform)
		}

		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s has empty check command", platform)
		}

		if cmd.PackageSource == "" {
			t.Errorf("Platform %s has empty package source", platform)
		}

		// Verify that check command is consistent across platforms
		if cmd.CheckCommand != "zsh" {
			t.Errorf("Platform %s has unexpected check command '%s', expected 'zsh'", platform, cmd.CheckCommand)
		}
	}
}

func TestInstallZshAction_MultipleCalls(t *testing.T) {
	// Test that multiple calls to NewInstallZshAction return equivalent actions
	action1 := NewInstallZshAction()
	action2 := NewInstallZshAction()

	if action1.ID != action2.ID {
		t.Errorf("Expected action IDs to be consistent, got '%s' and '%s'", action1.ID, action2.ID)
	}

	if action1.Name != action2.Name {
		t.Errorf("Expected action names to be consistent, got '%s' and '%s'", action1.Name, action2.Name)
	}

	if action1.Description != action2.Description {
		t.Errorf("Expected action descriptions to be consistent, got '%s' and '%s'", action1.Description, action2.Description)
	}

	if action1.Type != action2.Type {
		t.Errorf("Expected action types to be consistent, got '%s' and '%s'", action1.Type, action2.Type)
	}

	// Verify platform commands are equivalent
	if len(action1.PlatformCommands) != len(action2.PlatformCommands) {
		t.Errorf("Expected same number of platform commands, got %d and %d", len(action1.PlatformCommands), len(action2.PlatformCommands))
	}

	for platform, cmd1 := range action1.PlatformCommands {
		cmd2, exists := action2.PlatformCommands[platform]
		if !exists {
			t.Errorf("Platform %s exists in action1 but not in action2", platform)
			continue
		}

		if cmd1.Command != cmd2.Command {
			t.Errorf("Platform %s commands differ: '%s' vs '%s'", platform, cmd1.Command, cmd2.Command)
		}

		if cmd1.PackageSource != cmd2.PackageSource {
			t.Errorf("Platform %s package sources differ: '%s' vs '%s'", platform, cmd1.PackageSource, cmd2.PackageSource)
		}

		if cmd1.CheckCommand != cmd2.CheckCommand {
			t.Errorf("Platform %s check commands differ: '%s' vs '%s'", platform, cmd1.CheckCommand, cmd2.CheckCommand)
		}
	}
}

func TestInstallZshAction_LinuxPlatform(t *testing.T) {
	action := NewInstallZshAction()

	// Test Linux platform specifically
	linuxCmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	if linuxCmd.Command != "brew install zsh" {
		t.Errorf("Expected Linux command to be 'brew install zsh', got '%s'", linuxCmd.Command)
	}

	if linuxCmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected Linux package source to be '%s', got '%s'", actions.PackageSourceBrew, linuxCmd.PackageSource)
	}

	if linuxCmd.CheckCommand != "zsh" {
		t.Errorf("Expected Linux check command to be 'zsh', got '%s'", linuxCmd.CheckCommand)
	}
}

func TestInstallZshAction_ExpectedPlatforms(t *testing.T) {
	action := NewInstallZshAction()

	expectedPlatforms := []actions.Platform{
		actions.PlatformDebian,
		actions.PlatformArch,
		actions.PlatformMacOS,
		actions.PlatformLinux,
	}

	// Check that all expected platforms are present
	for _, platform := range expectedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if _, exists := action.PlatformCommands[platform]; !exists {
				t.Errorf("Expected platform %s to be present in PlatformCommands", platform)
			}
		})
	}

	// Check that we have exactly the expected number of platforms
	if len(action.PlatformCommands) != len(expectedPlatforms) {
		t.Errorf("Expected %d platforms, got %d", len(expectedPlatforms), len(action.PlatformCommands))
	}
}
