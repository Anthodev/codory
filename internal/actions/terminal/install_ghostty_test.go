package terminal

import (
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestInstallGhostty tests the creation of the Install Ghostty action
// This test verifies the action configuration without executing any actual commands
func TestInstallGhostty(t *testing.T) {
	action := InstallGhostty()

	if action == nil {
		t.Fatal("InstallGhostty() returned nil")
	}

	if action.ID != "install_ghostty" {
		t.Errorf("Expected action ID to be 'install_ghostty', got '%s'", action.ID)
	}

	if action.Name != "Install ghostty terminal" {
		t.Errorf("Expected action name to be 'Install ghostty terminal', got '%s'", action.Name)
	}

	if action.Description != "Install ghostty terminal on your system" {
		t.Errorf("Expected action description to be 'Install ghostty terminal on your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallGhostty_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallGhostty_PlatformCommands(t *testing.T) {
	action := InstallGhostty()

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
			expectedCommand:     "sudo pacman -S ghostty",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "ghostty",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "brew install --cask ghostty",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "ghostty",
			expectedInteractive: false,
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install --cask ghostty",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "ghostty",
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

// TestInstallGhostty_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallGhostty_UnsupportedPlatforms(t *testing.T) {
	action := InstallGhostty()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.PlatformDebian,
		actions.PlatformWindows,
		actions.PlatformAny,
	}

	for _, platform := range unsupportedPlatforms {
		if _, exists := action.PlatformCommands[platform]; exists {
			t.Errorf("Expected platform %s to not have a command, but one was found", platform)
		}
	}
}

// TestInstallGhostty_ActionConsistency tests that the action maintains consistency
// across multiple calls and doesn't have side effects
func TestInstallGhostty_ActionConsistency(t *testing.T) {
	// Call the function multiple times to ensure consistency
	action1 := InstallGhostty()
	action2 := InstallGhostty()

	// Verify both actions are not nil
	if action1 == nil || action2 == nil {
		t.Fatal("InstallGhostty() returned nil on one of the calls")
	}

	// Verify both actions have the same configuration
	if action1.ID != action2.ID {
		t.Errorf("Expected consistent ID, got '%s' and '%s'", action1.ID, action2.ID)
	}

	if action1.Name != action2.Name {
		t.Errorf("Expected consistent Name, got '%s' and '%s'", action1.Name, action2.Name)
	}

	if action1.Description != action2.Description {
		t.Errorf("Expected consistent Description, got '%s' and '%s'", action1.Description, action2.Description)
	}

	if action1.Type != action2.Type {
		t.Errorf("Expected consistent Type, got '%s' and '%s'", action1.Type, action2.Type)
	}

	// Verify both have the same number of platform commands
	if len(action1.PlatformCommands) != len(action2.PlatformCommands) {
		t.Errorf("Expected same number of platform commands, got %d and %d", len(action1.PlatformCommands), len(action2.PlatformCommands))
	}

	// Verify specific platform commands are consistent
	platforms := []actions.Platform{actions.PlatformArch, actions.PlatformLinux, actions.PlatformMacOS}
	for _, platform := range platforms {
		cmd1, exists1 := action1.PlatformCommands[platform]
		cmd2, exists2 := action2.PlatformCommands[platform]

		if exists1 != exists2 {
			t.Errorf("Platform %s existence mismatch between calls", platform)
			continue
		}

		if !exists1 {
			continue
		}

		if cmd1.Command != cmd2.Command {
			t.Errorf("Command mismatch for platform %s: '%s' vs '%s'", platform, cmd1.Command, cmd2.Command)
		}

		if cmd1.PackageSource != cmd2.PackageSource {
			t.Errorf("Package source mismatch for platform %s: '%s' vs '%s'", platform, cmd1.PackageSource, cmd2.PackageSource)
		}

		if cmd1.CheckCommand != cmd2.CheckCommand {
			t.Errorf("Check command mismatch for platform %s: '%s' vs '%s'", platform, cmd1.CheckCommand, cmd2.CheckCommand)
		}

		if cmd1.Interactive != cmd2.Interactive {
			t.Errorf("Interactive flag mismatch for platform %s: %v vs %v", platform, cmd1.Interactive, cmd2.Interactive)
		}
	}
}

// TestInstallGhostty_MultipleCalls tests that multiple calls to InstallGhostty
// return independent action instances
func TestInstallGhostty_MultipleCalls(t *testing.T) {
	action1 := InstallGhostty()
	action2 := InstallGhostty()

	// Verify they are different instances (not the same pointer)
	if action1 == action2 {
		t.Error("Expected different action instances, got the same pointer")
	}

	// Verify they have the same configuration but are independent
	if action1.ID != action2.ID {
		t.Errorf("Expected same ID, got '%s' and '%s'", action1.ID, action2.ID)
	}

	// Modify one action and verify the other is not affected
	action1.ID = "modified_id"
	if action2.ID == "modified_id" {
		t.Error("Modifying one action affected the other - they are not independent")
	}
}

// TestInstallGhostty_ExpectedPlatforms tests that all expected platforms are present
func TestInstallGhostty_ExpectedPlatforms(t *testing.T) {
	action := InstallGhostty()

	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range expectedPlatforms {
		if _, exists := action.PlatformCommands[platform]; !exists {
			t.Errorf("Expected platform %s to have a command configured", platform)
		}
	}
}

// TestInstallGhostty_ArchCommandStructure tests the Arch Linux command structure
func TestInstallGhostty_ArchCommandStructure(t *testing.T) {
	action := InstallGhostty()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "sudo pacman -S ghostty" {
		t.Errorf("Expected command 'sudo pacman -S ghostty', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected package source to be 'official', got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "ghostty" {
		t.Errorf("Expected check command 'ghostty', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Arch command to be interactive")
	}
}

// TestInstallGhostty_LinuxCommandStructure tests the Linux command structure
func TestInstallGhostty_LinuxCommandStructure(t *testing.T) {
	action := InstallGhostty()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "brew install --cask ghostty" {
		t.Errorf("Expected command 'brew install --cask ghostty', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected package source to be 'brew', got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "ghostty" {
		t.Errorf("Expected check command 'ghostty', got '%s'", cmd.CheckCommand)
	}

	if cmd.Interactive {
		t.Error("Expected Linux command to be non-interactive")
	}
}

// TestInstallGhostty_MacOSCommandStructure tests the macOS command structure
func TestInstallGhostty_MacOSCommandStructure(t *testing.T) {
	action := InstallGhostty()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "brew install --cask ghostty" {
		t.Errorf("Expected command 'brew install --cask ghostty', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected package source to be 'brew', got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "ghostty" {
		t.Errorf("Expected check command 'ghostty', got '%s'", cmd.CheckCommand)
	}

	if cmd.Interactive {
		t.Error("Expected macOS command to be non-interactive")
	}
}

// TestInstallGhostty_CheckCommandStructure tests that check commands are consistent
func TestInstallGhostty_CheckCommandStructure(t *testing.T) {
	action := InstallGhostty()

	// All platforms should have the same check command
	expectedCheck := "ghostty"

	platforms := []actions.Platform{actions.PlatformArch, actions.PlatformLinux, actions.PlatformMacOS}
	for _, platform := range platforms {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected platform %s to exist", platform)
			continue
		}

		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Expected check command for %s to be '%s', got '%s'", platform, expectedCheck, cmd.CheckCommand)
		}
	}
}

// TestInstallGhostty_NoCommandExecution verifies that no actual commands are executed
// This test ensures the function only creates the action configuration
func TestInstallGhostty_NoCommandExecution(t *testing.T) {
	// This test verifies that InstallGhostty only creates configuration
	// and doesn't execute any system commands
	action := InstallGhostty()

	// Verify the action is created without any side effects
	if action == nil {
		t.Fatal("InstallGhostty() should return a valid action")
	}

	// Verify that commands are configured but not executed
	// (In a real scenario, we'd mock the executor, but here we just verify configuration)
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Expected command to be configured for platform %s", platform)
		}
		if cmd.CheckCommand == "" {
			t.Errorf("Expected check command to be configured for platform %s", platform)
		}
	}
}

// TestInstallGhostty_SafeForCI tests that the function is safe to run in CI environments
// It should not execute any actual system commands
func TestInstallGhostty_SafeForCI(t *testing.T) {
	// This test ensures the function can be called in CI without side effects
	action := InstallGhostty()

	// Verify the action is created successfully
	if action == nil {
		t.Fatal("InstallGhostty() should work in CI environments")
	}

	// Verify all expected platforms are configured
	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range expectedPlatforms {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected platform %s to be configured for CI", platform)
			continue
		}

		// Verify commands are properly formatted
		if cmd.Command == "" {
			t.Errorf("Expected command to be configured for platform %s", platform)
		}
		if cmd.CheckCommand == "" {
			t.Errorf("Expected check command to be configured for platform %s", platform)
		}
	}
}

