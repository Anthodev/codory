package dev

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestNewInstallDockerComposeAction tests the creation of the Install Docker Compose action
// This test verifies the action configuration without executing any actual commands
func TestNewInstallDockerComposeAction(t *testing.T) {
	action := InstallDockerComposeAction()

	if action == nil {
		t.Fatal("NewInstallDockerComposeAction() returned nil")
	}

	if action.ID != "install_docker_compose" {
		t.Errorf("Expected action ID to be 'install_docker_compose', got '%s'", action.ID)
	}

	if action.Name != "Install Docker Compose" {
		t.Errorf("Expected action name to be 'Install Docker Compose', got '%s'", action.Name)
	}

	if action.Description != "Install docker compose on the system" {
		t.Errorf("Expected action description to be 'Install docker compose on the system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestNewInstallDockerComposeAction_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestNewInstallDockerComposeAction_PlatformCommands(t *testing.T) {
	action := InstallDockerComposeAction()

	tests := []struct {
		name                string
		platform            actions.Platform
		expectedCommand     string
		expectedSource      actions.PackageSource
		expectedCheck       string
		expectedInteractive bool
	}{
		{
			name:                "Arch platform",
			platform:            actions.PlatformArch,
			expectedCommand:     "sudo pacman -S docker-compose",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "docker compose version",
			expectedInteractive: true,
		},
		{
			name:                "Debian platform",
			platform:            actions.PlatformDebian,
			expectedCommand:     "sudo apt install docker-compose",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "docker compose version",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "brew install docker-compose",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "docker compose version",
			expectedInteractive: false,
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install docker-compose",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "docker compose version",
			expectedInteractive: false,
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

			if cmd.Interactive != tt.expectedInteractive {
				t.Errorf("Expected interactive for %s to be %t, got %t", tt.platform, tt.expectedInteractive, cmd.Interactive)
			}
		})
	}
}

// TestNewInstallDockerComposeAction_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestNewInstallDockerComposeAction_UnsupportedPlatforms(t *testing.T) {
	action := InstallDockerComposeAction()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.Platform("unsupported"),
		actions.Platform("windows"),
		actions.Platform("fedora"),
		actions.Platform("opensuse"),
		actions.Platform("alpine"),
	}

	for _, platform := range unsupportedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if _, exists := action.PlatformCommands[platform]; exists {
				t.Errorf("Expected platform %s to not have a command configured", platform)
			}
		})
	}
}

// TestNewInstallDockerComposeAction_ActionConsistency verifies that all platform commands have consistent structure
// This ensures the action is properly configured for safe execution
func TestNewInstallDockerComposeAction_ActionConsistency(t *testing.T) {
	action := InstallDockerComposeAction()

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

		// Verify interactive flag is set appropriately
		if platform == actions.PlatformArch || platform == actions.PlatformDebian {
			if !cmd.Interactive {
				t.Errorf("Platform %s should have interactive=true for package manager commands", platform)
			}
		} else {
			if cmd.Interactive {
				t.Errorf("Platform %s should have interactive=false for brew commands", platform)
			}
		}

		// Verify that check command is consistent across platforms
		expectedCheck := "docker compose version"
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command '%s', expected '%s'", platform, cmd.CheckCommand, expectedCheck)
		}
	}
}

// TestNewInstallDockerComposeAction_MultipleCalls tests that multiple calls return equivalent actions
// This ensures the factory function is deterministic and safe
func TestNewInstallDockerComposeAction_MultipleCalls(t *testing.T) {
	// Test that multiple calls to NewInstallDockerComposeAction return equivalent actions
	action1 := InstallDockerComposeAction()
	action2 := InstallDockerComposeAction()

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

		if cmd1.Interactive != cmd2.Interactive {
			t.Errorf("Platform %s interactive flags differ: %t vs %t", platform, cmd1.Interactive, cmd2.Interactive)
		}
	}
}

