package terminal

import (
	"strings"
	"testing"

	"anthodev/codory/internal/actions"
)

// TestInstallWezterm tests the creation of the Install Wezterm action
// This test verifies the action configuration without executing any actual commands
func TestInstallWezterm(t *testing.T) {
	action := InstallWezterm()

	if action == nil {
		t.Fatal("InstallWezterm() returned nil")
	}

	if action.ID != "install_wezterm" {
		t.Errorf("Expected action ID to be 'install_wezterm', got '%s'", action.ID)
	}

	if action.Name != "Install Wezterm" {
		t.Errorf("Expected action name to be 'Install Wezterm', got '%s'", action.Name)
	}

	if action.Description != "Install Wezterm, a terminal emulator, on your system" {
		t.Errorf("Expected action description to be 'Install Wezterm, a terminal emulator, on your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallWezterm_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallWezterm_PlatformCommands(t *testing.T) {
	action := InstallWezterm()

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
			expectedCommand:     "sudo pacman -S wezterm",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "wezterm",
			expectedInteractive: true,
		},
		{
			name:                "Debian platform",
			platform:            actions.PlatformDebian,
			expectedCommand:     "curl -fsSL https://apt.fury.io/wez/gpg.key | sudo gpg --yes --dearmor -o /usr/share/keyrings/wezterm-fury.gpg && echo 'deb [signed-by=/usr/share/keyrings/wezterm-fury.gpg] https://apt.fury.io/wez/ * *' | sudo tee /etc/apt/sources.list.d/wezterm.list && sudo chmod 644 /usr/share/keyrings/wezterm-fury.gpg && sudo apt update && sudo apt install wezterm",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "wezterm",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "brew install --cask wezterm",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "wezterm",
			expectedInteractive: false,
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install --cask wezterm",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "wezterm",
			expectedInteractive: false,
		},
		{
			name:                "Windows platform",
			platform:            actions.PlatformWindows,
			expectedCommand:     "winget install -e --id wez.wezterm",
			expectedSource:      actions.PackageSourceWinget,
			expectedCheck:       "wezterm",
			expectedInteractive: true,
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

// TestInstallWezterm_ActionConsistency verifies that the action maintains consistency across multiple calls
func TestInstallWezterm_ActionConsistency(t *testing.T) {
	action1 := InstallWezterm()
	action2 := InstallWezterm()

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

// TestInstallWezterm_MultipleCalls tests that multiple calls to InstallWezterm return equivalent actions
func TestInstallWezterm_MultipleCalls(t *testing.T) {
	// Call the function multiple times to ensure it returns equivalent configurations
	actions := make([]*actions.Action, 5)
	for i := 0; i < 5; i++ {
		actions[i] = InstallWezterm()
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

		if len(actions[0].PlatformCommands) != len(actions[i].PlatformCommands) {
			t.Errorf("Platform commands count should be consistent across calls: %d != %d", len(actions[0].PlatformCommands), len(actions[i].PlatformCommands))
		}
	}
}

// TestInstallWezterm_ExpectedPlatforms verifies that all expected platforms are supported
func TestInstallWezterm_ExpectedPlatforms(t *testing.T) {
	action := InstallWezterm()

	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformLinux,
		actions.PlatformMacOS,
		actions.PlatformWindows,
	}

	for _, platform := range expectedPlatforms {
		if _, exists := action.PlatformCommands[platform]; !exists {
			t.Errorf("Expected command for platform %s to exist", platform)
		}
	}
}

// TestInstallWezterm_ArchCommandStructure tests Arch-specific command details
func TestInstallWezterm_ArchCommandStructure(t *testing.T) {
	action := InstallWezterm()
	cmd := action.PlatformCommands[actions.PlatformArch]

	if !strings.Contains(cmd.Command, "pacman") {
		t.Errorf("Arch command should use pacman, got: %s", cmd.Command)
	}

	if !strings.Contains(cmd.Command, "wezterm") {
		t.Errorf("Arch command should reference wezterm, got: %s", cmd.Command)
	}

	if !strings.Contains(cmd.Command, "sudo") {
		t.Errorf("Arch command should include sudo, got: %s", cmd.Command)
	}
}

// TestInstallWezterm_DebianCommandStructure tests Debian-specific command details
func TestInstallWezterm_DebianCommandStructure(t *testing.T) {
	action := InstallWezterm()
	cmd := action.PlatformCommands[actions.PlatformDebian]

	if !strings.Contains(cmd.Command, "curl") {
		t.Errorf("Debian command should use curl, got: %s", cmd.Command)
	}

	if !strings.Contains(cmd.Command, "wezterm") {
		t.Errorf("Debian command should reference wezterm, got: %s", cmd.Command)
	}

	if !strings.Contains(cmd.Command, "apt") {
		t.Errorf("Debian command should use apt, got: %s", cmd.Command)
	}
}

// TestInstallWezterm_LinuxCommandStructure tests Linux-specific command details
func TestInstallWezterm_LinuxCommandStructure(t *testing.T) {
	action := InstallWezterm()
	cmd := action.PlatformCommands[actions.PlatformLinux]

	if !strings.Contains(cmd.Command, "brew") {
		t.Errorf("Linux command should use brew, got: %s", cmd.Command)
	}

	if !strings.Contains(cmd.Command, "wezterm") {
		t.Errorf("Linux command should reference wezterm, got: %s", cmd.Command)
	}
}

// TestInstallWezterm_MacOSCommandStructure tests macOS-specific command details
func TestInstallWezterm_MacOSCommandStructure(t *testing.T) {
	action := InstallWezterm()
	cmd := action.PlatformCommands[actions.PlatformMacOS]

	if !strings.Contains(cmd.Command, "brew") {
		t.Errorf("macOS command should use brew, got: %s", cmd.Command)
	}

	if !strings.Contains(cmd.Command, "cask") {
		t.Errorf("macOS command should use cask, got: %s", cmd.Command)
	}

	if !strings.Contains(cmd.Command, "wezterm") {
		t.Errorf("macOS command should reference wezterm, got: %s", cmd.Command)
	}
}

// TestInstallWezterm_WindowsCommandStructure tests Windows-specific command details
func TestInstallWezterm_WindowsCommandStructure(t *testing.T) {
	action := InstallWezterm()
	cmd := action.PlatformCommands[actions.PlatformWindows]

	if !strings.Contains(cmd.Command, "winget") {
		t.Errorf("Windows command should use winget, got: %s", cmd.Command)
	}

	if !strings.Contains(cmd.Command, "wez.wezterm") {
		t.Errorf("Windows command should reference wez.wezterm, got: %s", cmd.Command)
	}
}

// TestInstallWezterm_CheckCommandStructure verifies that all check commands are consistent
func TestInstallWezterm_CheckCommandStructure(t *testing.T) {
	action := InstallWezterm()

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != "wezterm" {
			t.Errorf("Check command for platform %s should be 'wezterm', got '%s'", platform, cmd.CheckCommand)
		}
	}
}

// TestInstallWezterm_NoCommandExecution ensures that the action only returns command strings
// without actually executing them
func TestInstallWezterm_NoCommandExecution(t *testing.T) {
	action := InstallWezterm()

	for platform, cmd := range action.PlatformCommands {
		// Verify command is a string, not an execution result
		if cmd.Command == "" {
			t.Errorf("Command for platform %s should not be empty", platform)
		}

		// Verify the command is stored as expected, not executed
		if cmd.Command != strings.TrimSpace(cmd.Command) {
			t.Errorf("Command for platform %s should be trimmed", platform)
		}
	}
}

// TestInstallWezterm_SafeForCI ensures that commands are suitable for CI environments
func TestInstallWezterm_SafeForCI(t *testing.T) {
	action := InstallWezterm()

	// Verify that commands don't have hardcoded user-specific paths
	unsafePatterns := []string{"/home/", "/Users/", "~"}

	for platform, cmd := range action.PlatformCommands {
		for _, pattern := range unsafePatterns {
			if strings.Contains(cmd.Command, pattern) {
				t.Errorf("Command for platform %s contains unsafe pattern %s: %s", platform, pattern, cmd.Command)
			}
		}
	}
}

// TestInstallWezterm_BrewConsistency verifies that brew commands use the same format
func TestInstallWezterm_BrewConsistency(t *testing.T) {
	action := InstallWezterm()

	brewPlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range brewPlatforms {
		cmd := action.PlatformCommands[platform]
		if !strings.Contains(cmd.Command, "brew install --cask wezterm") {
			t.Errorf("Brew command for platform %s should be consistent, got: %s", platform, cmd.Command)
		}

		if cmd.PackageSource != actions.PackageSourceBrew {
			t.Errorf("Package source for platform %s should be Brew, got: %s", platform, cmd.PackageSource)
		}
	}
}

// TestInstallWezterm_PackageSourceValidation verifies that package sources are appropriate for each platform
func TestInstallWezterm_PackageSourceValidation(t *testing.T) {
	action := InstallWezterm()

	tests := []struct {
		platform       actions.Platform
		expectedSource actions.PackageSource
	}{
		{actions.PlatformArch, actions.PackageSourceOfficial},
		{actions.PlatformDebian, actions.PackageSourceOfficial},
		{actions.PlatformLinux, actions.PackageSourceBrew},
		{actions.PlatformMacOS, actions.PackageSourceBrew},
		{actions.PlatformWindows, actions.PackageSourceWinget},
	}

	for _, tt := range tests {
		cmd := action.PlatformCommands[tt.platform]
		if cmd.PackageSource != tt.expectedSource {
			t.Errorf("Package source for platform %s should be %s, got %s", tt.platform, tt.expectedSource, cmd.PackageSource)
		}
	}
}

// TestInstallWezterm_InteractiveFlagConsistency verifies that interactive flags are set appropriately
func TestInstallWezterm_InteractiveFlagConsistency(t *testing.T) {
	action := InstallWezterm()

	// Platforms that should be interactive (require user confirmation or input)
	interactivePlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformWindows,
	}

	// Platforms that should not be interactive
	nonInteractivePlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range interactivePlatforms {
		cmd := action.PlatformCommands[platform]
		if !cmd.Interactive {
			t.Errorf("Platform %s should be interactive", platform)
		}
	}

	for _, platform := range nonInteractivePlatforms {
		cmd := action.PlatformCommands[platform]
		if cmd.Interactive {
			t.Errorf("Platform %s should not be interactive", platform)
		}
	}
}

// TestInstallWezterm_CheckCommandConsistency verifies that check commands are consistent
func TestInstallWezterm_CheckCommandConsistency(t *testing.T) {
	action := InstallWezterm()

	expectedCheckCmd := "wezterm"
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheckCmd {
			t.Errorf("Check command for platform %s should be '%s', got '%s'", platform, expectedCheckCmd, cmd.CheckCommand)
		}
	}
}

