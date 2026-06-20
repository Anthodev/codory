package tools

import (
	"fmt"
	"strings"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

func assertCommandContains(t *testing.T, command string, parts ...string) {
	t.Helper()
	for _, part := range parts {
		if !strings.Contains(command, part) {
			t.Fatalf("command %q missing %q", command, part)
		}
	}
}

// TestInstallBat tests the creation of the Install bat action
// This test verifies the action configuration without executing any actual commands
func TestInstallBat(t *testing.T) {
	action := InstallBat()

	if action == nil {
		t.Fatal("InstallBat() returned nil")
	}

	if action.ID != "install_bat" {
		t.Errorf("Expected action ID to be 'install_bat', got '%s'", action.ID)
	}

	if action.Name != "Install bat" {
		t.Errorf("Expected action name to be 'Install bat', got '%s'", action.Name)
	}

	if action.Description != "Install the bat content viewer tool" {
		t.Errorf("Expected action description to be 'Install the bat content viewer tool', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallBat_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallBat_PlatformCommands(t *testing.T) {
	action := InstallBat()

	tests := []struct {
		name                 string
		platform             actions.Platform
		expectedCommandParts []string
		expectedSource       actions.PackageSource
		expectedCheck        string
		expectedInteractive  bool
	}{
		{
			name:                 "Arch platform",
			platform:             actions.PlatformArch,
			expectedCommandParts: []string{"pacman", "bat"},
			expectedSource:       actions.PackageSourceOfficial,
			expectedCheck:        "bat",
			expectedInteractive:  true,
		},
		{
			name:                 "Debian platform",
			platform:             actions.PlatformDebian,
			expectedCommandParts: []string{"apt", "bat"},
			expectedSource:       actions.PackageSourceOfficial,
			expectedCheck:        "bat",
			expectedInteractive:  true,
		},
		{
			name:                 "Linux platform",
			platform:             actions.PlatformLinux,
			expectedCommandParts: []string{"brew", "bat"},
			expectedSource:       actions.PackageSourceBrew,
			expectedCheck:        "bat",
			expectedInteractive:  false,
		},
		{
			name:                 "macOS platform",
			platform:             actions.PlatformMacOS,
			expectedCommandParts: []string{"brew", "bat"},
			expectedSource:       actions.PackageSourceBrew,
			expectedCheck:        "bat",
			expectedInteractive:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, exists := action.PlatformCommands[tt.platform]
			if !exists {
				t.Fatalf("Expected platform command for %s to exist", tt.platform)
			}

			assertCommandContains(t, cmd.Command, tt.expectedCommandParts...)

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

// TestInstallBat_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallBat_UnsupportedPlatforms(t *testing.T) {
	action := InstallBat()

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

// TestInstallBat_ActionConsistency verifies that the action maintains consistency across multiple calls
func TestInstallBat_ActionConsistency(t *testing.T) {
	action1 := InstallBat()
	action2 := InstallBat()

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

// TestInstallBat_MultipleCalls tests that multiple calls to InstallBat return equivalent actions
func TestInstallBat_MultipleCalls(t *testing.T) {
	// Call the function multiple times to ensure it returns equivalent configurations
	actions := make([]*actions.Action, 5)
	for i := 0; i < 5; i++ {
		actions[i] = InstallBat()
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

// TestInstallBat_ExpectedPlatforms verifies that all expected platforms are supported
func TestInstallBat_ExpectedPlatforms(t *testing.T) {
	action := InstallBat()

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

// TestInstallBat_ArchCommandStructure tests the specific command structure for Arch Linux
func TestInstallBat_ArchCommandStructure(t *testing.T) {
	action := InstallBat()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	// Verify command structure
	assertCommandContains(t, cmd.Command, "pacman", "bat")

	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected Arch package source to be official, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "bat" {
		t.Errorf("Expected Arch check command to be 'bat', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Arch command to be interactive")
	}
}

// TestInstallBat_DebianCommandStructure tests the specific command structure for Debian
func TestInstallBat_DebianCommandStructure(t *testing.T) {
	action := InstallBat()

	cmd, exists := action.PlatformCommands[actions.PlatformDebian]
	if !exists {
		t.Fatal("Expected Debian platform command to exist")
	}

	// Verify command structure
	assertCommandContains(t, cmd.Command, "apt", "bat")

	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected Debian package source to be official, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "bat" {
		t.Errorf("Expected Debian check command to be 'bat', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Debian command to be interactive")
	}
}

// TestInstallBat_LinuxCommandStructure tests the specific command structure for generic Linux
func TestInstallBat_LinuxCommandStructure(t *testing.T) {
	action := InstallBat()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	// Verify command structure
	assertCommandContains(t, cmd.Command, "brew", "bat")

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected Linux package source to be brew, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "bat" {
		t.Errorf("Expected Linux check command to be 'bat', got '%s'", cmd.CheckCommand)
	}

	if cmd.Interactive {
		t.Error("Expected Linux command to be non-interactive")
	}
}

// TestInstallBat_MacOSCommandStructure tests the specific command structure for macOS
func TestInstallBat_MacOSCommandStructure(t *testing.T) {
	action := InstallBat()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	// Verify command structure
	assertCommandContains(t, cmd.Command, "brew", "bat")

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected macOS package source to be brew, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "bat" {
		t.Errorf("Expected macOS check command to be 'bat', got '%s'", cmd.CheckCommand)
	}

	if cmd.Interactive {
		t.Error("Expected macOS command to be non-interactive")
	}
}

// TestInstallBat_CheckCommandStructure tests that check commands are properly configured for each platform
func TestInstallBat_CheckCommandStructure(t *testing.T) {
	action := InstallBat()

	// Define expected check commands per platform
	expectedCheckCommands := map[actions.Platform]string{
		actions.PlatformArch:   "bat",
		actions.PlatformDebian: "bat",
		actions.PlatformLinux:  "bat",
		actions.PlatformMacOS:  "bat",
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

// TestInstallBat_NoCommandExecution verifies that no actual commands are executed during testing
func TestInstallBat_NoCommandExecution(t *testing.T) {
	action := InstallBat()

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

// TestInstallBat_SafeForCI verifies the test is safe for CI environments
func TestInstallBat_SafeForCI(t *testing.T) {
	action := InstallBat()

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

// TestInstallBat_BrewConsistency tests that brew commands are consistent where expected
func TestInstallBat_BrewConsistency(t *testing.T) {
	action := InstallBat()

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

		assertCommandContains(t, cmd.Command, "brew", "bat")
	}
}

// TestInstallBat_PackageSourceValidation tests that package sources are valid
func TestInstallBat_PackageSourceValidation(t *testing.T) {
	action := InstallBat()

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

// TestInstallBat_InteractiveFlagConsistency tests interactive flag consistency
func TestInstallBat_InteractiveFlagConsistency(t *testing.T) {
	action := InstallBat()

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

// TestInstallBat_CheckCommandConsistency tests that check commands are consistent
func TestInstallBat_CheckCommandConsistency(t *testing.T) {
	action := InstallBat()

	// All platforms check for "bat"
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s should have a check command", platform)
		}

		// Verify check command is not arbitrary
		if cmd.CheckCommand != "bat" {
			t.Errorf("Expected %s check command to be 'bat', got '%s'", platform, cmd.CheckCommand)
		}
	}
}

// TestInstallBat_CommandComplexity tests that commands have appropriate complexity
func TestInstallBat_CommandComplexity(t *testing.T) {
	action := InstallBat()

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

// TestInstallBat_ArchCommandSudo tests that Arch and Debian commands properly use sudo
func TestInstallBat_ArchCommandSudo(t *testing.T) {
	action := InstallBat()

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

// TestInstallBat_MockExecutorBehavior tests the action with a mock executor
// This ensures the action can be safely "executed" in test environments without running real commands
func TestInstallBat_MockExecutorBehavior(t *testing.T) {
	action := InstallBat()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_bat" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "bat installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "bat installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}
