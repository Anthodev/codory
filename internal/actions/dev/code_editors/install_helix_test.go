package dev

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestInstallHelix tests the creation of the Install Helix action
// This test verifies the action configuration without executing any actual commands
func TestInstallHelix(t *testing.T) {
	action := InstallHelix()

	if action == nil {
		t.Fatal("InstallHelix() returned nil")
	}

	if action.ID != "install_helix" {
		t.Errorf("Expected action ID to be 'install_helix', got '%s'", action.ID)
	}

	if action.Name != "Install Helix Editor" {
		t.Errorf("Expected action name to be 'Install Helix Editor', got '%s'", action.Name)
	}

	if action.Description != "Install Helix Editor on your system" {
		t.Errorf("Expected action description to be 'Install Helix Editor on your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallHelix_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallHelix_PlatformCommands(t *testing.T) {
	action := InstallHelix()

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
			expectedCommand:     "sudo pacman -S helix",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "helix",
			expectedInteractive: true,
		},
		{
			name:                "Debian platform",
			platform:            actions.PlatformDebian,
			expectedCommand:     "sudo add-apt-repository ppa:maveonair/helix-editor && sudo apt-get update && sudo apt install helix",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "helix",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "brew install helix",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "helix",
			expectedInteractive: false,
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install helix",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "helix",
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
				t.Errorf("Expected interactive for %s to be %v, got %v", tt.platform, tt.expectedInteractive, cmd.Interactive)
			}
		})
	}
}

// TestInstallHelix_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallHelix_UnsupportedPlatforms(t *testing.T) {
	action := InstallHelix()

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

// TestInstallHelix_ActionConsistency verifies that all platform commands have consistent structure
// This ensures the action is properly configured for safe execution
func TestInstallHelix_ActionConsistency(t *testing.T) {
	action := InstallHelix()

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
		expectedCheck := "helix"
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command '%s', expected '%s'", platform, cmd.CheckCommand, expectedCheck)
		}
	}
}

// TestInstallHelix_MultipleCalls tests that multiple calls return equivalent actions
// This ensures the factory function is deterministic and safe
func TestInstallHelix_MultipleCalls(t *testing.T) {
	// Test that multiple calls to InstallHelix return equivalent actions
	action1 := InstallHelix()
	action2 := InstallHelix()

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
			t.Errorf("Platform %s interactive flags differ: %v vs %v", platform, cmd1.Interactive, cmd2.Interactive)
		}
	}
}

