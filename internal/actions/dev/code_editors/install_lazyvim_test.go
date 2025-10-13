package dev

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestInstallLazyVim tests the creation of the Install LazyVim action
// This test verifies the action configuration without executing any actual commands
func TestInstallLazyVim(t *testing.T) {
	action := InstallLazyVim()

	if action == nil {
		t.Fatal("InstallLazyVim() returned nil")
	}

	if action.ID != "install_lazyvim" {
		t.Errorf("Expected action ID to be 'install_lazyvim', got '%s'", action.ID)
	}

	if action.Name != "Install LazyVim" {
		t.Errorf("Expected action name to be 'Install LazyVim', got '%s'", action.Name)
	}

	if action.Description != "Install LazyVim on your system" {
		t.Errorf("Expected action description to be 'Install LazyVim on your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallLazyVim_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallLazyVim_PlatformCommands(t *testing.T) {
	action := InstallLazyVim()

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
			expectedCommand: "git clone https://github.com/LazyVim/starter ~/.config/nvim && rm -rf ~/.config/nvim/.git",
			expectedSource:  actions.PackageSourceAny,
			expectedCheck:   "which git && which nvim && test -d ~/.config/nvim && test -f ~/.config/nvim/init.lua",
		},
		{
			name:            "macOS platform",
			platform:        actions.PlatformMacOS,
			expectedCommand: "git clone https://github.com/LazyVim/starter ~/.config/nvim && rm -rf ~/.config/nvim/.git",
			expectedSource:  actions.PackageSourceAny,
			expectedCheck:   "which git && which nvim && test -d ~/.config/nvim && test -f ~/.config/nvim/init.lua",
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

// TestInstallLazyVim_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallLazyVim_UnsupportedPlatforms(t *testing.T) {
	action := InstallLazyVim()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.Platform("unsupported"),
		actions.Platform("fedora"),
		actions.Platform("opensuse"),
		actions.Platform("alpine"),
		actions.Platform("windows"),
		actions.PlatformArch,
		actions.PlatformDebian,
	}

	for _, platform := range unsupportedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if _, exists := action.PlatformCommands[platform]; exists {
				t.Errorf("Expected platform %s to not have a command configured", platform)
			}
		})
	}
}

// TestInstallLazyVim_ActionConsistency verifies that all platform commands have consistent structure
// This ensures the action is properly configured for safe execution
func TestInstallLazyVim_ActionConsistency(t *testing.T) {
	action := InstallLazyVim()

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
		expectedCheck := "which git && which nvim && test -d ~/.config/nvim && test -f ~/.config/nvim/init.lua"
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command '%s', expected '%s'", platform, cmd.CheckCommand, expectedCheck)
		}
	}
}

// TestInstallLazyVim_MultipleCalls tests that multiple calls return equivalent actions
// This ensures the factory function is deterministic and safe
func TestInstallLazyVim_MultipleCalls(t *testing.T) {
	// Test that multiple calls to InstallLazyVim return equivalent actions
	action1 := InstallLazyVim()
	action2 := InstallLazyVim()

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

// TestInstallLazyVim_ExpectedPlatforms verifies that only expected platforms are configured
// This prevents accidental execution on unexpected platforms
func TestInstallLazyVim_ExpectedPlatforms(t *testing.T) {
	action := InstallLazyVim()

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

// TestInstallLazyVim_CommandStructure tests that the command properly clones LazyVim repository
func TestInstallLazyVim_CommandStructure(t *testing.T) {
	action := InstallLazyVim()

	for platform, cmd := range action.PlatformCommands {
		t.Run(string(platform), func(t *testing.T) {
			// Verify the command uses git clone
			if !utils.Contains(cmd.Command, "git clone") {
				t.Error("Command should use git clone to download LazyVim")
			}

			// Verify the command clones from the correct repository
			if !utils.Contains(cmd.Command, "https://github.com/LazyVim/starter") {
				t.Error("Command should clone from LazyVim starter repository")
			}

			// Verify the command removes the .git directory
			if !utils.Contains(cmd.Command, "rm -rf ~/.config/nvim/.git") {
				t.Error("Command should remove the .git directory after cloning")
			}

			// Verify it uses any package source
			if cmd.PackageSource != actions.PackageSourceAny {
				t.Errorf("Expected %s to use any package source, got '%s'", platform, cmd.PackageSource)
			}
		})
	}
}

// TestInstallLazyVim_CheckCommandStructure tests that the check command properly validates dependencies
func TestInstallLazyVim_CheckCommandStructure(t *testing.T) {
	action := InstallLazyVim()

	// Test that the check command properly validates all dependencies
	expectedCheck := "which git && which nvim && test -d ~/.config/nvim && test -f ~/.config/nvim/init.lua"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
		}

		// Verify the check command tests for git
		if !utils.Contains(cmd.CheckCommand, "which git") {
			t.Errorf("Platform %s check command should verify git dependency first", platform)
		}

		// Verify the check command tests for nvim
		if !utils.Contains(cmd.CheckCommand, "which nvim") {
			t.Errorf("Platform %s check command should verify nvim dependency", platform)
		}

		// Verify the check command tests for nvim config directory
		if !utils.Contains(cmd.CheckCommand, "test -d ~/.config/nvim") {
			t.Errorf("Platform %s check command should verify nvim config directory exists", platform)
		}

		// Verify the check command tests for init.lua file
		if !utils.Contains(cmd.CheckCommand, "test -f ~/.config/nvim/init.lua") {
			t.Errorf("Platform %s check command should verify init.lua file exists", platform)
		}
	}
}