// TestInstallWezterm_CommandComplexity tests that commands have reasonable complexity
func TestInstallWezterm_CommandComplexity(t *testing.T) {
	action := InstallWezterm()

	for platform, cmd := range action.PlatformCommands {
		// Verify command is not empty
		if len(cmd.Command) == 0 {
			t.Errorf("Command for platform %s should not be empty", platform)
		}

		// Verify command length is reasonable (not excessively long or short)
		if len(cmd.Command) < 5 {
			t.Errorf("Command for platform %s seems too short: %s", platform, cmd.Command)
		}

		// Most commands should be less than 1000 characters
		if len(cmd.Command) > 2000 {
			t.Errorf("Command for platform %s seems too long (%d characters)", platform, len(cmd.Command))
		}
	}
}

// TestInstallWezterm_ArchCommandSudo verifies that Arch command uses sudo
func TestInstallWezterm_ArchCommandSudo(t *testing.T) {
	action := InstallWezterm()
	cmd := action.PlatformCommands[actions.PlatformArch]

	if !strings.HasPrefix(cmd.Command, "sudo") {
		t.Errorf("Arch command should start with 'sudo', got: %s", cmd.Command)
	}
}

// TestInstallWezterm_DebianCommandMultiStep verifies that Debian command has multiple steps
func TestInstallWezterm_DebianCommandMultiStep(t *testing.T) {
	action := InstallWezterm()
	cmd := action.PlatformCommands[actions.PlatformDebian]

	// Debian command should have multiple steps connected with &&
	if !strings.Contains(cmd.Command, "&&") {
		t.Errorf("Debian command should have multiple steps connected with &&, got: %s", cmd.Command)
	}

	// Count the number of steps (should be at least 3 for gpg, apt, etc.)
	steps := strings.Count(cmd.Command, "&&")
	if steps < 2 {
		t.Errorf("Debian command should have multiple steps, got only %d separators", steps)
	}
}
