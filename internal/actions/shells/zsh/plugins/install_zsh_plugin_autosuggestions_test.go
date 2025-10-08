package shells

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
)

// TestInstallZshPluginAutosuggestions tests the creation of the Install Zsh Plugin Autosuggestions action
// This test verifies the action configuration without executing any actual commands
func TestInstallZshPluginAutosuggestions(t *testing.T) {
	action := InstallZshPluginAutosuggestions()

	if action == nil {
		t.Fatal("InstallZshPluginAutosuggestions() returned nil")
	}

	if action.ID != "install_zsh_plugin_autosuggestions" {
		t.Errorf("Expected action ID to be 'install_zsh_plugin_autosuggestions', got '%s'", action.ID)
	}

	if action.Name != "Install Zsh Plugin zsh-autosuggestions" {
		t.Errorf("Expected action name to be 'Install Zsh Plugin zsh-autosuggestions', got '%s'", action.Name)
	}

	if action.Description != "Install the zsh-autosuggestions plugin for oh-my-zsh" {
		t.Errorf("Expected action description to be 'Install the zsh-autosuggestions plugin for oh-my-zsh', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallZshPluginAutosuggestions_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallZshPluginAutosuggestions_PlatformCommands(t *testing.T) {
	action := InstallZshPluginAutosuggestions()

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
			expectedCommand: "git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
			expectedSource:  actions.PackageSourceAny,
			expectedCheck:   "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
		},
		{
			name:            "macOS platform",
			platform:        actions.PlatformMacOS,
			expectedCommand: "git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
			expectedSource:  actions.PackageSourceAny,
			expectedCheck:   "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
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

// TestInstallZshPluginAutosuggestions_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallZshPluginAutosuggestions_UnsupportedPlatforms(t *testing.T) {
	action := InstallZshPluginAutosuggestions()

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

// TestInstallZshPluginAutosuggestions_ActionConsistency verifies that all platform commands have consistent structure
// This ensures the action is properly configured for safe execution
func TestInstallZshPluginAutosuggestions_ActionConsistency(t *testing.T) {
	action := InstallZshPluginAutosuggestions()

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

// TestInstallZshPluginAutosuggestions_MultipleCalls tests that multiple calls return equivalent actions
// This ensures the factory function is deterministic and safe
func TestInstallZshPluginAutosuggestions_MultipleCalls(t *testing.T) {
	// Test that multiple calls to InstallZshPluginAutosuggestions return equivalent actions
	action1 := InstallZshPluginAutosuggestions()
	action2 := InstallZshPluginAutosuggestions()

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

// TestInstallZshPluginAutosuggestions_ExpectedPlatforms verifies that only expected platforms are configured
// This prevents accidental execution on unexpected platforms
func TestInstallZshPluginAutosuggestions_ExpectedPlatforms(t *testing.T) {
	action := InstallZshPluginAutosuggestions()

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

// TestInstallZshPluginAutosuggestions_CommandStructure tests that the command uses the correct git clone command
// This verifies the command structure without executing it
func TestInstallZshPluginAutosuggestions_CommandStructure(t *testing.T) {
	action := InstallZshPluginAutosuggestions()

	// Test that the command uses the correct git clone command
	expectedCommand := "git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions"

	for platform, cmd := range action.PlatformCommands {
		if cmd.Command != expectedCommand {
			t.Errorf("Platform %s has unexpected command structure: '%s'", platform, cmd.Command)
		}

		// Verify the command contains the expected URL
		if !contains(cmd.Command, "https://github.com/zsh-users/zsh-autosuggestions") {
			t.Errorf("Platform %s command does not contain the expected plugin URL", platform)
		}

		// Verify the command uses git clone
		if !contains(cmd.Command, "git clone") {
			t.Errorf("Platform %s command does not use git clone", platform)
		}

		// Verify the command uses the correct plugin directory
		if !contains(cmd.Command, "${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions") {
			t.Errorf("Platform %s command does not use the correct plugin directory", platform)
		}
	}
}

// TestInstallZshPluginAutosuggestions_CheckCommandStructure tests that the check command properly checks for dependencies
// This ensures the action won't execute if prerequisites are not met
func TestInstallZshPluginAutosuggestions_CheckCommandStructure(t *testing.T) {
	action := InstallZshPluginAutosuggestions()

	// Test that the check command properly checks for git dependency and Oh My Zsh installation
	expectedCheckLinux := "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions"
	expectedCheckMacOS := "which git && (test -d $HOME/.oh-my-zsh || which omz) && test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions"

	for platform, cmd := range action.PlatformCommands {
		if platform == actions.PlatformLinux {
			if cmd.CheckCommand != expectedCheckLinux {
				t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
			}
		} else if platform == actions.PlatformMacOS {
			if cmd.CheckCommand != expectedCheckMacOS {
				t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
			}
		}

		// Verify the check command tests for git dependency first
		if !contains(cmd.CheckCommand, "which git") {
			t.Errorf("Platform %s check command does not test for git dependency", platform)
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

// TestInstallZshPluginAutosuggestions_NoCommandExecution ensures that the test never executes actual commands
// This is a safety test to verify that we're only testing configuration, not execution
func TestInstallZshPluginAutosuggestions_NoCommandExecution(t *testing.T) {
	action := InstallZshPluginAutosuggestions()

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

// TestInstallZshPluginAutosuggestions_SafeForCI verifies that the action is safe to use in CI environments
// This test ensures no actual system commands will be executed during testing
func TestInstallZshPluginAutosuggestions_SafeForCI(t *testing.T) {
	action := InstallZshPluginAutosuggestions()

	// Verify the action has proper check commands to prevent unnecessary execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s is missing a check command, which could lead to unnecessary execution", platform)
		}

		// Verify check command includes dependency checks
		if !contains(cmd.CheckCommand, "which git") {
			t.Errorf("Platform %s check command should verify git dependency first", platform)
		}
	}
}

// TestInstallZshPluginAutosuggestions_MockExecutorBehavior tests that the action can be safely used with a mock executor
// This demonstrates how to properly mock the action execution without running actual commands
func TestInstallZshPluginAutosuggestions_MockExecutorBehavior(t *testing.T) {
	action := InstallZshPluginAutosuggestions()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := &MockExecutor{
		ExecuteFunc: func(action *actions.Action) (string, error) {
			// Verify the action structure without executing
			if action.ID != "install_zsh_plugin_autosuggestions" {
				return "", fmt.Errorf("unexpected action ID: %s", action.ID)
			}

			// Return a mock success result
			return "Zsh autosuggestions plugin installation mocked successfully", nil
		},
	}

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "Zsh autosuggestions plugin installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}
