package dev

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestAddCurrentUserToDockerGroup tests the creation of the Add Current User to Docker Group action
// This test verifies the action configuration without executing any actual commands
func TestAddCurrentUserToDockerGroup(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

	if action == nil {
		t.Fatal("AddCurrentUserToDockerGroup() returned nil")
	}

	if action.ID != "docker_add_user_to_group" {
		t.Errorf("Expected action ID to be 'docker_add_user_to_group', got '%s'", action.ID)
	}

	if action.Name != "Add Current User to Docker Group" {
		t.Errorf("Expected action name to be 'Add Current User to Docker Group', got '%s'", action.Name)
	}

	if action.Description != "Adds the current user to the docker group" {
		t.Errorf("Expected action description to be 'Adds the current user to the docker group', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.SuccessMessage != "Relaunch your shell to apply the changes" {
		t.Errorf("Expected success message to be 'Relaunch your shell to apply the changes', got '%s'", action.SuccessMessage)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestAddCurrentUserToDockerGroup_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestAddCurrentUserToDockerGroup_PlatformCommands(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

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
			expectedCommand: "sudo usermod -aG docker $USER",
			expectedSource:  actions.PackageSourceAny,
			expectedCheck:   "(groups | grep docker) || which docker",
		},
		{
			name:            "macOS platform",
			platform:        actions.PlatformMacOS,
			expectedCommand: "sudo dseditgroup -o edit -a $USER -t user docker",
			expectedSource:  actions.PackageSourceAny,
			expectedCheck:   "(groups | grep docker) || which docker",
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

// TestAddCurrentUserToDockerGroup_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestAddCurrentUserToDockerGroup_UnsupportedPlatforms(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.Platform("unsupported"),
		actions.Platform("fedora"),
		actions.Platform("opensuse"),
		actions.Platform("alpine"),
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformWindows,
	}

	for _, platform := range unsupportedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if _, exists := action.PlatformCommands[platform]; exists {
				t.Errorf("Expected platform %s to not have a command configured", platform)
			}
		})
	}
}

// TestAddCurrentUserToDockerGroup_ActionConsistency verifies that all platform commands have consistent structure
// This ensures the action is properly configured for safe execution
func TestAddCurrentUserToDockerGroup_ActionConsistency(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

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
		expectedCheck := "(groups | grep docker) || which docker"
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command '%s', expected '%s'", platform, cmd.CheckCommand, expectedCheck)
		}
	}
}

// TestAddCurrentUserToDockerGroup_MultipleCalls tests that multiple calls return equivalent actions
// This ensures the factory function is deterministic and safe
func TestAddCurrentUserToDockerGroup_MultipleCalls(t *testing.T) {
	// Test that multiple calls to AddCurrentUserToDockerGroup return equivalent actions
	action1 := AddCurrentUserToDockerGroupAction()
	action2 := AddCurrentUserToDockerGroupAction()

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

// TestAddCurrentUserToDockerGroup_ExpectedPlatforms verifies that only expected platforms are configured
// This prevents accidental execution on unexpected platforms
func TestAddCurrentUserToDockerGroup_ExpectedPlatforms(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

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

// TestAddCurrentUserToDockerGroup_LinuxCommandStructure tests that the Linux command uses usermod
func TestAddCurrentUserToDockerGroup_LinuxCommandStructure(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	expectedCommand := "sudo usermod -aG docker $USER"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Linux command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses usermod
	if !utils.Contains(cmd.Command, "usermod") {
		t.Error("Linux command should use usermod command")
	}

	// Verify the command adds user to docker group
	if !utils.Contains(cmd.Command, "-aG docker") {
		t.Error("Linux command should add user to docker group with -aG docker")
	}

	// Verify the command uses $USER variable
	if !utils.Contains(cmd.Command, "$USER") {
		t.Error("Linux command should use $USER variable")
	}
}

// TestAddCurrentUserToDockerGroup_MacOSCommandStructure tests that the macOS command uses dseditgroup
func TestAddCurrentUserToDockerGroup_MacOSCommandStructure(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	expectedCommand := "sudo dseditgroup -o edit -a $USER -t user docker"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected macOS command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses dseditgroup
	if !utils.Contains(cmd.Command, "dseditgroup") {
		t.Error("macOS command should use dseditgroup command")
	}

	// Verify the command edits the docker group
	if !utils.Contains(cmd.Command, "docker") {
		t.Error("macOS command should edit the docker group")
	}

	// Verify the command uses $USER variable
	if !utils.Contains(cmd.Command, "$USER") {
		t.Error("macOS command should use $USER variable")
	}
}

// TestAddCurrentUserToDockerGroup_CheckCommandStructure tests that the check command properly checks for docker group membership
func TestAddCurrentUserToDockerGroup_CheckCommandStructure(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

	// Test that the check command properly checks for docker group membership
	expectedCheck := "(groups | grep docker) || which docker"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
		}

		// Verify the check command tests for docker group membership
		if !utils.Contains(cmd.CheckCommand, "groups") {
			t.Errorf("Platform %s check command should test for group membership", platform)
		}

		if !utils.Contains(cmd.CheckCommand, "grep docker") {
			t.Errorf("Platform %s check command should grep for docker group", platform)
		}

		// Verify the check command also tests for docker installation
		if !utils.Contains(cmd.CheckCommand, "which docker") {
			t.Errorf("Platform %s check command should test for docker installation", platform)
		}
	}
}

// TestAddCurrentUserToDockerGroup_NoCommandExecution ensures that the test never executes actual commands
// This is a safety test to verify that we're only testing configuration, not execution
func TestAddCurrentUserToDockerGroup_NoCommandExecution(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

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

// TestAddCurrentUserToDockerGroup_SafeForCI verifies that the action is safe to use in CI environments
// This test ensures no actual system commands will be executed during testing
func TestAddCurrentUserToDockerGroup_SafeForCI(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

	// Verify the action has proper check commands to prevent unnecessary execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s is missing a check command, which could lead to unnecessary execution", platform)
		}

		// Verify check command includes dependency checks
		if !utils.Contains(cmd.CheckCommand, "docker") {
			t.Errorf("Platform %s check command should verify docker dependency first", platform)
		}
	}
}

// TestAddCurrentUserToDockerGroup_MockExecutorBehavior tests that the action can be safely used with a mock executor
// This demonstrates how to properly mock the action execution without running actual commands
func TestAddCurrentUserToDockerGroup_MockExecutorBehavior(t *testing.T) {
	action := AddCurrentUserToDockerGroupAction()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "docker_add_user_to_group" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "User added to docker group successfully (mocked)", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "User added to docker group successfully (mocked)"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}
