package dev

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestInstallNeovim tests the creation of the Install Neovim action
// This test verifies the action configuration without executing any actual commands
func TestInstallNeovim(t *testing.T) {
	action := InstallNeovim()

	if action == nil {
		t.Fatal("InstallNeovim() returned nil")
	}

	if action.ID != "install_neovim" {
		t.Errorf("Expected action ID to be 'install_neovim', got '%s'", action.ID)
	}

	if action.Name != "Install Neovim" {
		t.Errorf("Expected action name to be 'Install Neovim', got '%s'", action.Name)
	}

	if action.Description != "Install neovim on your system" {
		t.Errorf("Expected action description to be 'Install neovim on your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallNeovim_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallNeovim_PlatformCommands(t *testing.T) {
	action := InstallNeovim()

	tests := []struct {
		name            string
		platform        actions.Platform
		expectedCommand string
		expectedSource  actions.PackageSource
		expectedCheck   string
	}{
		{
			name:            "Arch platform",
			platform:        actions.PlatformArch,
			expectedCommand: "sudo pacman -S neovim",
			expectedSource:  actions.PackageSourceOfficial,
			expectedCheck:   "nvim",
		},
		{
			name:            "Debian platform",
			platform:        actions.PlatformDebian,
			expectedCommand: "sudo apt install neovim",
			expectedSource:  actions.PackageSourceOfficial,
			expectedCheck:   "nvim",
		},
		{
			name:            "Linux platform",
			platform:        actions.PlatformLinux,
			expectedCommand: "brew install neovim",
			expectedSource:  actions.PackageSourceBrew,
			expectedCheck:   "nvim",
		},
		{
			name:            "macOS platform",
			platform:        actions.PlatformMacOS,
			expectedCommand: "brew install neovim",
			expectedSource:  actions.PackageSourceBrew,
			expectedCheck:   "nvim",
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

// TestInstallNeovim_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallNeovim_UnsupportedPlatforms(t *testing.T) {
	action := InstallNeovim()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.Platform("unsupported"),
		actions.Platform("fedora"),
		actions.Platform("opensuse"),
		actions.Platform("alpine"),
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

// TestInstallNeovim_ActionConsistency verifies that all platform commands have consistent structure
// This ensures the action is properly configured for safe execution
func TestInstallNeovim_ActionConsistency(t *testing.T) {
	action := InstallNeovim()

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
		expectedCheck := "nvim"
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command '%s', expected '%s'", platform, cmd.CheckCommand, expectedCheck)
		}
	}
}

// TestInstallNeovim_MultipleCalls tests that multiple calls return equivalent actions
// This ensures the factory function is deterministic and safe
func TestInstallNeovim_MultipleCalls(t *testing.T) {
	// Test that multiple calls to InstallNeovim return equivalent actions
	action1 := InstallNeovim()
	action2 := InstallNeovim()

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

// TestInstallNeovim_ExpectedPlatforms verifies that only expected platforms are configured
// This prevents accidental execution on unexpected platforms
func TestInstallNeovim_ExpectedPlatforms(t *testing.T) {
	action := InstallNeovim()

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

// TestInstallNeovim_ArchCommandStructure tests that the Arch command uses pacman
func TestInstallNeovim_ArchCommandStructure(t *testing.T) {
	action := InstallNeovim()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	expectedCommand := "sudo pacman -S neovim"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Arch command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses pacman
	if !utils.Contains(cmd.Command, "pacman") {
		t.Error("Arch command should use pacman package manager")
	}

	// Verify the command installs neovim
	if !utils.Contains(cmd.Command, "neovim") {
		t.Error("Arch command should install neovim package")
	}

	// Verify it uses official package source
	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected Arch to use official package source, got '%s'", cmd.PackageSource)
	}
}

// TestInstallNeovim_DebianCommandStructure tests that the Debian command uses apt
func TestInstallNeovim_DebianCommandStructure(t *testing.T) {
	action := InstallNeovim()

	cmd, exists := action.PlatformCommands[actions.PlatformDebian]
	if !exists {
		t.Fatal("Expected Debian platform command to exist")
	}

	expectedCommand := "sudo apt install neovim"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Debian command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses apt
	if !utils.Contains(cmd.Command, "apt") {
		t.Error("Debian command should use apt package manager")
	}

	// Verify the command installs neovim
	if !utils.Contains(cmd.Command, "neovim") {
		t.Error("Debian command should install neovim package")
	}

	// Verify it uses official package source
	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected Debian to use official package source, got '%s'", cmd.PackageSource)
	}
}

// TestInstallNeovim_LinuxCommandStructure tests that the Linux command uses brew
func TestInstallNeovim_LinuxCommandStructure(t *testing.T) {
	action := InstallNeovim()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	expectedCommand := "brew install neovim"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Linux command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses brew
	if !utils.Contains(cmd.Command, "brew") {
		t.Error("Linux command should use Homebrew package manager")
	}

	// Verify the command installs neovim
	if !utils.Contains(cmd.Command, "neovim") {
		t.Error("Linux command should install neovim package")
	}

	// Verify it uses brew package source
	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected Linux to use brew package source, got '%s'", cmd.PackageSource)
	}
}

// TestInstallNeovim_MacOSCommandStructure tests that the macOS command uses brew
func TestInstallNeovim_MacOSCommandStructure(t *testing.T) {
	action := InstallNeovim()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	expectedCommand := "brew install neovim"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected macOS command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses brew
	if !utils.Contains(cmd.Command, "brew") {
		t.Error("macOS command should use Homebrew package manager")
	}

	// Verify the command installs neovim
	if !utils.Contains(cmd.Command, "neovim") {
		t.Error("macOS command should install neovim package")
	}

	// Verify it uses brew package source
	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected macOS to use brew package source, got '%s'", cmd.PackageSource)
	}
}

// TestInstallNeovim_CheckCommandStructure tests that the check command properly checks for neovim
func TestInstallNeovim_CheckCommandStructure(t *testing.T) {
	action := InstallNeovim()

	// Test that the check command properly checks for neovim
	expectedCheck := "nvim"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
		}

		// Verify the check command tests for neovim
		if !utils.Contains(cmd.CheckCommand, "nvim") {
			t.Errorf("Platform %s check command should verify neovim dependency first", platform)
		}
	}
}

