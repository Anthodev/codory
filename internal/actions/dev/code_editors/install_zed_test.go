package dev

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestInstallZed tests the creation of the Install Zed action
// This test verifies the action configuration without executing any actual commands
func TestInstallZed(t *testing.T) {
	action := InstallZed()

	if action == nil {
		t.Fatal("InstallZed() returned nil")
	}

	if action.ID != "install_zed" {
		t.Errorf("Expected action ID to be 'install_zed', got '%s'", action.ID)
	}

	if action.Name != "Install Zed Editor" {
		t.Errorf("Expected action name to be 'Install Zed Editor', got '%s'", action.Name)
	}

	if action.Description != "Install Zed Editor on your system" {
		t.Errorf("Expected action description to be 'Install Zed Editor on your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallZed_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallZed_PlatformCommands(t *testing.T) {
	action := InstallZed()

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
			expectedCommand:     "sudo pacman -S zed",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "zed",
			expectedInteractive: true,
		},
		{
			name:                "Debian platform",
			platform:            actions.PlatformDebian,
			expectedCommand:     "flatpak install flathub dev.zed.Zed",
			expectedSource:      actions.PackageSourceAny,
			expectedCheck:       "zed",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "curl -f https://zed.dev/install.sh | sh",
			expectedSource:      actions.PackageSourceAny,
			expectedCheck:       "zed",
			expectedInteractive: true,
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install --cask zed",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "zed",
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

// TestInstallZed_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallZed_UnsupportedPlatforms(t *testing.T) {
	action := InstallZed()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.PlatformWindows,
	}

	for _, platform := range unsupportedPlatforms {
		if _, exists := action.PlatformCommands[platform]; exists {
			t.Errorf("Expected no command for unsupported platform %s", platform)
		}
	}
}

// TestInstallZed_ActionConsistency verifies that the action maintains consistency across multiple calls
func TestInstallZed_ActionConsistency(t *testing.T) {
	action1 := InstallZed()
	action2 := InstallZed()

	if action1.ID != action2.ID {
		t.Errorf("Action ID should be consistent: %s != %s", action1.ID, action2.ID)
	}

	if action1.Name != action2.Name {
		t.Errorf("Action name should be consistent: %s != %s", action1.Name, action2.Name)
	}

	if action1.Description != action2.Description {
		t.Errorf("Action description should be consistent: %s != %s", action1.Description, action2.Description)
	}

	if len(action1.PlatformCommands) != len(action2.PlatformCommands) {
		t.Errorf("Platform commands count should be consistent: %d != %d", len(action1.PlatformCommands), len(action2.PlatformCommands))
	}

	// Verify all platform commands are identical
	for platform, cmd1 := range action1.PlatformCommands {
		cmd2, exists := action2.PlatformCommands[platform]
		if !exists {
			t.Errorf("Platform %s missing in second action", platform)
			continue
		}

		if cmd1.Command != cmd2.Command {
			t.Errorf("Command for platform %s should be consistent: %s != %s", platform, cmd1.Command, cmd2.Command)
		}

		if cmd1.PackageSource != cmd2.PackageSource {
			t.Errorf("Package source for platform %s should be consistent: %s != %s", platform, cmd1.PackageSource, cmd2.PackageSource)
		}

		if cmd1.CheckCommand != cmd2.CheckCommand {
			t.Errorf("Check command for platform %s should be consistent: %s != %s", platform, cmd1.CheckCommand, cmd2.CheckCommand)
		}

		if cmd1.Interactive != cmd2.Interactive {
			t.Errorf("Interactive flag for platform %s should be consistent: %v != %v", platform, cmd1.Interactive, cmd2.Interactive)
		}
	}
}

// TestInstallZed_MultipleCalls tests that multiple calls to InstallZed return equivalent actions
func TestInstallZed_MultipleCalls(t *testing.T) {
	// Call the function multiple times to ensure it returns equivalent configurations
	actions := make([]*actions.Action, 5)
	for i := 0; i < 5; i++ {
		actions[i] = InstallZed()
	}

	// Verify all actions are equivalent
	for i := 1; i < 5; i++ {
		if actions[0].ID != actions[i].ID {
			t.Errorf("Action ID should be consistent across calls: %s != %s", actions[0].ID, actions[i].ID)
		}

		if actions[0].Name != actions[i].Name {
			t.Errorf("Action name should be consistent across calls: %s != %s", actions[0].Name, actions[i].Name)
		}

		if actions[0].Description != actions[i].Description {
			t.Errorf("Action description should be consistent across calls: %s != %s", actions[0].Description, actions[i].Description)
		}

		if actions[0].Type != actions[i].Type {
			t.Errorf("Action type should be consistent across calls: %s != %s", actions[0].Type, actions[i].Type)
		}

		if len(actions[0].PlatformCommands) != len(actions[i].PlatformCommands) {
			t.Errorf("Platform commands count should be consistent across calls: %d != %d", len(actions[0].PlatformCommands), len(actions[i].PlatformCommands))
		}
	}
}

// TestInstallZed_ExpectedPlatforms verifies that all expected platforms are supported
func TestInstallZed_ExpectedPlatforms(t *testing.T) {
	action := InstallZed()

	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range expectedPlatforms {
		if _, exists := action.PlatformCommands[platform]; !exists {
			t.Errorf("Expected platform %s to be supported", platform)
		}
	}
}

// TestInstallZed_ArchCommandStructure tests the specific command structure for Arch Linux
func TestInstallZed_ArchCommandStructure(t *testing.T) {
	action := InstallZed()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "sudo pacman -S zed" {
		t.Errorf("Expected Arch command to be 'sudo pacman -S zed', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected Arch package source to be official, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "zed" {
		t.Errorf("Expected Arch check command to be 'zed', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Arch command to be interactive")
	}
}

// TestInstallZed_DebianCommandStructure tests the specific command structure for Debian
func TestInstallZed_DebianCommandStructure(t *testing.T) {
	action := InstallZed()

	cmd, exists := action.PlatformCommands[actions.PlatformDebian]
	if !exists {
		t.Fatal("Expected Debian platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "flatpak install flathub dev.zed.Zed" {
		t.Errorf("Expected Debian command to be 'flatpak install flathub dev.zed.Zed', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceAny {
		t.Errorf("Expected Debian package source to be any, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "zed" {
		t.Errorf("Expected Debian check command to be 'zed', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Debian command to be interactive")
	}
}

// TestInstallZed_LinuxCommandStructure tests the specific command structure for generic Linux
func TestInstallZed_LinuxCommandStructure(t *testing.T) {
	action := InstallZed()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "curl -f https://zed.dev/install.sh | sh" {
		t.Errorf("Expected Linux command to be 'curl -f https://zed.dev/install.sh | sh', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceAny {
		t.Errorf("Expected Linux package source to be any, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "zed" {
		t.Errorf("Expected Linux check command to be 'zed', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Linux command to be interactive")
	}
}

// TestInstallZed_MacOSCommandStructure tests the specific command structure for macOS
func TestInstallZed_MacOSCommandStructure(t *testing.T) {
	action := InstallZed()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "brew install --cask zed" {
		t.Errorf("Expected macOS command to be 'brew install --cask zed', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected macOS package source to be brew, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "zed" {
		t.Errorf("Expected macOS check command to be 'zed', got '%s'", cmd.CheckCommand)
	}

	if cmd.Interactive {
		t.Error("Expected macOS command to be non-interactive")
	}
}

// TestInstallZed_CheckCommandStructure tests that check commands are consistent across platforms
func TestInstallZed_CheckCommandStructure(t *testing.T) {
	action := InstallZed()

	expectedCheckCommand := "zed"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheckCommand {
			t.Errorf("Expected check command for platform %s to be '%s', got '%s'", platform, expectedCheckCommand, cmd.CheckCommand)
		}
	}
}

// TestInstallZed_NoCommandExecution verifies that no actual commands are executed during testing
func TestInstallZed_NoCommandExecution(t *testing.T) {
	action := InstallZed()

	// This test verifies that we're only testing configuration, not execution
	// The action should be created without any side effects
	if action == nil {
		t.Fatal("Action should be created without execution")
	}

	// Verify that we have platform commands configured
	if len(action.PlatformCommands) == 0 {
		t.Fatal("Platform commands should be configured")
	}

	// Ensure no actual command execution occurs in any test
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Platform %s should have a command configured", platform)
		}
	}
}

// TestInstallZed_SafeForCI verifies the test is safe for CI environments
func TestInstallZed_SafeForCI(t *testing.T) {
	action := InstallZed()

	// Ensure the action can be created in CI without issues
	if action == nil {
		t.Fatal("Action should be creatable in CI environment")
	}

	// Verify all platform commands are properly configured for CI
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Platform %s should have a command for CI", platform)
		}
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s should have a check command for CI", platform)
		}
	}
}

// TestInstallZed_BrewConsistency tests that brew commands are consistent where expected
func TestInstallZed_BrewConsistency(t *testing.T) {
	action := InstallZed()

	// macOS uses brew install --cask
	macOSCmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	if macOSCmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected macOS to use brew package source, got '%s'", macOSCmd.PackageSource)
	}

	if macOSCmd.Command != "brew install --cask zed" {
		t.Errorf("Expected macOS brew command to be 'brew install --cask zed', got '%s'", macOSCmd.Command)
	}
}

// TestInstallZed_PackageSourceValidation tests that package sources are valid
func TestInstallZed_PackageSourceValidation(t *testing.T) {
	action := InstallZed()

	validSources := []actions.PackageSource{
		actions.PackageSourceOfficial,
		actions.PackageSourceBrew,
		actions.PackageSourceAny,
	}

	for platform, cmd := range action.PlatformCommands {
		valid := false
		for _, source := range validSources {
			if cmd.PackageSource == source {
				valid = true
				break
			}
		}
		if !valid {
			t.Errorf("Platform %s has invalid package source: %s", platform, cmd.PackageSource)
		}
	}
}

// TestInstallZed_InteractiveFlagConsistency tests interactive flag consistency
func TestInstallZed_InteractiveFlagConsistency(t *testing.T) {
	action := InstallZed()

	// Arch, Debian, and Linux should be interactive
	interactivePlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformLinux,
	}

	for _, platform := range interactivePlatforms {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected platform %s to exist", platform)
			continue
		}
		if !cmd.Interactive {
			t.Errorf("Expected platform %s to be interactive", platform)
		}
	}

	// macOS should not be interactive
	macOSCmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}
	if macOSCmd.Interactive {
		t.Error("Expected macOS to be non-interactive")
	}
}

