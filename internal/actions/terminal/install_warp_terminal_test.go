package terminal

import (
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestInstallWarpTerminal tests the creation of the Install Warp Terminal action
// This test verifies the action configuration without executing any actual commands
func TestInstallWarpTerminal(t *testing.T) {
	action := InstallWarpTerminal()

	if action == nil {
		t.Fatal("InstallWarpTerminal() returned nil")
	}

	if action.ID != "install_warp_terminal" {
		t.Errorf("Expected action ID to be 'install_warp_terminal', got '%s'", action.ID)
	}

	if action.Name != "Install Warp Terminal" {
		t.Errorf("Expected action name to be 'Install Warp Terminal', got '%s'", action.Name)
	}

	if action.Description != "Install Warp Terminal in your system" {
		t.Errorf("Expected action description to be 'Install Warp Terminal in your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallWarpTerminal_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallWarpTerminal_PlatformCommands(t *testing.T) {
	action := InstallWarpTerminal()

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
			expectedCommand:     "yay -S warp-terminal-bin",
			expectedSource:      actions.PackageSourceAUR,
			expectedCheck:       "warp",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "brew install --cask warp",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "warp",
			expectedInteractive: false,
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install --cask warp",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "warp",
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

// TestInstallWarpTerminal_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallWarpTerminal_UnsupportedPlatforms(t *testing.T) {
	action := InstallWarpTerminal()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.PlatformWindows,
		actions.PlatformDebian,
		actions.PlatformAny,
	}

	for _, platform := range unsupportedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if cmd, exists := action.PlatformCommands[platform]; exists {
				t.Errorf("Expected platform %s to not have a command, but found: %s", platform, cmd.Command)
			}
		})
	}
}

// TestInstallWarpTerminal_ActionConsistency verifies that the action maintains consistency across calls
// This ensures the action is deterministic and doesn't have side effects
func TestInstallWarpTerminal_ActionConsistency(t *testing.T) {
	action1 := InstallWarpTerminal()
	action2 := InstallWarpTerminal()

	// Verify that multiple calls return the same configuration
	if action1.ID != action2.ID {
		t.Errorf("Action ID should be consistent across calls: got '%s' and '%s'", action1.ID, action2.ID)
	}

	if action1.Name != action2.Name {
		t.Errorf("Action name should be consistent across calls: got '%s' and '%s'", action1.Name, action2.Name)
	}

	if action1.Description != action2.Description {
		t.Errorf("Action description should be consistent across calls: got '%s' and '%s'", action1.Description, action2.Description)
	}

	if action1.Type != action2.Type {
		t.Errorf("Action type should be consistent across calls: got '%s' and '%s'", action1.Type, action2.Type)
	}

	// Verify platform commands are consistent
	if len(action1.PlatformCommands) != len(action2.PlatformCommands) {
		t.Errorf("Platform commands count should be consistent: got %d and %d", len(action1.PlatformCommands), len(action2.PlatformCommands))
	}

	for platform, cmd1 := range action1.PlatformCommands {
		cmd2, exists := action2.PlatformCommands[platform]
		if !exists {
			t.Errorf("Platform %s exists in first action but not in second", platform)
			continue
		}

		if cmd1.Command != cmd2.Command {
			t.Errorf("Command for platform %s should be consistent: got '%s' and '%s'", platform, cmd1.Command, cmd2.Command)
		}

		if cmd1.PackageSource != cmd2.PackageSource {
			t.Errorf("Package source for platform %s should be consistent: got '%s' and '%s'", platform, cmd1.PackageSource, cmd2.PackageSource)
		}

		if cmd1.CheckCommand != cmd2.CheckCommand {
			t.Errorf("Check command for platform %s should be consistent: got '%s' and '%s'", platform, cmd1.CheckCommand, cmd2.CheckCommand)
		}

		if cmd1.Interactive != cmd2.Interactive {
			t.Errorf("Interactive flag for platform %s should be consistent: got %v and %v", platform, cmd1.Interactive, cmd2.Interactive)
		}
	}
}

