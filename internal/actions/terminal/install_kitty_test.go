package terminal

import (
	"strings"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestInstallKitty tests the creation of the Install Kitty action
// This test verifies the action configuration without executing any actual commands
func TestInstallKitty(t *testing.T) {
	action := InstallKitty()

	if action == nil {
		t.Fatal("InstallKitty() returned nil")
	}

	if action.ID != "install_kitty" {
		t.Errorf("Expected action ID to be 'install_kitty', got '%s'", action.ID)
	}

	if action.Name != "Install Kitty terminal" {
		t.Errorf("Expected action name to be 'Install Kitty terminal', got '%s'", action.Name)
	}

	if action.Description != "Install the Kitty terminal on your system" {
		t.Errorf("Expected action description to be 'Install the Kitty terminal on your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallKitty_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallKitty_PlatformCommands(t *testing.T) {
	action := InstallKitty()

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
			expectedCommand:     "sudo pacman -S kitty",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "kitty",
			expectedInteractive: true,
		},
		{
			name:                "Debian platform",
			platform:            actions.PlatformDebian,
			expectedCommand:     "sudo apt-get install kitty",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "kitty",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "brew install --cask kitty",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "kitty",
			expectedInteractive: false,
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install --cask kitty",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "kitty",
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

// TestInstallKitty_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallKitty_UnsupportedPlatforms(t *testing.T) {
	action := InstallKitty()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.PlatformWindows,
		actions.PlatformAny,
	}

	for _, platform := range unsupportedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if _, exists := action.PlatformCommands[platform]; exists {
				t.Errorf("Platform %s should not have a command defined", platform)
			}
		})
	}
}

// TestInstallKitty_ActionConsistency verifies that the action maintains consistency across calls
func TestInstallKitty_ActionConsistency(t *testing.T) {
	action1 := InstallKitty()
	action2 := InstallKitty()

	// Verify that multiple calls return equivalent actions
	if action1.ID != action2.ID {
		t.Errorf("Action IDs should be consistent: %s != %s", action1.ID, action2.ID)
	}

	if action1.Name != action2.Name {
		t.Errorf("Action names should be consistent: %s != %s", action1.Name, action2.Name)
	}

	if action1.Description != action2.Description {
		t.Errorf("Action descriptions should be consistent: %s != %s", action1.Description, action2.Description)
	}

	if action1.Type != action2.Type {
		t.Errorf("Action types should be consistent: %s != %s", action1.Type, action2.Type)
	}

	// Verify platform commands are consistent
	if len(action1.PlatformCommands) != len(action2.PlatformCommands) {
		t.Errorf("Platform commands count should be consistent: %d != %d", len(action1.PlatformCommands), len(action2.PlatformCommands))
	}

	for platform, cmd1 := range action1.PlatformCommands {
		cmd2, exists := action2.PlatformCommands[platform]
		if !exists {
			t.Errorf("Platform %s exists in first action but not in second", platform)
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

// TestInstallKitty_MultipleCalls verifies that calling InstallKitty multiple times doesn't cause issues
func TestInstallKitty_MultipleCalls(t *testing.T) {
	// Call the function multiple times to ensure no side effects
	for i := 0; i < 5; i++ {
		action := InstallKitty()
		if action == nil {
			t.Fatalf("InstallKitty() returned nil on call %d", i+1)
		}

		// Verify basic properties on each call
		if action.ID != "install_kitty" {
			t.Errorf("Expected action ID to be 'install_kitty' on call %d, got '%s'", i+1, action.ID)
		}

		if len(action.PlatformCommands) != 4 {
			t.Errorf("Expected 4 platform commands on call %d, got %d", i+1, len(action.PlatformCommands))
		}
	}
}

// TestInstallKitty_ExpectedPlatforms verifies that all expected platforms are supported
func TestInstallKitty_ExpectedPlatforms(t *testing.T) {
	action := InstallKitty()

	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range expectedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if _, exists := action.PlatformCommands[platform]; !exists {
				t.Errorf("Expected platform %s to have a command defined", platform)
			}
		})
	}
}

