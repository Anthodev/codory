package dev

import (
	"strings"
	"testing"

	"anthodev/codory/internal/actions"
)

func assertCommandContains(t *testing.T, command string, parts ...string) {
	t.Helper()
	for _, part := range parts {
		if !strings.Contains(command, part) {
			t.Fatalf("command %q missing %q", command, part)
		}
	}
}

// TestInstallVsCode tests the creation of the Install VSCode action
// This test verifies the action configuration without executing any actual commands
func TestInstallVsCode(t *testing.T) {
	action := InstallVsCode()

	if action == nil {
		t.Fatal("InstallVsCode() returned nil")
	}

	if action.ID != "install_vscode" {
		t.Errorf("Expected action ID to be 'install_vscode', got '%s'", action.ID)
	}

	if action.Name != "Install VSCode" {
		t.Errorf("Expected action name to be 'Install VSCode', got '%s'", action.Name)
	}

	if action.Description != "Install Visual Studio Code on your system" {
		t.Errorf("Expected action description to be 'Install Visual Studio Code on your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallVsCode_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallVsCode_PlatformCommands(t *testing.T) {
	action := InstallVsCode()

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
			expectedCommandParts: []string{"pacman", "vscode"},
			expectedSource:       actions.PackageSourceOfficial,
			expectedCheck:        "code",
			expectedInteractive:  true,
		},
		{
			name:                 "Debian platform",
			platform:             actions.PlatformDebian,
			expectedCommandParts: []string{"apt", "code"},
			expectedSource:       actions.PackageSourceOfficial,
			expectedCheck:        "code",
			expectedInteractive:  true,
		},
		{
			name:                 "Linux platform",
			platform:             actions.PlatformLinux,
			expectedCommandParts: []string{"brew", "vscode"},
			expectedSource:       actions.PackageSourceBrew,
			expectedCheck:        "code",
			expectedInteractive:  false, // Note: Linux platform doesn't explicitly set Interactive
		},
		{
			name:                 "macOS platform",
			platform:             actions.PlatformMacOS,
			expectedCommandParts: []string{"brew", "vscode"},
			expectedSource:       actions.PackageSourceBrew,
			expectedCheck:        "code",
			expectedInteractive:  true,
		},
		{
			name:                 "Windows platform",
			platform:             actions.PlatformWindows,
			expectedCommandParts: []string{"winget", "Microsoft.VisualStudioCode"},
			expectedSource:       actions.PackageSourceWinget,
			expectedCheck:        "code",
			expectedInteractive:  true,
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

// TestInstallVsCode_ActionConsistency verifies that the action maintains consistency across multiple calls
func TestInstallVsCode_ActionConsistency(t *testing.T) {
	action1 := InstallVsCode()
	action2 := InstallVsCode()

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

// TestInstallVsCode_MultipleCalls tests that multiple calls to InstallVsCode return equivalent actions
func TestInstallVsCode_MultipleCalls(t *testing.T) {
	// Call the function multiple times to ensure it returns equivalent configurations
	actions := make([]*actions.Action, 5)
	for i := 0; i < 5; i++ {
		actions[i] = InstallVsCode()
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

// TestInstallVsCode_ExpectedPlatforms verifies that all expected platforms are supported
func TestInstallVsCode_ExpectedPlatforms(t *testing.T) {
	action := InstallVsCode()

	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformLinux,
		actions.PlatformMacOS,
		actions.PlatformWindows,
	}

	for _, platform := range expectedPlatforms {
		if _, exists := action.PlatformCommands[platform]; !exists {
			t.Errorf("Expected platform %s to be supported", platform)
		}
	}

	if len(action.PlatformCommands) != len(expectedPlatforms) {
		t.Errorf("Expected %d supported platforms, got %d", len(expectedPlatforms), len(action.PlatformCommands))
	}
}

// TestInstallVsCode_ArchCommandStructure tests the Arch Linux command structure
func TestInstallVsCode_ArchCommandStructure(t *testing.T) {
	action := InstallVsCode()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Arch platform command not found")
	}

	// Verify command structure
	assertCommandContains(t, cmd.Command, "pacman", "vscode")

	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected Arch package source to be '%s', got '%s'", actions.PackageSourceOfficial, cmd.PackageSource)
	}

	if cmd.CheckCommand != "code" {
		t.Errorf("Expected Arch check command to be 'code', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Arch command to be interactive")
	}
}

// TestInstallVsCode_DebianCommandStructure tests the Debian command structure
func TestInstallVsCode_DebianCommandStructure(t *testing.T) {
	action := InstallVsCode()

	cmd, exists := action.PlatformCommands[actions.PlatformDebian]
	if !exists {
		t.Fatal("Debian platform command not found")
	}

	// Verify command structure
	assertCommandContains(t, cmd.Command, "apt", "code")

	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("Expected Debian package source to be '%s', got '%s'", actions.PackageSourceOfficial, cmd.PackageSource)
	}

	if cmd.CheckCommand != "code" {
		t.Errorf("Expected Debian check command to be 'code', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Debian command to be interactive")
	}
}

// TestInstallVsCode_LinuxCommandStructure tests the Linux command structure
func TestInstallVsCode_LinuxCommandStructure(t *testing.T) {
	action := InstallVsCode()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Linux platform command not found")
	}

	// Verify command structure
	assertCommandContains(t, cmd.Command, "brew", "vscode")

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected Linux package source to be '%s', got '%s'", actions.PackageSourceBrew, cmd.PackageSource)
	}

	if cmd.CheckCommand != "code" {
		t.Errorf("Expected Linux check command to be 'code', got '%s'", cmd.CheckCommand)
	}

	// Linux platform doesn't explicitly set Interactive, so it should be false (default)
	if cmd.Interactive {
		t.Error("Expected Linux command to not be interactive by default")
	}
}