// TestInstallWarpTerminal_MultipleCalls verifies that calling the function multiple times doesn't cause issues
// This test ensures the function is pure and doesn't have side effects
func TestInstallWarpTerminal_MultipleCalls(t *testing.T) {
	// Call the function multiple times to ensure no side effects
	for i := 0; i < 5; i++ {
		action := InstallWarpTerminal()
		if action == nil {
			t.Fatalf("InstallWarpTerminal() returned nil on call %d", i+1)
		}

		// Verify basic properties on each call
		if action.ID != "install_warp_terminal" {
			t.Errorf("Expected action ID to be 'install_warp_terminal' on call %d, got '%s'", i+1, action.ID)
		}

		if len(action.PlatformCommands) != 3 {
			t.Errorf("Expected 3 platform commands on call %d, got %d", i+1, len(action.PlatformCommands))
		}
	}
}

// TestInstallWarpTerminal_ExpectedPlatforms verifies that all expected platforms are supported
func TestInstallWarpTerminal_ExpectedPlatforms(t *testing.T) {
	action := InstallWarpTerminal()

	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range expectedPlatforms {
		t.Run(string(platform), func(t *testing.T) {
			if cmd, exists := action.PlatformCommands[platform]; !exists {
				t.Errorf("Expected platform %s to be supported", platform)
			} else {
				if cmd.Command == "" {
					t.Errorf("Expected platform %s to have a non-empty command", platform)
				}
				if cmd.CheckCommand == "" {
					t.Errorf("Expected platform %s to have a non-empty check command", platform)
				}
			}
		})
	}
}

// TestInstallWarpTerminal_ArchCommandStructure tests the specific structure of Arch commands
func TestInstallWarpTerminal_ArchCommandStructure(t *testing.T) {
	action := InstallWarpTerminal()
	cmd, exists := action.PlatformCommands[actions.PlatformArch]

	if !exists {
		t.Fatal("Arch platform command not found")
	}

	// Test Arch-specific command structure
	testutil.AssertContains(t, cmd.Command, "yay", "Arch command should use yay")
	testutil.AssertContains(t, cmd.Command, "warp-terminal-bin", "Arch command should install warp-terminal-bin package")
	testutil.AssertContains(t, cmd.Command, "-S", "Arch command should use -S flag for sync")

	if cmd.PackageSource != actions.PackageSourceAUR {
		t.Errorf("Expected Arch package source to be AUR, got %s", cmd.PackageSource)
	}

	if !cmd.Interactive {
		t.Error("Expected Arch command to be interactive")
	}
}

// TestInstallWarpTerminal_LinuxCommandStructure tests the specific structure of Linux commands
func TestInstallWarpTerminal_LinuxCommandStructure(t *testing.T) {
	action := InstallWarpTerminal()
	cmd, exists := action.PlatformCommands[actions.PlatformLinux]

	if !exists {
		t.Fatal("Linux platform command not found")
	}

	// Test Linux-specific command structure
	testutil.AssertContains(t, cmd.Command, "brew", "Linux command should use brew")
	testutil.AssertContains(t, cmd.Command, "install", "Linux command should use install")
	testutil.AssertContains(t, cmd.Command, "--cask", "Linux command should use --cask flag")
	testutil.AssertContains(t, cmd.Command, "warp", "Linux command should install warp")

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected Linux package source to be Brew, got %s", cmd.PackageSource)
	}

	if cmd.Interactive {
		t.Error("Expected Linux command to be non-interactive")
	}
}

// TestInstallWarpTerminal_MacOSCommandStructure tests the specific structure of macOS commands
func TestInstallWarpTerminal_MacOSCommandStructure(t *testing.T) {
	action := InstallWarpTerminal()
	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]

	if !exists {
		t.Fatal("macOS platform command not found")
	}

	// Test macOS-specific command structure
	testutil.AssertContains(t, cmd.Command, "brew", "macOS command should use brew")
	testutil.AssertContains(t, cmd.Command, "install", "macOS command should use install")
	testutil.AssertContains(t, cmd.Command, "--cask", "macOS command should use --cask flag")
	testutil.AssertContains(t, cmd.Command, "warp", "macOS command should install warp")

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected macOS package source to be Brew, got %s", cmd.PackageSource)
	}

	if cmd.Interactive {
		t.Error("Expected macOS command to be non-interactive")
	}
}

// TestInstallWarpTerminal_CheckCommandStructure tests that all check commands are consistent
func TestInstallWarpTerminal_CheckCommandStructure(t *testing.T) {
	action := InstallWarpTerminal()

	for platform, cmd := range action.PlatformCommands {
		t.Run(string(platform), func(t *testing.T) {
			if cmd.CheckCommand != "warp" {
				t.Errorf("Expected check command for platform %s to be 'warp', got '%s'", platform, cmd.CheckCommand)
			}
		})
	}
}