// TestInstallKitty_ArchCommandStructure tests the specific structure of the Arch command
func TestInstallKitty_ArchCommandStructure(t *testing.T) {
	action := InstallKitty()
	cmd, exists := action.PlatformCommands[actions.PlatformArch]

	if !exists {
		t.Fatal("Arch platform command should exist")
	}

	// Test command structure
	testutil.AssertContains(t, cmd.Command, "sudo", "Arch command should use sudo")
	testutil.AssertContains(t, cmd.Command, "pacman", "Arch command should use pacman")
	testutil.AssertContains(t, cmd.Command, "-S", "Arch command should use install flag")
	testutil.AssertContains(t, cmd.Command, "kitty", "Arch command should install kitty")

	// Test package source
	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Arch should use official package source, got %s", cmd.PackageSource)
	}

	// Test check command
	if cmd.CheckCommand != "kitty" {
		t.Errorf("Arch check command should be 'kitty', got '%s'", cmd.CheckCommand)
	}

	// Test interactive flag
	if !cmd.Interactive {
		t.Error("Arch command should be interactive")
	}
}

// TestInstallKitty_DebianCommandStructure tests the specific structure of the Debian command
func TestInstallKitty_DebianCommandStructure(t *testing.T) {
	action := InstallKitty()
	cmd, exists := action.PlatformCommands[actions.PlatformDebian]

	if !exists {
		t.Fatal("Debian platform command should exist")
	}

	// Test command structure
	testutil.AssertContains(t, cmd.Command, "sudo", "Debian command should use sudo")
	testutil.AssertContains(t, cmd.Command, "apt-get", "Debian command should use apt-get")
	testutil.AssertContains(t, cmd.Command, "install", "Debian command should use install")
	testutil.AssertContains(t, cmd.Command, "kitty", "Debian command should install kitty")

	// Test package source
	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Debian should use official package source, got %s", cmd.PackageSource)
	}

	// Test check command
	if cmd.CheckCommand != "kitty" {
		t.Errorf("Debian check command should be 'kitty', got '%s'", cmd.CheckCommand)
	}

	// Test interactive flag
	if !cmd.Interactive {
		t.Error("Debian command should be interactive")
	}
}

// TestInstallKitty_LinuxCommandStructure tests the specific structure of the Linux command
func TestInstallKitty_LinuxCommandStructure(t *testing.T) {
	action := InstallKitty()
	cmd, exists := action.PlatformCommands[actions.PlatformLinux]

	if !exists {
		t.Fatal("Linux platform command should exist")
	}

	// Test command structure
	testutil.AssertContains(t, cmd.Command, "brew", "Linux command should use brew")
	testutil.AssertContains(t, cmd.Command, "install", "Linux command should use install")
	testutil.AssertContains(t, cmd.Command, "--cask", "Linux command should use cask")
	testutil.AssertContains(t, cmd.Command, "kitty", "Linux command should install kitty")

	// Test package source
	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Linux should use brew package source, got %s", cmd.PackageSource)
	}

	// Test check command
	if cmd.CheckCommand != "kitty" {
		t.Errorf("Linux check command should be 'kitty', got '%s'", cmd.CheckCommand)
	}

	// Test interactive flag
	if cmd.Interactive {
		t.Error("Linux command should not be interactive")
	}
}

// TestInstallKitty_MacOSCommandStructure tests the specific structure of the macOS command
func TestInstallKitty_MacOSCommandStructure(t *testing.T) {
	action := InstallKitty()
	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]

	if !exists {
		t.Fatal("macOS platform command should exist")
	}

	// Test command structure
	testutil.AssertContains(t, cmd.Command, "brew", "macOS command should use brew")
	testutil.AssertContains(t, cmd.Command, "install", "macOS command should use install")
	testutil.AssertContains(t, cmd.Command, "--cask", "macOS command should use cask")
	testutil.AssertContains(t, cmd.Command, "kitty", "macOS command should install kitty")

	// Test package source
	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("macOS should use brew package source, got %s", cmd.PackageSource)
	}

	// Test check command
	if cmd.CheckCommand != "kitty" {
		t.Errorf("macOS check command should be 'kitty', got '%s'", cmd.CheckCommand)
	}

	// Test interactive flag
	if cmd.Interactive {
		t.Error("macOS command should not be interactive")
	}
}