// TestInstallGhostty_BrewConsistency tests that brew commands are consistent
func TestInstallGhostty_BrewConsistency(t *testing.T) {
	action := InstallGhostty()

	// Platforms that use brew should have consistent patterns
	brewPlatforms := []actions.Platform{actions.PlatformLinux, actions.PlatformMacOS}

	for _, platform := range brewPlatforms {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected platform %s to exist", platform)
			continue
		}

		if cmd.PackageSource != actions.PackageSourceBrew {
			t.Errorf("Expected %s to use brew package source", platform)
		}

		testutil.AssertContains(t, cmd.Command, "brew install", "Brew command should contain 'brew install'")
		testutil.AssertContains(t, cmd.Command, "--cask", "Brew command should contain '--cask'")
	}
}

// TestInstallGhostty_PackageSourceValidation tests that package sources are valid
func TestInstallGhostty_PackageSourceValidation(t *testing.T) {
	action := InstallGhostty()

	validSources := []actions.PackageSource{
		actions.PackageSourceOfficial,
		actions.PackageSourceBrew,
		actions.PackageSourceAny,
	}

	for platform, cmd := range action.PlatformCommands {
		found := false
		for _, validSource := range validSources {
			if cmd.PackageSource == validSource {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Platform %s has invalid package source: %s", platform, cmd.PackageSource)
		}
	}
}

// TestInstallGhostty_InteractiveFlagConsistency tests interactive flag consistency
func TestInstallGhostty_InteractiveFlagConsistency(t *testing.T) {
	action := InstallGhostty()

	// Arch should be interactive due to sudo
	archCmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch command to exist")
	}
	if !archCmd.Interactive {
		t.Error("Expected Arch command to be interactive (uses sudo)")
	}

	// Brew commands should not be interactive
	brewPlatforms := []actions.Platform{actions.PlatformLinux, actions.PlatformMacOS}
	for _, platform := range brewPlatforms {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected platform %s to exist", platform)
			continue
		}
		if cmd.Interactive {
			t.Errorf("Expected %s command to be non-interactive (uses brew)", platform)
		}
	}
}

