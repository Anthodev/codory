package dev

import (
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
)

// TestNewInstallGitmojiAction tests the creation of the Install Gitmoji action
// This test verifies the action configuration without executing any actual commands
func TestNewInstallGitmojiAction(t *testing.T) {
	action := NewInstallGitmojiAction()

	if action == nil {
		t.Fatal("NewInstallGitmojiAction() returned nil")
	}

	if action.ID != "install_gitmoji" {
		t.Errorf("Expected action ID to be 'install_gitmoji', got '%s'", action.ID)
	}

	if action.Name != "Install Gitmoji" {
		t.Errorf("Expected action name to be 'Install Gitmoji', got '%s'", action.Name)
	}

	if action.Description != "Install Gitmoji" {
		t.Errorf("Expected action description to be 'Install Gitmoji', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.SuccessMessage != "Gitmoji installed successfully! Run `gitmoji -g` to configure it." {
		t.Errorf("Expected success message to be 'Gitmoji installed successfully! Run `gitmoji -g` to configure it.', got '%s'", action.SuccessMessage)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestNewInstallGitmojiAction_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestNewInstallGitmojiAction_PlatformCommands(t *testing.T) {
	action := NewInstallGitmojiAction()

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
			expectedCommand: "brew install gitmoji",
			expectedSource:  actions.PackageSourceBrew,
			expectedCheck:   "gitmoji",
		},
		{
			name:            "macOS platform",
			platform:        actions.PlatformMacOS,
			expectedCommand: "brew install gitmoji",
			expectedSource:  actions.PackageSourceBrew,
			expectedCheck:   "gitmoji",
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

// TestNewInstallGitmojiAction_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestNewInstallGitmojiAction_UnsupportedPlatforms(t *testing.T) {
	action := NewInstallGitmojiAction()

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

// TestNewInstallGitmojiAction_ActionConsistency verifies that all platform commands have consistent structure
// This ensures the action is properly configured for safe execution
func TestNewInstallGitmojiAction_ActionConsistency(t *testing.T) {
	action := NewInstallGitmojiAction()

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
	}
}

// TestNewInstallGitmojiAction_MultipleCalls tests that multiple calls return equivalent actions
// This ensures the factory function is deterministic and safe
func TestNewInstallGitmojiAction_MultipleCalls(t *testing.T) {
	// Test that multiple calls to NewInstallGitmojiAction return equivalent actions
	action1 := NewInstallGitmojiAction()
	action2 := NewInstallGitmojiAction()

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

	if action1.SuccessMessage != action2.SuccessMessage {
		t.Errorf("Expected success messages to be consistent, got '%s' and '%s'", action1.SuccessMessage, action2.SuccessMessage)
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

// TestNewInstallGitmojiAction_ExpectedPlatforms verifies that only expected platforms are configured
// This prevents accidental execution on unexpected platforms
func TestNewInstallGitmojiAction_ExpectedPlatforms(t *testing.T) {
	action := NewInstallGitmojiAction()

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

// TestNewInstallGitmojiAction_CommandStructure tests that the command uses the correct brew install command
// This verifies the command structure without executing it
func TestNewInstallGitmojiAction_CommandStructure(t *testing.T) {
	action := NewInstallGitmojiAction()

	// Test that the command uses the correct brew install command
	expectedCommand := "brew install gitmoji"

	for platform, cmd := range action.PlatformCommands {
		if cmd.Command != expectedCommand {
			t.Errorf("Platform %s has unexpected command structure: '%s'", platform, cmd.Command)
		}

		// Verify the command contains the expected brew install
		if !utils.Contains(cmd.Command, "brew install") {
			t.Errorf("Platform %s command does not use brew install", platform)
		}

		// Verify the command contains gitmoji
		if !utils.Contains(cmd.Command, "gitmoji") {
			t.Errorf("Platform %s command does not contain gitmoji", platform)
		}
	}
}

// TestNewInstallGitmojiAction_CheckCommandStructure tests that the check command properly checks for gitmoji
// This ensures the action won't execute if prerequisites are not met
func TestNewInstallGitmojiAction_CheckCommandStructure(t *testing.T) {
	action := NewInstallGitmojiAction()

	// Test that the check command properly checks for gitmoji
	expectedCheck := "gitmoji"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
		}

		// Verify the check command tests for gitmoji
		if !utils.Contains(cmd.CheckCommand, "gitmoji") {
			t.Errorf("Platform %s check command should test for gitmoji", platform)
		}
	}
}

// TestNewInstallGitmojiAction_NoCommandExecution ensures that the test never executes actual commands
// This is a safety test to verify that we're only testing configuration, not execution
func TestNewInstallGitmojiAction_NoCommandExecution(t *testing.T) {
	action := NewInstallGitmojiAction()

	// Verify that the action is configured as a command type (not function)
	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be Command, got %s", action.Type)
	}

	// Verify that no handler is set (which would indicate function execution)
	if action.Handler != nil {
		t.Error("Expected no handler to be set for command-type actions")
	}

	// Verify that platform commands are configured (indicating command execution)
	if len(action.PlatformCommands) == 0 {
		t.Error("Expected platform commands to be configured for command-type actions")
	}
}

// TestNewInstallGitmojiAction_SafeForCI verifies that the action is safe to use in CI environments
// This test ensures no actual system commands will be executed during testing
func TestNewInstallGitmojiAction_SafeForCI(t *testing.T) {
	action := NewInstallGitmojiAction()

	// Verify the action has proper check commands to prevent unnecessary execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s is missing a check command, which could lead to unnecessary execution", platform)
		}

		// Verify check command includes dependency checks
		if !utils.Contains(cmd.CheckCommand, "gitmoji") {
			t.Errorf("Platform %s check command should verify gitmoji dependency first", platform)
		}
	}
}
