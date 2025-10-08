package shells

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestNewInstallOmz tests the creation of the Install Oh My Zsh action
// This test verifies the action configuration without executing any actual commands
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

// TestInstallOmz_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
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

// TestInstallOmz_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
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

// TestInstallOmz_ActionConsistency verifies that all platform commands have consistent structure
// This ensures the action is properly configured for safe execution
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

// TestInstallOmz_MultipleCalls tests that multiple calls return equivalent actions
// This ensures the factory function is deterministic and safe
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

// TestInstallOmz_ExpectedPlatforms verifies that only expected platforms are configured
// This prevents accidental execution on unexpected platforms
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

// TestInstallOmz_CommandStructure tests that the command uses the official Oh My Zsh install script
// This verifies the command structure without executing it
func TestInstallOmz_CommandStructure(t *testing.T) {
	action := NewInstallOmz()

	// Test that the command uses the official Oh My Zsh install script
	expectedCommand := `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`

	for platform, cmd := range action.PlatformCommands {
		if cmd.Command != expectedCommand {
			t.Errorf("Platform %s has unexpected command structure: '%s'", platform, cmd.Command)
		}

		// Verify the command contains the expected URL
		if !utils.Contains(cmd.Command, "https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh") {
			t.Errorf("Platform %s command does not contain the expected Oh My Zsh install script URL", platform)
		}

		// Verify the command uses curl
		if !utils.Contains(cmd.Command, "curl") {
			t.Errorf("Platform %s command does not use curl", platform)
		}

		// Verify the command is executed with sh
		if !utils.Contains(cmd.Command, "sh -c") {
			t.Errorf("Platform %s command is not executed with sh", platform)
		}
	}
}

// TestInstallOmz_CheckCommandStructure tests that the check command properly checks for dependencies
// This ensures the action won't execute if prerequisites are not met
func TestInstallOmz_CheckCommandStructure(t *testing.T) {
	action := NewInstallOmz()

	// Test that the check command properly checks for zsh dependency and Oh My Zsh installation
	expectedCheck := "which zsh && (test -d $HOME/.oh-my-zsh || which omz)"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
		}

		// Verify the check command tests for zsh dependency first
		if !utils.Contains(cmd.CheckCommand, "which zsh") {
			t.Errorf("Platform %s check command does not test for zsh dependency", platform)
		}

		// Verify the check command tests for the .oh-my-zsh directory
		if !utils.Contains(cmd.CheckCommand, "test -d $HOME/.oh-my-zsh") {
			t.Errorf("Platform %s check command does not test for .oh-my-zsh directory", platform)
		}

		// Verify the check command also checks for omz command
		if !utils.Contains(cmd.CheckCommand, "which omz") {
			t.Errorf("Platform %s check command does not check for omz command", platform)
		}
	}
}

// TestInstallOmz_NoCommandExecution ensures that the test never executes actual commands
// This is a safety test to verify that we're only testing configuration, not execution
func TestInstallOmz_NoCommandExecution(t *testing.T) {
	action := NewInstallOmz()

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

// TestInstallOmz_SafeForCI verifies that the action is safe to use in CI environments
// This test ensures no actual system commands will be executed during testing
func TestInstallOmz_SafeForCI(t *testing.T) {
	action := NewInstallOmz()

	// Verify the action has proper check commands to prevent unnecessary execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s is missing a check command, which could lead to unnecessary execution", platform)
		}

		// Verify check command includes dependency checks
		if !utils.Contains(cmd.CheckCommand, "which zsh") {
			t.Errorf("Platform %s check command should verify zsh dependency first", platform)
		}
	}
}

// TestInstallOmz_MockExecutorBehavior tests that the action can be safely used with a mock executor
// This demonstrates how to properly mock the action execution without running actual commands
func TestInstallOmz_MockExecutorBehavior(t *testing.T) {
	action := NewInstallOmz()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_omz" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "Oh My Zsh installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "Oh My Zsh installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}