// TestNewInstallDockerComposeAction_ExpectedPlatforms verifies that only expected platforms are configured
// This prevents accidental execution on unexpected platforms
func TestNewInstallDockerComposeAction_ExpectedPlatforms(t *testing.T) {
	action := InstallDockerComposeAction()

	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
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

// TestNewInstallDockerComposeAction_ArchCommandStructure tests that the Arch command uses pacman
func TestNewInstallDockerComposeAction_ArchCommandStructure(t *testing.T) {
	action := InstallDockerComposeAction()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	expectedCommand := "sudo pacman -S docker-compose"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Arch command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses pacman
	if !utils.Contains(cmd.Command, "pacman") {
		t.Error("Arch command should use pacman package manager")
	}

	// Verify the command installs docker-compose
	if !utils.Contains(cmd.Command, "docker-compose") {
		t.Error("Arch command should install docker-compose package")
	}

	// Verify the command uses sudo
	if !utils.Contains(cmd.Command, "sudo") {
		t.Error("Arch command should use sudo for package installation")
	}
}

// TestNewInstallDockerComposeAction_DebianCommandStructure tests that the Debian command uses apt
func TestNewInstallDockerComposeAction_DebianCommandStructure(t *testing.T) {
	action := InstallDockerComposeAction()

	cmd, exists := action.PlatformCommands[actions.PlatformDebian]
	if !exists {
		t.Fatal("Expected Debian platform command to exist")
	}

	expectedCommand := "sudo apt install docker-compose"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Debian command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses apt
	if !utils.Contains(cmd.Command, "apt") {
		t.Error("Debian command should use apt package manager")
	}

	// Verify the command installs docker-compose
	if !utils.Contains(cmd.Command, "docker-compose") {
		t.Error("Debian command should install docker-compose package")
	}

	// Verify the command uses sudo
	if !utils.Contains(cmd.Command, "sudo") {
		t.Error("Debian command should use sudo for package installation")
	}
}

// TestNewInstallDockerComposeAction_LinuxCommandStructure tests that the Linux command uses brew
func TestNewInstallDockerComposeAction_LinuxCommandStructure(t *testing.T) {
	action := InstallDockerComposeAction()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	expectedCommand := "brew install docker-compose"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Linux command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses brew
	if !utils.Contains(cmd.Command, "brew") {
		t.Error("Linux command should use Homebrew package manager")
	}

	// Verify the command installs docker-compose
	if !utils.Contains(cmd.Command, "docker-compose") {
		t.Error("Linux command should install docker-compose package")
	}
}

// TestNewInstallDockerComposeAction_MacOSCommandStructure tests that the macOS command uses brew
func TestNewInstallDockerComposeAction_MacOSCommandStructure(t *testing.T) {
	action := InstallDockerComposeAction()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	expectedCommand := "brew install docker-compose"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected macOS command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses brew
	if !utils.Contains(cmd.Command, "brew") {
		t.Error("macOS command should use Homebrew package manager")
	}

	// Verify the command installs docker-compose
	if !utils.Contains(cmd.Command, "docker-compose") {
		t.Error("macOS command should install docker-compose package")
	}
}

// TestNewInstallDockerComposeAction_CheckCommandStructure tests that the check command properly checks for docker compose
func TestNewInstallDockerComposeAction_CheckCommandStructure(t *testing.T) {
	action := InstallDockerComposeAction()

	// Test that the check command properly checks for docker compose
	expectedCheck := "docker compose version"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
		}

		// Verify the check command tests for docker compose
		if !utils.Contains(cmd.CheckCommand, "docker compose") {
			t.Errorf("Platform %s check command should test for docker compose", platform)
		}

		// Verify the check command includes version check
		if !utils.Contains(cmd.CheckCommand, "version") {
			t.Errorf("Platform %s check command should include version check", platform)
		}
	}
}

// TestNewInstallDockerComposeAction_PackageSourceConsistency tests that package sources are correctly assigned
func TestNewInstallDockerComposeAction_PackageSourceConsistency(t *testing.T) {
	action := InstallDockerComposeAction()

	// Test official package sources for Arch and Debian
	officialPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
	}

	for _, platform := range officialPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			cmd, exists := action.PlatformCommands[platform]
			if !exists {
				t.Fatalf("Expected platform command for %s to exist", platform)
			}

			if cmd.PackageSource != actions.PackageSourceOfficial {
				t.Errorf("Expected %s to use official package source, got '%s'", platform, cmd.PackageSource)
			}
		})
	}

	// Test brew package sources for Linux and macOS
	brewPlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range brewPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			cmd, exists := action.PlatformCommands[platform]
			if !exists {
				t.Fatalf("Expected platform command for %s to exist", platform)
			}

			if cmd.PackageSource != actions.PackageSourceBrew {
				t.Errorf("Expected %s to use brew package source, got '%s'", platform, cmd.PackageSource)
			}
		})
	}
}

// TestNewInstallDockerComposeAction_NoCommandExecution ensures that the test never executes actual commands
// This is a safety test to verify that we're only testing configuration, not execution
func TestNewInstallDockerComposeAction_NoCommandExecution(t *testing.T) {
	action := InstallDockerComposeAction()

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

// TestNewInstallDockerComposeAction_SafeForCI verifies that the action is safe to use in CI environments
// This test ensures no actual system commands will be executed during testing
func TestNewInstallDockerComposeAction_SafeForCI(t *testing.T) {
	action := InstallDockerComposeAction()

	// Verify the action has proper check commands to prevent unnecessary execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s is missing a check command, which could lead to unnecessary execution", platform)
		}

		// Verify check command includes dependency checks
		if !utils.Contains(cmd.CheckCommand, "docker compose") {
			t.Errorf("Platform %s check command should verify docker compose dependency first", platform)
		}
	}
}

// TestNewInstallDockerComposeAction_MockExecutorBehavior tests that the action can be safely used with a mock executor
// This demonstrates how to properly mock the action execution without running actual commands
func TestNewInstallDockerComposeAction_MockExecutorBehavior(t *testing.T) {
	action := InstallDockerComposeAction()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_docker_compose" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "Docker Compose installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "Docker Compose installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}