// TestInstallVsCode_MacOSCommandStructure tests the macOS command structure
func TestInstallVsCode_MacOSCommandStructure(t *testing.T) {
	action := InstallVsCode()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("macOS platform command not found")
	}

	// Verify command structure
	assertCommandContains(t, cmd.Command, "brew", "vscode")

	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected macOS package source to be '%s', got '%s'", actions.PackageSourceBrew, cmd.PackageSource)
	}

	if cmd.CheckCommand != "code" {
		t.Errorf("Expected macOS check command to be 'code', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected macOS command to be interactive")
	}
}

// TestInstallVsCode_WindowsCommandStructure tests the Windows command structure
func TestInstallVsCode_WindowsCommandStructure(t *testing.T) {
	action := InstallVsCode()

	cmd, exists := action.PlatformCommands[actions.PlatformWindows]
	if !exists {
		t.Fatal("Windows platform command not found")
	}

	// Verify command structure
	assertCommandContains(t, cmd.Command, "winget", "Microsoft.VisualStudioCode")

	if cmd.PackageSource != actions.PackageSourceWinget {
		t.Errorf("Expected Windows package source to be '%s', got '%s'", actions.PackageSourceWinget, cmd.PackageSource)
	}

	if cmd.CheckCommand != "code" {
		t.Errorf("Expected Windows check command to be 'code', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Windows command to be interactive")
	}
}

// TestInstallVsCode_CheckCommandConsistency tests that check commands are consistent across platforms
func TestInstallVsCode_CheckCommandConsistency(t *testing.T) {
	action := InstallVsCode()

	expectedCheckCommand := "code"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheckCommand {
			t.Errorf("Expected check command for platform %s to be '%s', got '%s'", platform, expectedCheckCommand, cmd.CheckCommand)
		}
	}
}

// TestInstallVsCode_NoCommandExecution tests that no actual commands are executed during testing
func TestInstallVsCode_NoCommandExecution(t *testing.T) {
	// This test ensures that we're only testing the configuration, not executing commands
	action := InstallVsCode()

	// Verify that we can access all the command configurations without any execution
	for platform, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("Platform %s has empty command", platform)
		}

		if cmd.CheckCommand == "" {
			t.Errorf("Platform %s has empty check command", platform)
		}

		// Verify package source is valid
		validSources := []actions.PackageSource{
			actions.PackageSourceOfficial,
			actions.PackageSourceBrew,
			actions.PackageSourceWinget,
		}

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