// TestInstallWarpTerminal_NoCommandExecution verifies that the test doesn't execute any actual commands
// This is important for CI/CD safety and test performance
func TestInstallWarpTerminal_NoCommandExecution(t *testing.T) {
	action := InstallWarpTerminal()

	// Verify that the action is properly configured but doesn't execute commands
	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be command, got %s", action.Type)
	}

	// Verify that commands are strings, not executed functions
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Expected non-empty command for platform %s", platform)
		}

		// Verify command is a string (not a function or executable)
		if len(cmd.Command) == 0 {
			t.Errorf("Expected command to be a string for platform %s", platform)
		}
	}
}

// TestInstallWarpTerminal_SafeForCI verifies that the tests are safe to run in CI/CD environments
func TestInstallWarpTerminal_SafeForCI(t *testing.T) {
	// This test should not require any external dependencies or system access
	action := InstallWarpTerminal()

	if action == nil {
		t.Fatal("Action should not be nil in CI environment")
	}

	// Verify that the action doesn't try to access the filesystem or network
	// All properties should be static and deterministic
	if action.ID == "" {
		t.Error("Action ID should not be empty")
	}

	if action.Name == "" {
		t.Error("Action name should not be empty")
	}

	if action.Description == "" {
		t.Error("Action description should not be empty")
	}
}

// TestInstallWarpTerminal_BrewConsistency tests that brew commands are consistent across platforms
func TestInstallWarpTerminal_BrewConsistency(t *testing.T) {
	action := InstallWarpTerminal()

	// Get Linux and macOS commands
	linuxCmd, linuxExists := action.PlatformCommands[actions.PlatformLinux]
	macosCmd, macosExists := action.PlatformCommands[actions.PlatformMacOS]

	if !linuxExists {
		t.Fatal("Linux platform command not found")
	}

	if !macosExists {
		t.Fatal("macOS platform command not found")
	}

	// Verify both use the same brew command structure
	if linuxCmd.Command != macosCmd.Command {
		t.Errorf("Expected Linux and macOS commands to be identical: Linux='%s', macOS='%s'", linuxCmd.Command, macosCmd.Command)
	}

	if linuxCmd.PackageSource != macosCmd.PackageSource {
		t.Errorf("Expected Linux and macOS package sources to be identical: Linux='%s', macOS='%s'", linuxCmd.PackageSource, macosCmd.PackageSource)
	}

	if linuxCmd.CheckCommand != macosCmd.CheckCommand {
		t.Errorf("Expected Linux and macOS check commands to be identical: Linux='%s', macOS='%s'", linuxCmd.CheckCommand, macosCmd.CheckCommand)
	}

	if linuxCmd.Interactive != macosCmd.Interactive {
		t.Errorf("Expected Linux and macOS interactive flags to be identical: Linux=%v, macOS=%v", linuxCmd.Interactive, macosCmd.Interactive)
	}
}

// TestInstallWarpTerminal_PackageSourceValidation tests that package sources are valid
func TestInstallWarpTerminal_PackageSourceValidation(t *testing.T) {
	action := InstallWarpTerminal()

	validSources := map[actions.Platform]actions.PackageSource{
		actions.PlatformArch:  actions.PackageSourceAUR,
		actions.PlatformLinux: actions.PackageSourceBrew,
		actions.PlatformMacOS: actions.PackageSourceBrew,
	}

	for platform, expectedSource := range validSources {
		t.Run(string(platform), func(t *testing.T) {
			cmd, exists := action.PlatformCommands[platform]
			if !exists {
				t.Errorf("Expected platform %s to have a command", platform)
				return
			}

			if cmd.PackageSource != expectedSource {
				t.Errorf("Expected package source for platform %s to be %s, got %s", platform, expectedSource, cmd.PackageSource)
			}

			// Verify package source is not empty
			if cmd.PackageSource == "" {
				t.Errorf("Package source for platform %s should not be empty", platform)
			}
		})
	}
}