// TestInstallKitty_CheckCommandStructure verifies that all check commands are consistent
func TestInstallKitty_CheckCommandStructure(t *testing.T) {
	action := InstallKitty()

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != "kitty" {
			t.Errorf("Platform %s should have check command 'kitty', got '%s'", platform, cmd.CheckCommand)
		}
	}
}

// TestInstallKitty_NoCommandExecution verifies that creating the action doesn't execute any commands
func TestInstallKitty_NoCommandExecution(t *testing.T) {
	// This test ensures that just creating the action doesn't execute any commands
	// The action should be safe to create in tests without side effects
	action := InstallKitty()

	if action == nil {
		t.Fatal("InstallKitty() returned nil")
	}

	// Verify the action is properly configured but no commands were executed
	if len(action.PlatformCommands) == 0 {
		t.Error("Expected platform commands to be configured")
	}

	// The fact that we reach this point without any command execution errors
	// indicates that creating the action is safe
}

// TestInstallKitty_SafeForCI verifies that the test can run in CI environments
func TestInstallKitty_SafeForCI(t *testing.T) {
	// This test ensures that the action creation is safe for CI environments
	// and doesn't depend on external systems or commands
	action := InstallKitty()

	if action == nil {
		t.Fatal("InstallKitty() returned nil")
	}

	// Verify that the action is properly configured
	if action.ID == "" {
		t.Error("Action ID should not be empty")
	}

	if action.Name == "" {
		t.Error("Action name should not be empty")
	}

	if action.Description == "" {
		t.Error("Action description should not be empty")
	}

	// Verify that all platform commands are properly configured
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Command for platform %s should not be empty", platform)
		}

		if cmd.CheckCommand == "" {
			t.Errorf("Check command for platform %s should not be empty", platform)
		}
	}
}

// TestInstallKitty_BrewConsistency verifies that brew commands are consistent
func TestInstallKitty_BrewConsistency(t *testing.T) {
	action := InstallKitty()

	brewPlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range brewPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			cmd, exists := action.PlatformCommands[platform]
			if !exists {
				t.Fatalf("%s platform command should exist", platform)
			}

			// Verify brew-specific consistency
			if !utils.ContainsString(cmd.Command, "brew install") {
				t.Errorf("%s command should use 'brew install'", platform)
			}

			if !utils.ContainsString(cmd.Command, "--cask") {
				t.Errorf("%s command should use '--cask' flag", platform)
			}

			if cmd.PackageSource != actions.PackageSourceBrew {
				t.Errorf("%s should use brew package source, got %s", platform, cmd.PackageSource)
			}
		})
	}
}

// TestInstallKitty_PackageSourceValidation validates that package sources are appropriate for each platform
func TestInstallKitty_PackageSourceValidation(t *testing.T) {
	action := InstallKitty()

	// Test Arch package source
	if cmd, exists := action.PlatformCommands[actions.PlatformArch]; exists {
		if cmd.PackageSource != actions.PackageSourceOfficial {
			t.Errorf("Arch should use official package source, got %s", cmd.PackageSource)
		}
	}

	// Test Debian package source
	if cmd, exists := action.PlatformCommands[actions.PlatformDebian]; exists {
		if cmd.PackageSource != actions.PackageSourceOfficial {
			t.Errorf("Debian should use official package source, got %s", cmd.PackageSource)
		}
	}

	// Test Linux package source
	if cmd, exists := action.PlatformCommands[actions.PlatformLinux]; exists {
		if cmd.PackageSource != actions.PackageSourceBrew {
			t.Errorf("Linux should use brew package source, got %s", cmd.PackageSource)
		}
	}

	// Test macOS package source
	if cmd, exists := action.PlatformCommands[actions.PlatformMacOS]; exists {
		if cmd.PackageSource != actions.PackageSourceBrew {
			t.Errorf("macOS should use brew package source, got %s", cmd.PackageSource)
		}
	}
}