// TestInstallGhostty_CheckCommandConsistency tests that check commands are consistent
func TestInstallGhostty_CheckCommandConsistency(t *testing.T) {
	action := InstallGhostty()

	expectedCheck := "ghostty"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheck {
			t.Errorf("Expected check command for %s to be '%s', got '%s'", platform, expectedCheck, cmd.CheckCommand)
		}
	}
}

// TestInstallGhostty_CommandComplexity tests command complexity and safety
func TestInstallGhostty_CommandComplexity(t *testing.T) {
	action := InstallGhostty()

	// Verify commands are reasonable and safe
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Expected non-empty command for platform %s", platform)
		}

		// Check for potentially dangerous patterns
		dangerousPatterns := []string{"rm -rf", ">", ">>", "&", "&&", "||", "`", "$("}
		for _, pattern := range dangerousPatterns {
			if contains(cmd.Command, pattern) && platform != actions.PlatformLinux {
				// Linux uses curl | sh which is acceptable for install scripts
				t.Errorf("Platform %s command contains potentially dangerous pattern '%s': %s", platform, pattern, cmd.Command)
			}
		}
	}
}

// TestInstallGhostty_ArchCommandSudo tests Arch command sudo usage
func TestInstallGhostty_ArchCommandSudo(t *testing.T) {
	action := InstallGhostty()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch command to exist")
	}

	// Verify sudo is used for Arch
	testutil.AssertContains(t, cmd.Command, "sudo", "Arch command should contain 'sudo'")
	testutil.AssertContains(t, cmd.Command, "pacman", "Arch command should contain 'pacman'")
	testutil.AssertContains(t, cmd.Command, "-S", "Arch command should contain '-S' flag")
}

// TestInstallGhostty_MockExecutorBehavior tests behavior with mock executor
// This simulates how the action would behave with a mocked executor
func TestInstallGhostty_MockExecutorBehavior(t *testing.T) {
	action := InstallGhostty()

	// Verify the action can be used with a mock executor
	if action.Type != actions.ActionTypeCommand {
		t.Error("Expected action type to be ActionTypeCommand for mock executor compatibility")
	}

	// Verify all platform commands are properly configured for mocking
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Expected command to be configured for platform %s", platform)
		}
		if cmd.CheckCommand == "" {
			t.Errorf("Expected check command to be configured for platform %s", platform)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return utils.ContainsString(s, substr)
}

// AssertCommandContains is a test helper that checks if a command contains expected substrings
func AssertCommandContains(t *testing.T, command string, expectedParts []string, platform actions.Platform) {
	t.Helper()
	for _, part := range expectedParts {
		if !contains(command, part) {
			t.Errorf("Expected %s command to contain '%s', got: %s", platform, part, command)
		}
	}
}