// TestInstallHelix_ExpectedPlatforms verifies that only expected platforms are configured
// This prevents accidental execution on unexpected platforms
func TestInstallHelix_ExpectedPlatforms(t *testing.T) {
	action := InstallHelix()

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

// TestInstallHelix_ArchCommandStructure tests that the Arch command uses pacman
func TestInstallHelix_ArchCommandStructure(t *testing.T) {
	action := InstallHelix()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	expectedCommand := "sudo pacman -S helix"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Arch command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses pacman
	if !utils.Contains(cmd.Command, "pacman") {
		t.Error("Arch command should use pacman package manager")
	}

	// Verify the command installs helix
	if !utils.Contains(cmd.Command, "helix") {
		t.Error("Arch command should install helix package")
	}

	// Verify it's interactive
	if !cmd.Interactive {
		t.Error("Arch command should be interactive")
	}
}

// TestInstallHelix_DebianCommandStructure tests that the Debian command uses apt
func TestInstallHelix_DebianCommandStructure(t *testing.T) {
	action := InstallHelix()

	cmd, exists := action.PlatformCommands[actions.PlatformDebian]
	if !exists {
		t.Fatal("Expected Debian platform command to exist")
	}

	expectedCommand := "sudo add-apt-repository ppa:maveonair/helix-editor && sudo apt-get update && sudo apt install helix"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Debian command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses apt
	if !utils.Contains(cmd.Command, "apt") {
		t.Error("Debian command should use apt package manager")
	}

	// Verify the command adds the PPA
	if !utils.Contains(cmd.Command, "add-apt-repository") {
		t.Error("Debian command should add PPA repository")
	}

	// Verify the command installs helix
	if !utils.Contains(cmd.Command, "helix") {
		t.Error("Debian command should install helix package")
	}

	// Verify it's interactive
	if !cmd.Interactive {
		t.Error("Debian command should be interactive")
	}
}

// TestInstallHelix_LinuxCommandStructure tests that the Linux command uses brew
func TestInstallHelix_LinuxCommandStructure(t *testing.T) {
	action := InstallHelix()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	expectedCommand := "brew install helix"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Linux command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses brew
	if !utils.Contains(cmd.Command, "brew") {
		t.Error("Linux command should use Homebrew package manager")
	}

	// Verify the command installs helix
	if !utils.Contains(cmd.Command, "helix") {
		t.Error("Linux command should install helix package")
	}

	// Verify it's not interactive
	if cmd.Interactive {
		t.Error("Linux command should not be interactive")
	}
}

// TestInstallHelix_MacOSCommandStructure tests that the macOS command uses brew
func TestInstallHelix_MacOSCommandStructure(t *testing.T) {
	action := InstallHelix()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	expectedCommand := "brew install helix"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected macOS command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	// Verify the command uses brew
	if !utils.Contains(cmd.Command, "brew") {
		t.Error("macOS command should use Homebrew package manager")
	}

	// Verify the command installs helix
	if !utils.Contains(cmd.Command, "helix") {
		t.Error("macOS command should install helix package")
	}

	// Verify it's not interactive
	if cmd.Interactive {
		t.Error("macOS command should not be interactive")
	}
}

// TestInstallHelix_CheckCommandStructure tests that the check command properly checks for helix
func TestInstallHelix_CheckCommandStructure(t *testing.T) {
	action := InstallHelix()

	// Test that the check command properly checks for helix
	expectedCheck := "helix"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command: '%s'", platform, cmd.CheckCommand)
		}

		// Verify the check command tests for helix (helix)
		if !utils.Contains(cmd.CheckCommand, "helix") {
			t.Errorf("Platform %s check command should verify helix dependency first", platform)
		}
	}
}

// TestInstallHelix_NoCommandExecution ensures that the test never executes actual commands
// This is a safety test to verify that we're only testing configuration, not execution
func TestInstallHelix_NoCommandExecution(t *testing.T) {
	action := InstallHelix()

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

// TestInstallHelix_SafeForCI verifies that the action is safe to use in CI environments
// This test ensures no actual system commands will be executed during testing
func TestInstallHelix_SafeForCI(t *testing.T) {
	action := InstallHelix()

	// Verify the action has proper check commands to prevent unnecessary execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s is missing a check command, which could lead to unnecessary execution", platform)
		}

		// Verify check command includes dependency checks
		if !utils.Contains(cmd.CheckCommand, "helix") {
			t.Errorf("Platform %s check command should verify helix dependency first", platform)
		}
	}
}