// TestInstallLazyVim_NoCommandExecution ensures that the test never executes actual commands
// This is a safety test to verify that we're only testing configuration, not execution
func TestInstallLazyVim_NoCommandExecution(t *testing.T) {
	action := InstallLazyVim()

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

// TestInstallLazyVim_SafeForCI verifies that the action is safe to use in CI environments
// This test ensures no actual system commands will be executed during testing
func TestInstallLazyVim_SafeForCI(t *testing.T) {
	action := InstallLazyVim()

	// Verify the action has proper check commands to prevent unnecessary execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s is missing a check command, which could lead to unnecessary execution", platform)
		}

		// Verify check command includes dependency checks
		if !utils.Contains(cmd.CheckCommand, "which git") {
			t.Errorf("Platform %s check command should verify git dependency first", platform)
		}

		if !utils.Contains(cmd.CheckCommand, "which nvim") {
			t.Errorf("Platform %s check command should verify nvim dependency", platform)
		}
	}
}

// TestInstallLazyVim_MockExecutorBehavior tests that the action can be safely used with a mock executor
// This demonstrates how to properly mock the action execution without running actual commands
func TestInstallLazyVim_MockExecutorBehavior(t *testing.T) {
	action := InstallLazyVim()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_lazyvim" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "LazyVim installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "LazyVim installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}

// TestInstallLazyVim_PlatformConsistency tests that Linux and macOS both use the same commands
func TestInstallLazyVim_PlatformConsistency(t *testing.T) {
	action := InstallLazyVim()

	linuxCmd, linuxExists := action.PlatformCommands[actions.PlatformLinux]
	macOSCmd, macOSExists := action.PlatformCommands[actions.PlatformMacOS]

	if !linuxExists {
		t.Fatal("Expected Linux platform command to exist")
	}

	if !macOSExists {
		t.Fatal("Expected macOS platform command to exist")
	}

	// Both should use the same git clone command
	if linuxCmd.Command != macOSCmd.Command {
		t.Errorf("Linux and macOS commands should be identical, got '%s' and '%s'", linuxCmd.Command, macOSCmd.Command)
	}

	// Both should use the same check command
	if linuxCmd.CheckCommand != macOSCmd.CheckCommand {
		t.Errorf("Linux and macOS check commands should be identical, got '%s' and '%s'", linuxCmd.CheckCommand, macOSCmd.CheckCommand)
	}

	// Both should use any package source
	if linuxCmd.PackageSource != actions.PackageSourceAny {
		t.Errorf("Linux should use any package source, got '%s'", linuxCmd.PackageSource)
	}

	if macOSCmd.PackageSource != actions.PackageSourceAny {
		t.Errorf("macOS should use any package source, got '%s'", macOSCmd.PackageSource)
	}
}

// TestInstallLazyVim_PackageSourceValidation tests that package sources are appropriate
func TestInstallLazyVim_PackageSourceValidation(t *testing.T) {
	action := InstallLazyVim()

	// Both Linux and macOS should use any package source since LazyVim is installed via git
	for _, platform := range []actions.Platform{actions.PlatformLinux, actions.PlatformMacOS} {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected %s platform command to exist", platform)
			continue
		}

		if cmd.PackageSource != actions.PackageSourceAny {
			t.Errorf("Expected %s to use any package source, got '%s'", platform, cmd.PackageSource)
		}
	}
}
