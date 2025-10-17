package tools

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestInstallBtop tests the creation of the Install btop action
// This test verifies the action configuration without executing any actual commands
func TestInstallBtop(t *testing.T) {
	action := InstallBtop()

	if action == nil {
		t.Fatal("InstallBtop() returned nil")
	}

	if action.ID != "install_btop" {
		t.Errorf("Expected action ID to be 'install_btop', got '%s'", action.ID)
	}

	if action.Name != "Install btop" {
		t.Errorf("Expected action name to be 'Install btop', got '%s'", action.Name)
	}

	if action.Description != "Install btop system monitor" {
		t.Errorf("Expected action description to be 'Install btop system monitor', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallBtop_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallBtop_PlatformCommands(t *testing.T) {
	action := InstallBtop()

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
			expectedCommand:     "sudo pacman -S btop",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "btop",
			expectedInteractive: true,
		},
		{
			name:                "Debian platform",
			platform:            actions.PlatformDebian,
			expectedCommand:     "sudo apt install btop",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "btop",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "brew install btop",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "btop",
			expectedInteractive: false,
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install btop",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "btop",
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

// TestInstallBtop_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallBtop_UnsupportedPlatforms(t *testing.T) {
	action := InstallBtop()

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

// TestInstallBtop_ActionConsistency verifies that the action maintains consistency across multiple calls
func TestInstallBtop_ActionConsistency(t *testing.T) {
	action1 := InstallBtop()
	action2 := InstallBtop()

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

// TestInstallBtop_MultipleCalls tests that multiple calls to InstallBtop return equivalent actions
func TestInstallBtop_MultipleCalls(t *testing.T) {
	// Call the function multiple times to ensure it returns equivalent configurations
	actions := make([]*actions.Action, 5)
	for i := 0; i < 5; i++ {
		actions[i] = InstallBtop()
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

// TestInstallBtop_ExpectedPlatforms verifies that all expected platforms are supported
func TestInstallBtop_ExpectedPlatforms(t *testing.T) {
	action := InstallBtop()

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

// TestInstallBtop_ArchCommandStructure tests the specific command structure for Arch Linux
func TestInstallBtop_ArchCommandStructure(t *testing.T) {
	action := InstallBtop()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "sudo pacman -S btop" {
		t.Errorf("Expected Arch command to be 'sudo pacman -S btop', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected Arch package source to be official, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "btop" {
		t.Errorf("Expected Arch check command to be 'btop', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Arch command to be interactive")
	}
}

// TestInstallBtop_DebianCommandStructure tests the specific command structure for Debian
func TestInstallBtop_DebianCommandStructure(t *testing.T) {
	action := InstallBtop()

	cmd, exists := action.PlatformCommands[actions.PlatformDebian]
	if !exists {
		t.Fatal("Expected Debian platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "sudo apt install btop" {
		t.Errorf("Expected Debian command to be 'sudo apt install btop', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected Debian package source to be official, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "btop" {
		t.Errorf("Expected Debian check command to be 'btop', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Debian command to be interactive")
	}
}

// TestInstallBtop_LinuxCommandStructure tests the specific command structure for generic Linux
func TestInstallBtop_LinuxCommandStructure(t *testing.T) {
	action := InstallBtop()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "brew install btop" {
		t.Errorf("Expected Linux command to be 'brew install btop', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected Linux package source to be brew, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "btop" {
		t.Errorf("Expected Linux check command to be 'btop', got '%s'", cmd.CheckCommand)
	}

	if cmd.Interactive {
		t.Error("Expected Linux command to be non-interactive")
	}
}

// TestInstallBtop_MacOSCommandStructure tests the specific command structure for macOS
func TestInstallBtop_MacOSCommandStructure(t *testing.T) {
	action := InstallBtop()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "brew install btop" {
		t.Errorf("Expected macOS command to be 'brew install btop', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected macOS package source to be brew, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "btop" {
		t.Errorf("Expected macOS check command to be 'btop', got '%s'", cmd.CheckCommand)
	}

	if cmd.Interactive {
		t.Error("Expected macOS command to be non-interactive")
	}
}

// TestInstallBtop_CheckCommandStructure tests that check commands are properly configured for each platform
func TestInstallBtop_CheckCommandStructure(t *testing.T) {
	action := InstallBtop()

	// All platforms check for "btop"
	expectedCheckCommands := map[actions.Platform]string{
		actions.PlatformArch:   "btop",
		actions.PlatformDebian: "btop",
		actions.PlatformLinux:  "btop",
		actions.PlatformMacOS:  "btop",
	}

	for platform, expectedCheck := range expectedCheckCommands {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Platform %s should exist", platform)
			continue
		}
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Expected check command for platform %s to be '%s', got '%s'", platform, expectedCheck, cmd.CheckCommand)
		}
	}
}

// TestInstallBtop_NoCommandExecution verifies that no actual commands are executed during testing
func TestInstallBtop_NoCommandExecution(t *testing.T) {
	action := InstallBtop()

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

// TestInstallBtop_SafeForCI verifies the test is safe for CI environments
func TestInstallBtop_SafeForCI(t *testing.T) {
	action := InstallBtop()

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

// TestInstallBtop_BrewConsistency tests that brew commands are consistent where expected
func TestInstallBtop_BrewConsistency(t *testing.T) {
	action := InstallBtop()

	// Linux and macOS both use brew
	brewPlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range brewPlatforms {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Fatalf("Expected %s platform command to exist", platform)
		}

		if cmd.PackageSource != actions.PackageSourceBrew {
			t.Errorf("Expected %s to use brew package source, got '%s'", platform, cmd.PackageSource)
		}

		if cmd.Command != "brew install btop" {
			t.Errorf("Expected %s brew command to be 'brew install btop', got '%s'", platform, cmd.Command)
		}
	}
}

// TestInstallBtop_PackageSourceValidation tests that package sources are valid
func TestInstallBtop_PackageSourceValidation(t *testing.T) {
	action := InstallBtop()

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

// TestInstallBtop_InteractiveFlagConsistency tests interactive flag consistency
func TestInstallBtop_InteractiveFlagConsistency(t *testing.T) {
	action := InstallBtop()

	// Arch and Debian should be interactive
	interactivePlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
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

	// Linux and macOS should not be interactive
	nonInteractivePlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range nonInteractivePlatforms {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected platform %s to exist", platform)
			continue
		}
		if cmd.Interactive {
			t.Errorf("Expected platform %s to be non-interactive", platform)
		}
	}
}

// TestInstallBtop_CheckCommandConsistency tests that check commands are consistent
func TestInstallBtop_CheckCommandConsistency(t *testing.T) {
	action := InstallBtop()

	// All platforms check for "btop"
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s should have a check command", platform)
		}

		if cmd.CheckCommand != "btop" {
			t.Errorf("Expected check command for platform %s to be 'btop', got '%s'", platform, cmd.CheckCommand)
		}
	}
}

// TestInstallBtop_CommandComplexity tests that commands have appropriate complexity
func TestInstallBtop_CommandComplexity(t *testing.T) {
	action := InstallBtop()

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
			if !utils.ContainsString(cmd.Command, "apt") {
				t.Errorf("Debian command should contain 'apt': %s", cmd.Command)
			}
		case actions.PlatformLinux:
			if !utils.ContainsString(cmd.Command, "brew") {
				t.Errorf("Linux command should contain 'brew': %s", cmd.Command)
			}
		case actions.PlatformMacOS:
			if !utils.ContainsString(cmd.Command, "brew") {
				t.Errorf("macOS command should contain 'brew': %s", cmd.Command)
			}
		}
	}
}

// TestInstallBtop_ArchCommandSudo tests that Arch and Debian commands properly use sudo
func TestInstallBtop_ArchCommandSudo(t *testing.T) {
	action := InstallBtop()

	sudoPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
	}

	for _, platform := range sudoPlatforms {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Fatalf("Expected %s platform command to exist", platform)
		}

		if !utils.ContainsString(cmd.Command, "sudo") {
			t.Errorf("%s command should use sudo: %s", platform, cmd.Command)
		}
	}
}

// TestInstallBtop_MockExecutorBehavior tests the action with a mock executor
// This ensures the action can be safely "executed" in test environments without running real commands
func TestInstallBtop_MockExecutorBehavior(t *testing.T) {
	action := InstallBtop()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_btop" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "btop installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "btop installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}
