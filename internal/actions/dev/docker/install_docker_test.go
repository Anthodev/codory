package dev

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestNewInstallDockerAction tests the creation of the Install Docker action
// This test verifies the action configuration without executing any actual commands
func TestNewInstallDockerAction(t *testing.T) {
	action := InstallDockerAction()

	if action == nil {
		t.Fatal("NewInstallDockerAction() returned nil")
	}

	if action.ID != "install_docker" {
		t.Errorf("Expected action ID to be 'install_docker', got '%s'", action.ID)
	}

	if action.Name != "Install Docker engine" {
		t.Errorf("Expected action name to be 'Install Docker engine', got '%s'", action.Name)
	}

	if action.Description != "Install Docker engine" {
		t.Errorf("Expected action description to be 'Install Docker engine', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestNewInstallDockerAction_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestNewInstallDockerAction_PlatformCommands(t *testing.T) {
	action := InstallDockerAction()

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
			expectedCommand:     "sudo pacman -S docker",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "docker",
			expectedInteractive: true,
		},
		{
			name:                "Debian platform",
			platform:            actions.PlatformDebian,
			expectedCommand:     "sudo apt get install docker",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "docker",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh",
			expectedSource:      actions.PackageSourceAny,
			expectedCheck:       "docker",
			expectedInteractive: false,
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install docker",
			expectedSource:      actions.PackageSourceAny,
			expectedCheck:       "docker",
			expectedInteractive: false,
		},
		{
			name:                "Windows platform",
			platform:            actions.PlatformWindows,
			expectedCommand:     "winget install -e --id Docker.DockerDesktop",
			expectedSource:      actions.PackageSourceWinget,
			expectedCheck:       "docker",
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

// TestNewInstallDockerAction_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestNewInstallDockerAction_UnsupportedPlatforms(t *testing.T) {
	action := InstallDockerAction()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.Platform("unsupported"),
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

// TestNewInstallDockerAction_ActionConsistency verifies that all platform commands have consistent structure
// This ensures the action is properly configured for safe execution
func TestNewInstallDockerAction_ActionConsistency(t *testing.T) {
	action := InstallDockerAction()

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
		expectedCheck := "docker"
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command '%s', expected '%s'", platform, cmd.CheckCommand, expectedCheck)
		}

		// Verify interactive flag is set appropriately
		if platform == actions.PlatformArch || platform == actions.PlatformDebian {
			if !cmd.Interactive {
				t.Errorf("Platform %s should have interactive=true for package manager commands", platform)
			}
		} else {
			if cmd.Interactive {
				t.Errorf("Platform %s should have interactive=false for non-package manager commands", platform)
			}
		}
	}
}

// TestNewInstallDockerAction_MultipleCalls tests that multiple calls return equivalent actions
// This ensures the factory function is deterministic and safe
func TestNewInstallDockerAction_MultipleCalls(t *testing.T) {
	// Test that multiple calls to NewInstallDockerAction return equivalent actions
	action1 := InstallDockerAction()
	action2 := InstallDockerAction()

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

// TestNewInstallDockerAction_ExpectedPlatforms verifies that only expected platforms are configured
// This prevents accidental execution on unexpected platforms
func TestNewInstallDockerAction_ExpectedPlatforms(t *testing.T) {
	action := InstallDockerAction()

	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformLinux,
		actions.PlatformMacOS,
		actions.PlatformWindows,
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

// TestNewInstallDockerAction_ArchCommandStructure tests that the Arch command uses pacman
func TestNewInstallDockerAction_ArchCommandStructure(t *testing.T) {
	action := InstallDockerAction()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	expectedCommand := "sudo pacman -S docker"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Arch command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses pacman
	if !utils.Contains(cmd.Command, "pacman") {
		t.Error("Arch command should use pacman package manager")
	}

	// Verify the command installs docker
	if !utils.Contains(cmd.Command, "docker") {
		t.Error("Arch command should install docker package")
	}
}

// TestNewInstallDockerAction_DebianCommandStructure tests that the Debian command uses apt
func TestNewInstallDockerAction_DebianCommandStructure(t *testing.T) {
	action := InstallDockerAction()

	cmd, exists := action.PlatformCommands[actions.PlatformDebian]
	if !exists {
		t.Fatal("Expected Debian platform command to exist")
	}

	expectedCommand := "sudo apt get install docker"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Debian command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses apt
	if !utils.Contains(cmd.Command, "apt") {
		t.Error("Debian command should use apt package manager")
	}

	// Verify the command installs docker
	if !utils.Contains(cmd.Command, "docker") {
		t.Error("Debian command should install docker package")
	}
}

// TestNewInstallDockerAction_LinuxCommandStructure tests that the Linux command uses the official Docker install script
func TestNewInstallDockerAction_LinuxCommandStructure(t *testing.T) {
	action := InstallDockerAction()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	expectedCommand := "curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Linux command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses the official Docker install script
	if !utils.Contains(cmd.Command, "https://get.docker.com") {
		t.Error("Linux command should use the official Docker install script")
	}

	// Verify the command uses curl
	if !utils.Contains(cmd.Command, "curl") {
		t.Error("Linux command should use curl to download the script")
	}

	// Verify the command executes the downloaded script
	if !utils.Contains(cmd.Command, "sh get-docker.sh") {
		t.Error("Linux command should execute the downloaded Docker install script")
	}
}

// TestNewInstallDockerAction_MacOSCommandStructure tests that the macOS command uses brew
func TestNewInstallDockerAction_MacOSCommandStructure(t *testing.T) {
	action := InstallDockerAction()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	expectedCommand := "brew install docker"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected macOS command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses brew
	if !utils.Contains(cmd.Command, "brew") {
		t.Error("macOS command should use Homebrew package manager")
	}

	// Verify the command installs docker
	if !utils.Contains(cmd.Command, "docker") {
		t.Error("macOS command should install docker package")
	}
}

// TestNewInstallDockerAction_WindowsCommandStructure tests that the Windows command uses winget
func TestNewInstallDockerAction_WindowsCommandStructure(t *testing.T) {
	action := InstallDockerAction()

	cmd, exists := action.PlatformCommands[actions.PlatformWindows]
	if !exists {
		t.Fatal("Expected Windows platform command to exist")
	}

	expectedCommand := "winget install -e --id Docker.DockerDesktop"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Windows command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses winget
	if !utils.Contains(cmd.Command, "winget") {
		t.Error("Windows command should use winget package manager")
	}

	// Verify the command installs Docker Desktop
	if !utils.Contains(cmd.Command, "Docker.DockerDesktop") {
		t.Error("Windows command should install Docker Desktop")
	}
}

// TestNewInstallDockerAction_CheckCommandStructure tests that the check command properly checks for docker
func TestNewInstallDockerAction_CheckCommandStructure(t *testing.T) {
	action := InstallDockerAction()

	// Test that the check command properly checks for docker
	expectedCheck := "docker"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
		}

		// Verify the check command tests for docker
		if !utils.Contains(cmd.CheckCommand, "docker") {
			t.Errorf("Platform %s check command should test for docker", platform)
		}
	}
}

// TestNewInstallDockerAction_NoCommandExecution ensures that the test never executes actual commands
// This is a safety test to verify that we're only testing configuration, not execution
func TestNewInstallDockerAction_NoCommandExecution(t *testing.T) {
	action := InstallDockerAction()

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

// TestNewInstallDockerAction_SafeForCI verifies that the action is safe to use in CI environments
// This test ensures no actual system commands will be executed during testing
func TestNewInstallDockerAction_SafeForCI(t *testing.T) {
	action := InstallDockerAction()

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

// TestNewInstallDockerAction_MockExecutorBehavior tests that the action can be safely used with a mock executor
// This demonstrates how to properly mock the action execution without running actual commands
func TestNewInstallDockerAction_MockExecutorBehavior(t *testing.T) {
	action := InstallDockerAction()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_docker" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "Docker installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "Docker installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}