// TestInstallHelix_MockExecutorBehavior tests that the action can be safely used with a mock executor
// This demonstrates how to properly mock the action execution without running actual commands
func TestInstallHelix_MockExecutorBehavior(t *testing.T) {
	action := InstallHelix()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_helix" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "Helix installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "Helix installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}

// TestInstallHelix_BrewConsistency tests that Linux and macOS both use brew consistently
func TestInstallHelix_BrewConsistency(t *testing.T) {
	action := InstallHelix()

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

	// Both should not be interactive
	if linuxCmd.Interactive {
		t.Error("Linux command should not be interactive")
	}

	if macOSCmd.Interactive {
		t.Error("macOS command should not be interactive")
	}
}

// TestInstallHelix_PackageSourceValidation tests that package sources are appropriate for each platform
func TestInstallHelix_PackageSourceValidation(t *testing.T) {
	action := InstallHelix()

	// Arch and Debian should use official package sources
	for _, platform := range []actions.Platform{actions.PlatformArch, actions.PlatformDebian} {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected %s platform command to exist", platform)
			continue
		}

		if cmd.PackageSource != actions.PackageSourceOfficial {
			t.Errorf("Expected %s to use official package source, got '%s'", platform, cmd.PackageSource)
		}
	}

	// Linux and macOS should use brew
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

// TestInstallHelix_InteractiveFlagConsistency tests that interactive flags are set appropriately
func TestInstallHelix_InteractiveFlagConsistency(t *testing.T) {
	action := InstallHelix()

	// Arch and Debian should be interactive (require sudo)
	interactivePlatforms := []actions.Platform{actions.PlatformArch, actions.PlatformDebian}
	for _, platform := range interactivePlatforms {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected %s platform command to exist", platform)
			continue
		}

		if !cmd.Interactive {
			t.Errorf("Expected %s command to be interactive (requires sudo)", platform)
		}
	}

	// Linux and macOS should not be interactive (brew doesn't require sudo)
	nonInteractivePlatforms := []actions.Platform{actions.PlatformLinux, actions.PlatformMacOS}
	for _, platform := range nonInteractivePlatforms {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected %s platform command to exist", platform)
			continue
		}

		if cmd.Interactive {
			t.Errorf("Expected %s command to not be interactive (brew doesn't require sudo)", platform)
		}
	}
}

// TestInstallHelix_CheckCommandConsistency tests that all platforms use the same check command
func TestInstallHelix_CheckCommandConsistency(t *testing.T) {
	action := InstallHelix()

	expectedCheck := "helix"
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Platform %s has unexpected check command '%s', expected '%s'", platform, cmd.CheckCommand, expectedCheck)
		}
	}
}

// TestInstallHelix_CommandComplexity tests that complex commands are properly structured
func TestInstallHelix_CommandComplexity(t *testing.T) {
	action := InstallHelix()

	// Test Debian command complexity (multiple commands chained)
	debianCmd, exists := action.PlatformCommands[actions.PlatformDebian]
	if !exists {
		t.Fatal("Expected Debian platform command to exist")
	}

	// Should contain multiple commands separated by &&
	if !utils.Contains(debianCmd.Command, "&&") {
		t.Error("Debian command should contain multiple commands chained with &&")
	}

	// Should add PPA repository
	if !utils.Contains(debianCmd.Command, "add-apt-repository") {
		t.Error("Debian command should add PPA repository")
	}

	// Should update package list
	if !utils.Contains(debianCmd.Command, "apt-get update") {
		t.Error("Debian command should update package list")
	}

	// Should install the package
	if !utils.Contains(debianCmd.Command, "apt install helix") {
		t.Error("Debian command should install helix package")
	}
}

// TestInstallHelix_ArchCommandSudo tests that Arch command properly uses sudo
func TestInstallHelix_ArchCommandSudo(t *testing.T) {
	action := InstallHelix()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	if !utils.Contains(cmd.Command, "sudo") {
		t.Error("Arch command should use sudo for pacman")
	}
}

// TestInstallHelix_DebianCommandSudo tests that Debian command properly uses sudo
func TestInstallHelix_DebianCommandSudo(t *testing.T) {
	action := InstallHelix()

	cmd, exists := action.PlatformCommands[actions.PlatformDebian]
	if !exists {
		t.Fatal("Expected Debian platform command to exist")
	}

	// Should use sudo for add-apt-repository
	if !utils.Contains(cmd.Command, "sudo add-apt-repository") {
		t.Error("Debian command should use sudo for add-apt-repository")
	}

	// Should use sudo for apt install
	if !utils.Contains(cmd.Command, "sudo apt install") {
		t.Error("Debian command should use sudo for apt install")
	}
}