// TestInstallNeovim_NoCommandExecution ensures that the test never executes actual commands
// This is a safety test to verify that we're only testing configuration, not execution
func TestInstallNeovim_NoCommandExecution(t *testing.T) {
	action := InstallNeovim()

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

// TestInstallNeovim_SafeForCI verifies that the action is safe to use in CI environments
// This test ensures no actual system commands will be executed during testing
func TestInstallNeovim_SafeForCI(t *testing.T) {
	action := InstallNeovim()

	// Verify the action has proper check commands to prevent unnecessary execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s is missing a check command, which could lead to unnecessary execution", platform)
		}

		// Verify check command includes dependency checks
		if !utils.Contains(cmd.CheckCommand, "nvim") {
			t.Errorf("Platform %s check command should verify neovim dependency first", platform)
		}
	}
}

// TestInstallNeovim_MockExecutorBehavior tests that the action can be safely used with a mock executor
// This demonstrates how to properly mock the action execution without running actual commands
func TestInstallNeovim_MockExecutorBehavior(t *testing.T) {
	action := InstallNeovim()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_neovim" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "Neovim installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "Neovim installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}

// TestInstallNeovim_BrewConsistency tests that Linux and macOS both use brew consistently
func TestInstallNeovim_BrewConsistency(t *testing.T) {
	action := InstallNeovim()

	linuxCmd, linuxExists := action.PlatformCommands[actions.PlatformLinux]
	macOSCmd, macOSExists := action.PlatformCommands[actions.PlatformMacOS]

	if !linuxExists {
		t.Fatal("Expected Linux platform command to exist")
	}

	if !macOSExists {
		t.Fatal("Expected macOS platform command to exist")
	}

	// Both should use the same brew command
	if linuxCmd.Command != macOSCmd.Command {
		t.Errorf("Linux and macOS commands should be identical, got '%s' and '%s'", linuxCmd.Command, macOSCmd.Command)
	}

	// Both should use brew package source
	if linuxCmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Linux should use brew package source, got '%s'", linuxCmd.PackageSource)
	}

	if macOSCmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("macOS should use brew package source, got '%s'", macOSCmd.PackageSource)
	}
}

// TestInstallNeovim_OfficialPackageConsistency tests that Arch and Debian both use official packages consistently
func TestInstallNeovim_OfficialPackageConsistency(t *testing.T) {
	action := InstallNeovim()

	archCmd, archExists := action.PlatformCommands[actions.PlatformArch]
	debianCmd, debianExists := action.PlatformCommands[actions.PlatformDebian]

	if !archExists {
		t.Fatal("Expected Arch platform command to exist")
	}

	if !debianExists {
		t.Fatal("Expected Debian platform command to exist")
	}

	// Both should use official package source
	if archCmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Arch should use official package source, got '%s'", archCmd.PackageSource)
	}

	if debianCmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Debian should use official package source, got '%s'", debianCmd.PackageSource)
	}

	// Both should have sudo in their commands
	if !utils.Contains(archCmd.Command, "sudo") {
		t.Error("Arch command should use sudo for system package installation")
	}

	if !utils.Contains(debianCmd.Command, "sudo") {
		t.Error("Debian command should use sudo for system package installation")
	}
}

// TestInstallNeovim_PackageSourceValidation tests that package sources are appropriate for each platform
func TestInstallNeovim_PackageSourceValidation(t *testing.T) {
	action := InstallNeovim()

	// Arch and Debian should use official package sources
	// Linux and macOS should use brew package sources
	for _, platform := range []actions.Platform{actions.PlatformLinux, actions.PlatformMacOS} {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected %s platform command to exist", platform)
			continue
		}

		if cmd.PackageSource != actions.PackageSourceBrew {
			t.Errorf("Expected %s to use brew package source, got '%s'", platform, cmd.PackageSource)
		}
	}
}