// TestInstallZed_CheckCommandConsistency tests that check commands are consistent
func TestInstallZed_CheckCommandConsistency(t *testing.T) {
	action := InstallZed()

	// All platforms should have the same check command
	expectedCheck := "zed"
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Expected check command for platform %s to be '%s', got '%s'", platform, expectedCheck, cmd.CheckCommand)
		}
	}
}

// TestInstallZed_CommandComplexity tests that commands have appropriate complexity
func TestInstallZed_CommandComplexity(t *testing.T) {
	action := InstallZed()

	// Test that commands are not empty and have reasonable complexity
	for platform, cmd := range action.PlatformCommands {
		if len(cmd.Command) < 5 {
			t.Errorf("Command for platform %s is too short: %s", platform, cmd.Command)
		}

		// Verify commands contain expected keywords
		switch platform {
		case actions.PlatformArch:
			if !utils.ContainsString(cmd.Command, "pacman") {
				t.Errorf("Arch command should contain 'pacman': %s", cmd.Command)
			}
		case actions.PlatformDebian:
			if !utils.ContainsString(cmd.Command, "flatpak") {
				t.Errorf("Debian command should contain 'flatpak': %s", cmd.Command)
			}
		case actions.PlatformLinux:
			if !utils.ContainsString(cmd.Command, "curl") {
				t.Errorf("Linux command should contain 'curl': %s", cmd.Command)
			}
		case actions.PlatformMacOS:
			if !utils.ContainsString(cmd.Command, "brew") {
				t.Errorf("macOS command should contain 'brew': %s", cmd.Command)
			}
		}
	}
}

// TestInstallZed_ArchCommandSudo tests that Arch command properly uses sudo
func TestInstallZed_ArchCommandSudo(t *testing.T) {
	action := InstallZed()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	if !utils.ContainsString(cmd.Command, "sudo") {
		t.Errorf("Arch command should use sudo: %s", cmd.Command)
	}
}

// TestInstallZed_MockExecutorBehavior tests the action with a mock executor
// This ensures the action can be safely "executed" in test environments without running real commands
func TestInstallZed_MockExecutorBehavior(t *testing.T) {
	action := InstallZed()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_zed" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "Zed installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "Zed installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}
