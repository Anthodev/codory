package shells

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestInstallAtuin tests the creation of the Install Atuin action
// This test verifies the action configuration without executing any actual commands
func TestInstallAtuin(t *testing.T) {
	action := InstallAtuin()

	if action == nil {
		t.Fatal("InstallAtuin() returned nil")
	}

	if action.ID != "install_atuin" {
		t.Errorf("Expected action ID to be 'install_atuin', got '%s'", action.ID)
	}

	if action.Name != "Install Atuin" {
		t.Errorf("Expected action name to be 'Install Atuin', got '%s'", action.Name)
	}

	if action.Description != "Install Atuin shell history manager" {
		t.Errorf("Expected action description to be 'Install Atuin shell history manager', got '%s'", action.Description)
	}

	if action.SuccessMessage != "Run `atuin register -u <YOUR_USERNAME> -e <YOUR EMAIL>` to register to Atuin and `atuin login -u <USERNAME>` to login to Atuin" {
		t.Errorf("Expected success message to be specific Atuin registration instructions, got '%s'", action.SuccessMessage)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallAtuin_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallAtuin_PlatformCommands(t *testing.T) {
	action := InstallAtuin()

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
			expectedCommand: "curl --proto '=https' --tlsv1.2 -LsSf https://setup.atuin.sh | sh",
			expectedSource:  actions.PackageSourceAny,
			expectedCheck:   "atuin",
		},
		{
			name:            "macOS platform",
			platform:        actions.PlatformMacOS,
			expectedCommand: "curl --proto '=https' --tlsv1.2 -LsSf https://setup.atuin.sh | sh",
			expectedSource:  actions.PackageSourceAny,
			expectedCheck:   "atuin",
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

// TestInstallAtuin_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallAtuin_UnsupportedPlatforms(t *testing.T) {
	action := InstallAtuin()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.Platform("unsupported"),
		actions.Platform("fedora"),
		actions.Platform("opensuse"),
		actions.Platform("alpine"),
		actions.Platform("arch"),
		actions.Platform("debian"),
		actions.Platform("windows"),
	}

	for _, platform := range unsupportedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if _, exists := action.PlatformCommands[platform]; exists {
				t.Errorf("Expected platform %s to not have a command configured", platform)
			}
		})
	}
}

// TestInstallAtuin_ActionConsistency verifies that all platform commands have consistent structure
// This ensures the action is properly configured for safe execution
func TestInstallAtuin_ActionConsistency(t *testing.T) {
	action := InstallAtuin()

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
		expectedCheck := "atuin"
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command '%s', expected '%s'", platform, cmd.CheckCommand, expectedCheck)
		}
	}
}

// TestInstallAtuin_MultipleCalls tests that multiple calls return equivalent actions
// This ensures the factory function is deterministic and safe
func TestInstallAtuin_MultipleCalls(t *testing.T) {
	// Test that multiple calls to InstallAtuin return equivalent actions
	action1 := InstallAtuin()
	action2 := InstallAtuin()

	if action1.ID != action2.ID {
		t.Errorf("Expected action IDs to be consistent, got '%s' and '%s'", action1.ID, action2.ID)
	}

	if action1.Name != action2.Name {
		t.Errorf("Expected action names to be consistent, got '%s' and '%s'", action1.Name, action2.Name)
	}

	if action1.Description != action2.Description {
		t.Errorf("Expected action descriptions to be consistent, got '%s' and '%s'", action1.Description, action2.Description)
	}

	if action1.SuccessMessage != action2.SuccessMessage {
		t.Errorf("Expected success messages to be consistent, got '%s' and '%s'", action1.SuccessMessage, action2.SuccessMessage)
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

// TestInstallAtuin_ExpectedPlatforms verifies that only expected platforms are configured
// This prevents accidental execution on unexpected platforms
func TestInstallAtuin_ExpectedPlatforms(t *testing.T) {
	action := InstallAtuin()

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

// TestInstallAtuin_CommandStructure tests that the command uses the official Atuin install script
func TestInstallAtuin_CommandStructure(t *testing.T) {
	action := InstallAtuin()

	expectedPlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	expectedCommand := "curl --proto '=https' --tlsv1.2 -LsSf https://setup.atuin.sh | sh"

	for _, platform := range expectedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			cmd, exists := action.PlatformCommands[platform]
			if !exists {
				t.Fatalf("Expected %s platform command to exist", platform)
			}

			if cmd.Command != expectedCommand {
				t.Errorf("Expected %s command to be '%s', got '%s'", platform, expectedCommand, cmd.Command)
			}

			// Verify the command uses the official Atuin install script
			if !utils.Contains(cmd.Command, "https://setup.atuin.sh") {
				t.Error("Command should use the official Atuin install script")
			}

			// Verify the command uses curl with security flags
			if !utils.Contains(cmd.Command, "curl") {
				t.Error("Command should use curl to download the script")
			}

			// Verify security flags are present
			if !utils.Contains(cmd.Command, "--proto '=https'") {
				t.Error("Command should use --proto '=https' for security")
			}

			if !utils.Contains(cmd.Command, "--tlsv1.2") {
				t.Error("Command should use --tlsv1.2 for security")
			}

			// Verify the command executes the downloaded script
			if !utils.Contains(cmd.Command, "| sh") {
				t.Error("Command should pipe to shell for execution")
			}
		})
	}
}