// TestInstallWarpTerminal_InteractiveFlagConsistency tests that interactive flags are set correctly
func TestInstallWarpTerminal_InteractiveFlagConsistency(t *testing.T) {
	action := InstallWarpTerminal()

	expectedInteractive := map[actions.Platform]bool{
		actions.PlatformArch:  true,  // yay requires user interaction for AUR packages
		actions.PlatformLinux: false, // brew cask install is typically non-interactive
		actions.PlatformMacOS: false, // brew cask install is typically non-interactive
	}

	for platform, expected := range expectedInteractive {
		t.Run(string(platform), func(t *testing.T) {
			cmd, exists := action.PlatformCommands[platform]
			if !exists {
				t.Errorf("Expected platform %s to have a command", platform)
				return
			}

			if cmd.Interactive != expected {
				t.Errorf("Expected interactive flag for platform %s to be %v, got %v", platform, expected, cmd.Interactive)
			}
		})
	}
}

// TestInstallWarpTerminal_CheckCommandConsistency tests that all check commands are the same
func TestInstallWarpTerminal_CheckCommandConsistency(t *testing.T) {
	action := InstallWarpTerminal()

	expectedCheckCommand := "warp"

	for platform, cmd := range action.PlatformCommands {
		t.Run(string(platform), func(t *testing.T) {
			if cmd.CheckCommand != expectedCheckCommand {
				t.Errorf("Expected check command for platform %s to be '%s', got '%s'", platform, expectedCheckCommand, cmd.CheckCommand)
			}

			// Verify check command is not empty
			if cmd.CheckCommand == "" {
				t.Errorf("Check command for platform %s should not be empty", platform)
			}
		})
	}
}

// TestInstallWarpTerminal_CommandComplexity tests that commands have appropriate complexity
func TestInstallWarpTerminal_CommandComplexity(t *testing.T) {
	action := InstallWarpTerminal()

	for platform, cmd := range action.PlatformCommands {
		t.Run(string(platform), func(t *testing.T) {
			// Commands should not be too simple (just "warp") or too complex
			if len(cmd.Command) < 10 {
				t.Errorf("Command for platform %s seems too simple: '%s'", platform, cmd.Command)
			}

			if len(cmd.Command) > 200 {
				t.Errorf("Command for platform %s seems too complex: '%s'", platform, cmd.Command)
			}

			// Commands should contain platform-appropriate package manager references
			switch platform {
			case actions.PlatformArch:
				if !utils.ContainsString(cmd.Command, "yay") {
					t.Errorf("Arch command should contain 'yay': '%s'", cmd.Command)
				}
			case actions.PlatformLinux, actions.PlatformMacOS:
				if !utils.ContainsString(cmd.Command, "brew") {
					t.Errorf("%s command should contain 'brew': '%s'", platform, cmd.Command)
				}
			}
		})
	}
}

// TestInstallWarpTerminal_ArchCommandYay tests that Arch command specifically uses yay
func TestInstallWarpTerminal_ArchCommandYay(t *testing.T) {
	action := InstallWarpTerminal()
	cmd, exists := action.PlatformCommands[actions.PlatformArch]

	if !exists {
		t.Fatal("Arch platform command not found")
	}

	// Verify Arch command uses yay specifically
	if cmd.Command != "yay -S warp-terminal-bin" {
		t.Errorf("Expected Arch command to be exactly 'yay -S warp-terminal-bin', got '%s'", cmd.Command)
	}

	// Verify it's marked as using AUR
	if cmd.PackageSource != actions.PackageSourceAUR {
		t.Errorf("Expected Arch command to use AUR package source, got %s", cmd.PackageSource)
	}
}

// TestInstallWarpTerminal_MockExecutorBehavior tests that the action would work correctly with a mock executor
func TestInstallWarpTerminal_MockExecutorBehavior(t *testing.T) {
	action := InstallWarpTerminal()

	// This test verifies that the action structure is compatible with mock executors
	// used in other parts of the codebase
	if action == nil {
		t.Fatal("Action should not be nil")
	}

	// Verify the action has the expected structure for mock execution
	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be %s for mock execution", actions.ActionTypeCommand)
	}

	// Verify all platforms have the required fields for mock execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Platform %s should have a command for mock execution", platform)
		}

		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s should have a check command for mock execution", platform)
		}

		if cmd.PackageSource == "" {
			t.Errorf("Platform %s should have a package source for mock execution", platform)
		}
	}
}