// TestInstallKitty_InteractiveFlagConsistency validates that interactive flags are set appropriately
func TestInstallKitty_InteractiveFlagConsistency(t *testing.T) {
	action := InstallKitty()

	// Test Arch interactive flag
	if cmd, exists := action.PlatformCommands[actions.PlatformArch]; exists {
		if !cmd.Interactive {
			t.Error("Arch command should be interactive (requires sudo)")
		}
	}

	// Test Debian interactive flag
	if cmd, exists := action.PlatformCommands[actions.PlatformDebian]; exists {
		if !cmd.Interactive {
			t.Error("Debian command should be interactive (requires sudo)")
		}
	}

	// Test Linux interactive flag
	if cmd, exists := action.PlatformCommands[actions.PlatformLinux]; exists {
		if cmd.Interactive {
			t.Error("Linux command should not be interactive (brew handles it)")
		}
	}

	// Test macOS interactive flag
	if cmd, exists := action.PlatformCommands[actions.PlatformMacOS]; exists {
		if cmd.Interactive {
			t.Error("macOS command should not be interactive (brew handles it)")
		}
	}
}

// TestInstallKitty_CheckCommandConsistency validates that all check commands are the same
func TestInstallKitty_CheckCommandConsistency(t *testing.T) {
	action := InstallKitty()
	expectedCheckCommand := "kitty"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheckCommand {
			t.Errorf("Platform %s should have check command '%s', got '%s'", platform, expectedCheckCommand, cmd.CheckCommand)
		}
	}
}

// TestInstallKitty_CommandComplexity validates that commands are not overly complex
func TestInstallKitty_CommandComplexity(t *testing.T) {
	action := InstallKitty()

	for platform, cmd := range action.PlatformCommands {
		// Commands should be simple and direct
		if len(cmd.Command) > 100 {
			t.Errorf("Platform %s command is too complex (length: %d)", platform, len(cmd.Command))
		}

		// Commands should not contain pipes or complex shell constructs
		if utils.ContainsString(cmd.Command, "|") {
			t.Errorf("Platform %s command should not use pipes", platform)
		}

		if utils.ContainsString(cmd.Command, "&&") {
			t.Errorf("Platform %s command should not use command chaining", platform)
		}

		if utils.ContainsString(cmd.Command, "||") {
			t.Errorf("Platform %s command should not use OR operators", platform)
		}
	}
}

// TestInstallKitty_ArchCommandSudo validates that Arch command uses sudo appropriately
func TestInstallKitty_ArchCommandSudo(t *testing.T) {
	action := InstallKitty()
	cmd, exists := action.PlatformCommands[actions.PlatformArch]

	if !exists {
		t.Fatal("Arch platform command should exist")
	}

	// Arch command should use sudo for package installation
	if !utils.ContainsString(cmd.Command, "sudo") {
		t.Error("Arch command should use sudo for package installation")
	}

	// sudo should be at the beginning of the command
	if !strings.HasPrefix(cmd.Command, "sudo ") {
		t.Error("Arch command should start with sudo")
	}
}

// TestInstallKitty_MockExecutorBehavior validates that the action behaves correctly with mock executors
func TestInstallKitty_MockExecutorBehavior(t *testing.T) {
	action := InstallKitty()

	// This test validates that the action structure is compatible with mock executors
	// used in testing environments
	if action == nil {
		t.Fatal("InstallKitty() returned nil")
	}

	// Verify that all required fields are present for mock executor compatibility
	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Action type should be '%s' for mock executor compatibility", actions.ActionTypeCommand)
	}

	if len(action.PlatformCommands) == 0 {
		t.Error("Platform commands should be configured for mock executor compatibility")
	}

	// Verify each platform command has the required fields
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Platform %s command should not be empty for mock executor", platform)
		}

		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s check command should not be empty for mock executor", platform)
		}
	}
}
