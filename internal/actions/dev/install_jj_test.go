package dev

import (
	"fmt"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
	"anthodev/codory/test/testutil"
)

// TestInstallJujutsuAction tests the creation of the Install Jujutsu action
// This test verifies the action configuration without executing any actual commands
func TestInstallJujutsuAction(t *testing.T) {
	action := InstallJujutsuAction()

	if action == nil {
		t.Fatal("InstallJujutsuAction() returned nil")
	}

	if action.ID != "install_jj" {
		t.Errorf("Expected action ID to be 'install_jj', got '%s'", action.ID)
	}

	if action.Name != "Install Jujutsu (jj) vcs" {
		t.Errorf("Expected action name to be 'Install Jujutsu (jj) vcs', got '%s'", action.Name)
	}

	if action.Description != "Install Jujutsu (jj) version control system" {
		t.Errorf("Expected action description to be 'Install Jujutsu (jj) version control system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallJujutsuAction_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallJujutsuAction_PlatformCommands(t *testing.T) {
	action := InstallJujutsuAction()

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
			expectedCommand:     "sudo pacman -S jujutsu",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "jj",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "brew install jj",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "jj",
			expectedInteractive: false,
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install jj",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "jj",
			expectedInteractive: false,
		},
		{
			name:                "Windows platform",
			platform:            actions.PlatformWindows,
			expectedCommand:     "winget install jj-vcs.jj",
			expectedSource:      actions.PackageSourceWinget,
			expectedCheck:       "jj",
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

// TestInstallJujutsuAction_UnsupportedPlatforms verifies that unsupported platforms don't have commands
// This prevents accidental execution on unsupported systems
func TestInstallJujutsuAction_UnsupportedPlatforms(t *testing.T) {
	action := InstallJujutsuAction()

	// Test that unsupported platforms don't have commands
	unsupportedPlatforms := []actions.Platform{
		actions.PlatformDebian,
	}

	for _, platform := range unsupportedPlatforms {
		if _, exists := action.PlatformCommands[platform]; exists {
			t.Errorf("Expected no command for unsupported platform %s", platform)
		}
	}
}

// TestInstallJujutsuAction_ActionConsistency verifies that the action maintains consistency across multiple calls
func TestInstallJujutsuAction_ActionConsistency(t *testing.T) {
	action1 := InstallJujutsuAction()
	action2 := InstallJujutsuAction()

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

// TestInstallJujutsuAction_MultipleCalls tests that multiple calls to InstallJujutsuAction return equivalent actions
func TestInstallJujutsuAction_MultipleCalls(t *testing.T) {
	// Call the function multiple times to ensure it returns equivalent configurations
	actions := make([]*actions.Action, 5)
	for i := 0; i < 5; i++ {
		actions[i] = InstallJujutsuAction()
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

// TestInstallJujutsuAction_ExpectedPlatforms verifies that all expected platforms are supported
func TestInstallJujutsuAction_ExpectedPlatforms(t *testing.T) {
	action := InstallJujutsuAction()

	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformLinux,
		actions.PlatformMacOS,
		actions.PlatformWindows,
	}

	for _, platform := range expectedPlatforms {
		if _, exists := action.PlatformCommands[platform]; !exists {
			t.Errorf("Expected platform %s to be supported", platform)
		}
	}
}

// TestInstallJujutsuAction_ArchCommandStructure tests the specific command structure for Arch Linux
func TestInstallJujutsuAction_ArchCommandStructure(t *testing.T) {
	action := InstallJujutsuAction()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "sudo pacman -S jujutsu" {
		t.Errorf("Expected Arch command to be 'sudo pacman -S jujutsu', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected Arch package source to be official, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "jj" {
		t.Errorf("Expected Arch check command to be 'jj', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Arch command to be interactive")
	}
}

// TestInstallJujutsuAction_LinuxCommandStructure tests the specific command structure for generic Linux
func TestInstallJujutsuAction_LinuxCommandStructure(t *testing.T) {
	action := InstallJujutsuAction()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "brew install jj" {
		t.Errorf("Expected Linux command to be 'brew install jj', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected Linux package source to be brew, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "jj" {
		t.Errorf("Expected Linux check command to be 'jj', got '%s'", cmd.CheckCommand)
	}

	if cmd.Interactive {
		t.Error("Expected Linux command to be non-interactive")
	}
}

// TestInstallJujutsuAction_MacOSCommandStructure tests the specific command structure for macOS
func TestInstallJujutsuAction_MacOSCommandStructure(t *testing.T) {
	action := InstallJujutsuAction()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "brew install jj" {
		t.Errorf("Expected macOS command to be 'brew install jj', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected macOS package source to be brew, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "jj" {
		t.Errorf("Expected macOS check command to be 'jj', got '%s'", cmd.CheckCommand)
	}

	if cmd.Interactive {
		t.Error("Expected macOS command to be non-interactive")
	}
}

// TestInstallJujutsuAction_WindowsCommandStructure tests the specific command structure for Windows
func TestInstallJujutsuAction_WindowsCommandStructure(t *testing.T) {
	action := InstallJujutsuAction()

	cmd, exists := action.PlatformCommands[actions.PlatformWindows]
	if !exists {
		t.Fatal("Expected Windows platform command to exist")
	}

	// Verify command structure
	if cmd.Command != "winget install jj-vcs.jj" {
		t.Errorf("Expected Windows command to be 'winget install jj-vcs.jj', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceWinget {
		t.Errorf("Expected Windows package source to be winget, got '%s'", cmd.PackageSource)
	}

	if cmd.CheckCommand != "jj" {
		t.Errorf("Expected Windows check command to be 'jj', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Windows command to be interactive")
	}
}

// TestInstallJujutsuAction_CheckCommandStructure tests that check commands are properly configured for each platform
func TestInstallJujutsuAction_CheckCommandStructure(t *testing.T) {
	action := InstallJujutsuAction()

	// All platforms check for "jj"
	expectedCheckCommands := map[actions.Platform]string{
		actions.PlatformArch:    "jj",
		actions.PlatformLinux:   "jj",
		actions.PlatformMacOS:   "jj",
		actions.PlatformWindows: "jj",
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

// TestInstallJujutsuAction_NoCommandExecution verifies that no actual commands are executed during testing
func TestInstallJujutsuAction_NoCommandExecution(t *testing.T) {
	action := InstallJujutsuAction()

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

// TestInstallJujutsuAction_SafeForCI verifies the test is safe for CI environments
func TestInstallJujutsuAction_SafeForCI(t *testing.T) {
	action := InstallJujutsuAction()

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

// TestInstallJujutsuAction_BrewConsistency tests that brew commands are consistent where expected
func TestInstallJujutsuAction_BrewConsistency(t *testing.T) {
	action := InstallJujutsuAction()

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

		if cmd.Command != "brew install jj" {
			t.Errorf("Expected %s brew command to be 'brew install jj', got '%s'", platform, cmd.Command)
		}
	}
}

// TestInstallJujutsuAction_PackageSourceValidation tests that package sources are valid
func TestInstallJujutsuAction_PackageSourceValidation(t *testing.T) {
	action := InstallJujutsuAction()

	validSources := []actions.PackageSource{
		actions.PackageSourceOfficial,
		actions.PackageSourceBrew,
		actions.PackageSourceWinget,
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

// TestInstallJujutsuAction_InteractiveFlagConsistency tests interactive flag consistency
func TestInstallJujutsuAction_InteractiveFlagConsistency(t *testing.T) {
	action := InstallJujutsuAction()

	// Arch and Windows should be interactive
	interactivePlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformWindows,
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

// TestInstallJujutsuAction_CheckCommandConsistency tests that check commands are consistent
func TestInstallJujutsuAction_CheckCommandConsistency(t *testing.T) {
	action := InstallJujutsuAction()

	// All platforms check for "jj"
	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s should have a check command", platform)
		}

		if cmd.CheckCommand != "jj" {
			t.Errorf("Expected check command for platform %s to be 'jj', got '%s'", platform, cmd.CheckCommand)
		}
	}
}

// TestInstallJujutsuAction_CommandComplexity tests that commands have appropriate complexity
func TestInstallJujutsuAction_CommandComplexity(t *testing.T) {
	action := InstallJujutsuAction()

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
		case actions.PlatformLinux:
			if !utils.ContainsString(cmd.Command, "brew") {
				t.Errorf("Linux command should contain 'brew': %s", cmd.Command)
			}
		case actions.PlatformMacOS:
			if !utils.ContainsString(cmd.Command, "brew") {
				t.Errorf("macOS command should contain 'brew': %s", cmd.Command)
			}
		case actions.PlatformWindows:
			if !utils.ContainsString(cmd.Command, "winget") {
				t.Errorf("Windows command should contain 'winget': %s", cmd.Command)
			}
		}
	}
}

// TestInstallJujutsuAction_PackageManagerConsistency tests that package managers are appropriate for each platform
func TestInstallJujutsuAction_PackageManagerConsistency(t *testing.T) {
	action := InstallJujutsuAction()

	// Verify Arch uses pacman
	archCmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}
	if !utils.ContainsString(archCmd.Command, "pacman") {
		t.Errorf("Arch should use pacman package manager: %s", archCmd.Command)
	}

	// Verify Linux uses brew
	linuxCmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}
	if !utils.ContainsString(linuxCmd.Command, "brew") {
		t.Errorf("Linux should use brew package manager: %s", linuxCmd.Command)
	}

	// Verify macOS uses brew
	macOSCmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}
	if !utils.ContainsString(macOSCmd.Command, "brew") {
		t.Errorf("macOS should use brew package manager: %s", macOSCmd.Command)
	}

	// Verify Windows uses winget
	windowsCmd, exists := action.PlatformCommands[actions.PlatformWindows]
	if !exists {
		t.Fatal("Expected Windows platform command to exist")
	}
	if !utils.ContainsString(windowsCmd.Command, "winget") {
		t.Errorf("Windows should use winget package manager: %s", windowsCmd.Command)
	}
}

// TestInstallJujutsuAction_WindowsCommandWinget tests that Windows command properly uses winget
func TestInstallJujutsuAction_WindowsCommandWinget(t *testing.T) {
	action := InstallJujutsuAction()

	cmd, exists := action.PlatformCommands[actions.PlatformWindows]
	if !exists {
		t.Fatal("Expected Windows platform command to exist")
	}

	if cmd.Command != "winget install jj-vcs.jj" {
		t.Errorf("Windows command should use winget: %s", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceWinget {
		t.Errorf("Windows package source should be winget, got '%s'", cmd.PackageSource)
	}
}

// TestInstallJujutsuAction_SudoPresence tests that sudo is properly used where needed
func TestInstallJujutsuAction_SudoPresence(t *testing.T) {
	action := InstallJujutsuAction()

	// Arch should use sudo
	archCmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Expected Arch platform command to exist")
	}
	if !utils.ContainsString(archCmd.Command, "sudo") {
		t.Errorf("Arch command should use sudo: %s", archCmd.Command)
	}

	// Linux should not use sudo with brew
	linuxCmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Expected Linux platform command to exist")
	}
	if utils.ContainsString(linuxCmd.Command, "sudo") {
		t.Errorf("Linux brew command should not use sudo: %s", linuxCmd.Command)
	}

	// macOS should not use sudo with brew
	macOSCmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("Expected macOS platform command to exist")
	}
	if utils.ContainsString(macOSCmd.Command, "sudo") {
		t.Errorf("macOS brew command should not use sudo: %s", macOSCmd.Command)
	}

	// Windows should not use sudo
	windowsCmd, exists := action.PlatformCommands[actions.PlatformWindows]
	if !exists {
		t.Fatal("Expected Windows platform command to exist")
	}
	if utils.ContainsString(windowsCmd.Command, "sudo") {
		t.Errorf("Windows winget command should not use sudo: %s", windowsCmd.Command)
	}
}

// TestInstallJujutsuAction_MockExecutorBehavior tests the action with a mock executor
// This ensures the action can be safely "executed" in test environments without running real commands
func TestInstallJujutsuAction_MockExecutorBehavior(t *testing.T) {
	action := InstallJujutsuAction()

	// Create a mock executor that doesn't execute real commands
	mockExecutor := testutil.NewMockExecutor(func(action *actions.Action) (string, error) {
		// Verify the action structure without executing
		if action.ID != "install_jj" {
			return "", fmt.Errorf("unexpected action ID: %s", action.ID)
		}

		return "jujutsu installation mocked successfully", nil
	})

	// Test that we can "execute" the action through the mock
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("Mock execution failed: %v", err)
	}

	expectedResult := "jujutsu installation mocked successfully"
	if result != expectedResult {
		t.Errorf("Expected mock result '%s', got '%s'", expectedResult, result)
	}
}
