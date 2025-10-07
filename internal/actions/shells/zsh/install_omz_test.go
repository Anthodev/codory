package shells

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestNewInstallOmz(t *testing.T) {
	action := NewInstallOmz()

	if action == nil {
		t.Fatal("NewInstallOmz() returned nil")
	}

	if action.ID != "install_omz" {
		t.Errorf("Expected action ID to be 'install_omz', got '%s'", action.ID)
	}

	if action.Name != "Install Oh My Zsh" {
		t.Errorf("Expected action name to be 'Install Oh My Zsh', got '%s'", action.Name)
	}

	if action.Description != "Install Oh My Zsh on the system" {
		t.Errorf("Expected action description to be 'Install Oh My Zsh on the system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

func TestInstallOmz_PlatformCommands(t *testing.T) {
	action := NewInstallOmz()

	tests := []struct {
		name            string
		platform        actions.Platform
		expectedCommand string
		expectedSource  actions.PackageSource
		expectedCheck   string
	}{
		{
			name:            "Linux platform",
			platform:        actions.PlatformLinux,
			expectedCommand: `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`,
			expectedSource:  actions.PackageSourceAny,
			expectedCheck:   "which zsh && (test -d $HOME/.oh-my-zsh || which omz)",
		},
		{
			name:            "macOS platform",
			platform:        actions.PlatformMacOS,
			expectedCommand: `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`,
			expectedSource:  actions.PackageSourceAny,
			expectedCheck:   "which zsh && (test -d $HOME/.oh-my-zsh || which omz)",
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

func TestInstallOmz_UnsupportedPlatforms(t *testing.T) {
	action := NewInstallOmz()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.Platform("unsupported"),
		actions.Platform("windows"),
		actions.Platform("debian"),
		actions.Platform("arch"),
	}

	for _, platform := range unsupportedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if _, exists := action.PlatformCommands[platform]; exists {
				t.Errorf("Expected platform %s to not have a command configured", platform)
			}
		})
	}
}

func TestInstallOmz_ActionConsistency(t *testing.T) {
	action := NewInstallOmz()

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
		expectedCheck := "which zsh && (test -d $HOME/.oh-my-zsh || which omz)"
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command '%s', expected '%s'", platform, cmd.CheckCommand, expectedCheck)
		}
	}
}

func TestInstallOmz_MultipleCalls(t *testing.T) {
	// Test that multiple calls to NewInstallOmz return equivalent actions
	action1 := NewInstallOmz()
	action2 := NewInstallOmz()

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

func TestInstallOmz_ExpectedPlatforms(t *testing.T) {
	action := NewInstallOmz()

	expectedPlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformMacOS,
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

func TestInstallOmz_CommandStructure(t *testing.T) {
	action := NewInstallOmz()

	// Test that the command uses the official Oh My Zsh install script
	expectedCommand := `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`

	for platform, cmd := range action.PlatformCommands {
		if cmd.Command != expectedCommand {
			t.Errorf("Platform %s has unexpected command structure: '%s'", platform, cmd.Command)
		}

		// Verify the command contains the expected URL
		if !contains(cmd.Command, "https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh") {
			t.Errorf("Platform %s command does not contain the expected Oh My Zsh install script URL", platform)
		}

		// Verify the command uses curl
		if !contains(cmd.Command, "curl") {
			t.Errorf("Platform %s command does not use curl", platform)
		}

		// Verify the command is executed with sh
		if !contains(cmd.Command, "sh -c") {
			t.Errorf("Platform %s command is not executed with sh", platform)
		}
	}
}

func TestInstallOmz_CheckCommandStructure(t *testing.T) {
	action := NewInstallOmz()

	// Test that the check command properly checks for zsh dependency and Oh My Zsh installation
	expectedCheck := "which zsh && (test -d $HOME/.oh-my-zsh || which omz)"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
		}

		// Verify the check command tests for zsh dependency first
		if !contains(cmd.CheckCommand, "which zsh") {
			t.Errorf("Platform %s check command does not test for zsh dependency", platform)
		}

		// Verify the check command tests for the .oh-my-zsh directory
		if !contains(cmd.CheckCommand, "test -d $HOME/.oh-my-zsh") {
			t.Errorf("Platform %s check command does not test for .oh-my-zsh directory", platform)
		}

		// Verify the check command also checks for omz command
		if !contains(cmd.CheckCommand, "which omz") {
			t.Errorf("Platform %s check command does not check for omz command", platform)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