// TestInstallVsCode_SafeForCI tests that the function is safe to call in CI environments
func TestInstallVsCode_SafeForCI(t *testing.T) {
	// This test ensures the function can be called multiple times safely in CI
	// without any side effects or command execution

	// Call multiple times to ensure no state is maintained
	for i := 0; i < 10; i++ {
		action := InstallVsCode()

		if action == nil {
			t.Fatalf("InstallVsCode() returned nil on call %d", i)
		}

		if action.ID != "install_vscode" {
			t.Errorf("Expected action ID to be 'install_vscode' on call %d, got '%s'", i, action.ID)
		}

		// Verify all platforms are present
		expectedPlatforms := 5
		if len(action.PlatformCommands) != expectedPlatforms {
			t.Errorf("Expected %d platforms on call %d, got %d", expectedPlatforms, i, len(action.PlatformCommands))
		}
	}
}

// TestInstallVsCode_BrewConsistency tests that brew commands are consistent where expected
func TestInstallVsCode_BrewConsistency(t *testing.T) {
	action := InstallVsCode()

	// Both Linux and macOS use brew, verify they're consistent
	linuxCmd, linuxExists := action.PlatformCommands[actions.PlatformLinux]
	macOSCmd, macOSExists := action.PlatformCommands[actions.PlatformMacOS]

	if !linuxExists || !macOSExists {
		t.Fatal("Both Linux and macOS platforms should exist")
	}

	// Both should use the same brew command
	expectedBrewCommand := "brew install --cask vscode"
	if linuxCmd.Command != expectedBrewCommand {
		t.Errorf("Expected Linux brew command to be '%s', got '%s'", expectedBrewCommand, linuxCmd.Command)
	}

	if macOSCmd.Command != expectedBrewCommand {
		t.Errorf("Expected macOS brew command to be '%s', got '%s'", expectedBrewCommand, macOSCmd.Command)
	}

	// Both should use brew package source
	if linuxCmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected Linux package source to be '%s', got '%s'", actions.PackageSourceBrew, linuxCmd.PackageSource)
	}

	if macOSCmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("Expected macOS package source to be '%s', got '%s'", actions.PackageSourceBrew, macOSCmd.PackageSource)
	}
}

// TestInstallVsCode_PackageSourceValidation tests that package sources are valid
func TestInstallVsCode_PackageSourceValidation(t *testing.T) {
	action := InstallVsCode()

	// Define expected package sources for each platform
	expectedSources := map[actions.Platform]actions.PackageSource{
		actions.PlatformArch:    actions.PackageSourceOfficial,
		actions.PlatformDebian:  actions.PackageSourceOfficial,
		actions.PlatformLinux:   actions.PackageSourceBrew,
		actions.PlatformMacOS:   actions.PackageSourceBrew,
		actions.PlatformWindows: actions.PackageSourceWinget,
	}

	for platform, expectedSource := range expectedSources {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected platform %s to exist", platform)
			continue
		}

		if cmd.PackageSource != expectedSource {
			t.Errorf("Expected package source for %s to be '%s', got '%s'", platform, expectedSource, cmd.PackageSource)
		}
	}
}

// TestInstallVsCode_InteractiveFlagConsistency tests interactive flag consistency
func TestInstallVsCode_InteractiveFlagConsistency(t *testing.T) {
	action := InstallVsCode()

	// Define expected interactive flags for each platform
	expectedInteractive := map[actions.Platform]bool{
		actions.PlatformArch:    true,
		actions.PlatformDebian:  true,
		actions.PlatformLinux:   false, // Linux doesn't explicitly set Interactive
		actions.PlatformMacOS:   true,
		actions.PlatformWindows: true,
	}

	for platform, expected := range expectedInteractive {
		cmd, exists := action.PlatformCommands[platform]
		if !exists {
			t.Errorf("Expected platform %s to exist", platform)
			continue
		}

		if cmd.Interactive != expected {
			t.Errorf("Expected interactive flag for %s to be %v, got %v", platform, expected, cmd.Interactive)
		}
	}
}