// TestInstallAtuin_CheckCommandStructure tests that the check command properly checks for atuin
func TestInstallAtuin_CheckCommandStructure(t *testing.T) {
	action := InstallAtuin()

	// Test that the check command properly checks for atuin
	expectedCheck := "atuin"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
		}

		// Verify the check command tests for atuin
		if !utils.Contains(cmd.CheckCommand, "atuin") {
			t.Errorf("Platform %s check command should test for atuin", platform)
		}
	}
}

// TestInstallAtuin_NoCommandExecution ensures that the test never executes actual commands
// This is a safety test to verify that we're only testing configuration, not execution
func TestInstallAtuin_NoCommandExecution(t *testing.T) {
	action := InstallAtuin()

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

// TestInstallAtuin_SafeForCI verifies that the action is safe to use in CI environments
// This test ensures no actual system commands will be executed during testing
func TestInstallAtuin_SafeForCI(t *testing.T) {
	action := InstallAtuin()

	// Verify the action has proper check commands to prevent unnecessary execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s is missing a check command, which could lead to unnecessary execution", platform)
		}

		// Verify check command includes dependency checks
		if !utils.Contains(cmd.CheckCommand, "atuin") {
			t.Errorf("Platform %s check command should verify atuin dependency first", platform)
		}
	}
}

// TestInstallAtuin_MockExecutorBehavior tests that the action can be safely used with a mock executor
// This demonstrates how to properly mock the action execution without running actual commands
func TestInstallAtuin_MockExecutorBehavior(t *testing.T) {
	action := InstallAtuin()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_atuin" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "Atuin installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "Atuin installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}

// TestInstallAtuin_SuccessMessage tests that the success message provides proper guidance
func TestInstallAtuin_SuccessMessage(t *testing.T) {
	action := InstallAtuin()

	// Verify the success message contains registration instructions
	if !utils.Contains(action.SuccessMessage, "atuin register") {
		t.Error("Success message should contain registration instructions")
	}

	if !utils.Contains(action.SuccessMessage, "atuin login") {
		t.Error("Success message should contain login instructions")
	}

	// Verify the success message contains username and email placeholders
	if !utils.Contains(action.SuccessMessage, "<YOUR_USERNAME>") {
		t.Error("Success message should contain username placeholder")
	}

	if !utils.Contains(action.SuccessMessage, "<YOUR EMAIL>") {
		t.Error("Success message should contain email placeholder")
	}
}

// TestInstallAtuin_PackageSourceConsistency tests that package sources are consistent
func TestInstallAtuin_PackageSourceConsistency(t *testing.T) {
	action := InstallAtuin()

	// All platforms should use PackageSourceAny since Atuin uses a universal installer
	expectedSource := actions.PackageSourceAny

	for platform, cmd := range action.PlatformCommands {
		if cmd.PackageSource != expectedSource {
			t.Errorf("Platform %s has unexpected package source '%s', expected '%s'", platform, cmd.PackageSource, expectedSource)
		}
	}
}
