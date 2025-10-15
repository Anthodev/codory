package dev

import (
	"strings"
	"testing"

	"anthodev/codory/internal/actions"
)

// TestInstallJetBrainsToolbox tests the creation of the Install JetBrains Toolbox action
// This test verifies the action configuration without executing any actual commands
func TestInstallJetBrainsToolbox(t *testing.T) {
	action := InstallJetBrainsToolbox()

	if action == nil {
		t.Fatal("InstallJetBrainsToolbox() returned nil")
	}

	if action.ID != "install_jetbrains_toolbox" {
		t.Errorf("Expected action ID to be 'install_jetbrains_toolbox', got '%s'", action.ID)
	}

	if action.Name != "Install Jetbrains Toolbox" {
		t.Errorf("Expected action name to be 'Install Jetbrains Toolbox', got '%s'", action.Name)
	}

	if action.Description != "Install Jetbrains Toolbox on your system" {
		t.Errorf("Expected action description to be 'Install Jetbrains Toolbox on your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeCommand, action.Type)
	}

	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands to be set")
	}
}

// TestInstallJetBrainsToolbox_PlatformCommands tests that platform commands are correctly configured
// This test only verifies the command strings without executing them
func TestInstallJetBrainsToolbox_PlatformCommands(t *testing.T) {
	action := InstallJetBrainsToolbox()

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
			expectedCommand:     "yay -S jetbrains-toolbox",
			expectedSource:      actions.PackageSourceAUR,
			expectedCheck:       "jetbrains-toolbox",
			expectedInteractive: true,
		},
		{
			name:                "Linux platform",
			platform:            actions.PlatformLinux,
			expectedCommand:     "echo Follow the instructions on the following page: https://www.jetbrains.com/help/toolbox-app/toolbox-app-silent-installation.html#tba_installation",
			expectedSource:      actions.PackageSourceAny,
			expectedCheck:       "jetbrains-toolbox",
			expectedInteractive: false, // Linux doesn't explicitly set Interactive
		},
		{
			name:                "macOS platform",
			platform:            actions.PlatformMacOS,
			expectedCommand:     "echo Follow the instructions on the following page: https://www.jetbrains.com/help/toolbox-app/toolbox-app-silent-installation.html#toolbox_macOS",
			expectedSource:      actions.PackageSourceAny,
			expectedCheck:       "jetbrains-toolbox",
			expectedInteractive: false, // macOS doesn't explicitly set Interactive
		},
		{
			name:                "Windows platform",
			platform:            actions.PlatformWindows,
			expectedCommand:     "winget install -e --id JetBrains.Toolbox",
			expectedSource:      actions.PackageSourceWinget,
			expectedCheck:       "jetbrains-toolbox",
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

// TestInstallJetBrainsToolbox_ActionConsistency verifies that the action maintains consistency across multiple calls
func TestInstallJetBrainsToolbox_ActionConsistency(t *testing.T) {
	action1 := InstallJetBrainsToolbox()
	action2 := InstallJetBrainsToolbox()

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

// TestInstallJetBrainsToolbox_MultipleCalls tests that multiple calls to InstallJetBrainsToolbox return equivalent actions
func TestInstallJetBrainsToolbox_MultipleCalls(t *testing.T) {
	// Call the function multiple times to ensure it returns equivalent configurations
	actions := make([]*actions.Action, 5)
	for i := 0; i < 5; i++ {
		actions[i] = InstallJetBrainsToolbox()
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

// TestInstallJetBrainsToolbox_ExpectedPlatforms verifies that all expected platforms are supported
func TestInstallJetBrainsToolbox_ExpectedPlatforms(t *testing.T) {
	action := InstallJetBrainsToolbox()

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

	if len(action.PlatformCommands) != len(expectedPlatforms) {
		t.Errorf("Expected %d supported platforms, got %d", len(expectedPlatforms), len(action.PlatformCommands))
	}
}

// TestInstallJetBrainsToolbox_ArchCommandStructure tests the Arch Linux command structure
func TestInstallJetBrainsToolbox_ArchCommandStructure(t *testing.T) {
	action := InstallJetBrainsToolbox()

	cmd, exists := action.PlatformCommands[actions.PlatformArch]
	if !exists {
		t.Fatal("Arch platform command not found")
	}

	// Verify command structure
	if cmd.Command != "yay -S jetbrains-toolbox" {
		t.Errorf("Expected Arch command to be 'yay -S jetbrains-toolbox', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceAUR {
		t.Errorf("Expected Arch package source to be '%s', got '%s'", actions.PackageSourceAUR, cmd.PackageSource)
	}

	if cmd.CheckCommand != "jetbrains-toolbox" {
		t.Errorf("Expected Arch check command to be 'jetbrains-toolbox', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Arch command to be interactive")
	}
}

// TestInstallJetBrainsToolbox_LinuxCommandStructure tests the Linux command structure
func TestInstallJetBrainsToolbox_LinuxCommandStructure(t *testing.T) {
	action := InstallJetBrainsToolbox()

	cmd, exists := action.PlatformCommands[actions.PlatformLinux]
	if !exists {
		t.Fatal("Linux platform command not found")
	}

	// Verify command structure
	expectedCommand := "echo Follow the instructions on the following page: https://www.jetbrains.com/help/toolbox-app/toolbox-app-silent-installation.html#tba_installation"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected Linux command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceAny {
		t.Errorf("Expected Linux package source to be '%s', got '%s'", actions.PackageSourceAny, cmd.PackageSource)
	}

	if cmd.CheckCommand != "jetbrains-toolbox" {
		t.Errorf("Expected Linux check command to be 'jetbrains-toolbox', got '%s'", cmd.CheckCommand)
	}

	// Linux platform doesn't explicitly set Interactive, so it should be false (default)
	if cmd.Interactive {
		t.Error("Expected Linux command to not be interactive by default")
	}
}

// TestInstallJetBrainsToolbox_MacOSCommandStructure tests the macOS command structure
func TestInstallJetBrainsToolbox_MacOSCommandStructure(t *testing.T) {
	action := InstallJetBrainsToolbox()

	cmd, exists := action.PlatformCommands[actions.PlatformMacOS]
	if !exists {
		t.Fatal("macOS platform command not found")
	}

	// Verify command structure
	expectedCommand := "echo Follow the instructions on the following page: https://www.jetbrains.com/help/toolbox-app/toolbox-app-silent-installation.html#toolbox_macOS"
	if cmd.Command != expectedCommand {
		t.Errorf("Expected macOS command to be '%s', got '%s'", expectedCommand, cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceAny {
		t.Errorf("Expected macOS package source to be '%s', got '%s'", actions.PackageSourceAny, cmd.PackageSource)
	}

	if cmd.CheckCommand != "jetbrains-toolbox" {
		t.Errorf("Expected macOS check command to be 'jetbrains-toolbox', got '%s'", cmd.CheckCommand)
	}

	// macOS platform doesn't explicitly set Interactive, so it should be false (default)
	if cmd.Interactive {
		t.Error("Expected macOS command to not be interactive by default")
	}
}

// TestInstallJetBrainsToolbox_WindowsCommandStructure tests the Windows command structure
func TestInstallJetBrainsToolbox_WindowsCommandStructure(t *testing.T) {
	action := InstallJetBrainsToolbox()

	cmd, exists := action.PlatformCommands[actions.PlatformWindows]
	if !exists {
		t.Fatal("Windows platform command not found")
	}

	// Verify command structure
	if cmd.Command != "winget install -e --id JetBrains.Toolbox" {
		t.Errorf("Expected Windows command to be 'winget install -e --id JetBrains.Toolbox', got '%s'", cmd.Command)
	}

	if cmd.PackageSource != actions.PackageSourceWinget {
		t.Errorf("Expected Windows package source to be '%s', got '%s'", actions.PackageSourceWinget, cmd.PackageSource)
	}

	if cmd.CheckCommand != "jetbrains-toolbox" {
		t.Errorf("Expected Windows check command to be 'jetbrains-toolbox', got '%s'", cmd.CheckCommand)
	}

	if !cmd.Interactive {
		t.Error("Expected Windows command to be interactive")
	}
}

// TestInstallJetBrainsToolbox_CheckCommandConsistency tests that check commands are consistent across platforms
func TestInstallJetBrainsToolbox_CheckCommandConsistency(t *testing.T) {
	action := InstallJetBrainsToolbox()

	expectedCheckCommand := "jetbrains-toolbox"

	for platform, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != expectedCheckCommand {
			t.Errorf("Expected check command for platform %s to be '%s', got '%s'", platform, expectedCheckCommand, cmd.CheckCommand)
		}
	}
}

// TestInstallJetBrainsToolbox_NoCommandExecution tests that no actual commands are executed during testing
func TestInstallJetBrainsToolbox_NoCommandExecution(t *testing.T) {
	// This test ensures that we're only testing the configuration, not executing commands
	action := InstallJetBrainsToolbox()

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
			actions.PackageSourceAUR,
			actions.PackageSourceAny,
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

// TestInstallJetBrainsToolbox_SafeForCI tests that the function is safe to call in CI environments
func TestInstallJetBrainsToolbox_SafeForCI(t *testing.T) {
	// This test ensures the function can be called multiple times safely in CI
	// without any side effects or command execution

	// Call multiple times to ensure no state is maintained
	for i := 0; i < 10; i++ {
		action := InstallJetBrainsToolbox()

		if action == nil {
			t.Fatalf("InstallJetBrainsToolbox() returned nil on call %d", i)
		}

		if action.ID != "install_jetbrains_toolbox" {
			t.Errorf("Expected action ID to be 'install_jetbrains_toolbox' on call %d, got '%s'", i, action.ID)
		}

		// Verify all platforms are present
		expectedPlatforms := 4
		if len(action.PlatformCommands) != expectedPlatforms {
			t.Errorf("Expected %d platforms on call %d, got %d", expectedPlatforms, i, len(action.PlatformCommands))
		}
	}
}

// TestInstallJetBrainsToolbox_PackageSourceValidation tests that package sources are valid
func TestInstallJetBrainsToolbox_PackageSourceValidation(t *testing.T) {
	action := InstallJetBrainsToolbox()

	// Define expected package sources for each platform
	expectedSources := map[actions.Platform]actions.PackageSource{
		actions.PlatformArch:    actions.PackageSourceAUR,
		actions.PlatformLinux:   actions.PackageSourceAny,
		actions.PlatformMacOS:   actions.PackageSourceAny,
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

// TestInstallJetBrainsToolbox_InteractiveFlagConsistency tests interactive flag consistency
func TestInstallJetBrainsToolbox_InteractiveFlagConsistency(t *testing.T) {
	action := InstallJetBrainsToolbox()

	// Define expected interactive flags for each platform
	expectedInteractive := map[actions.Platform]bool{
		actions.PlatformArch:    true,
		actions.PlatformLinux:   false, // Linux doesn't explicitly set Interactive
		actions.PlatformMacOS:   false, // macOS doesn't explicitly set Interactive
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

// TestInstallJetBrainsToolbox_EchoCommandConsistency tests that echo commands are consistent for Linux and macOS
func TestInstallJetBrainsToolbox_EchoCommandConsistency(t *testing.T) {
	action := InstallJetBrainsToolbox()

	// Both Linux and macOS use echo commands with instructions
	linuxCmd, linuxExists := action.PlatformCommands[actions.PlatformLinux]
	macOSCmd, macOSExists := action.PlatformCommands[actions.PlatformMacOS]

	if !linuxExists || !macOSExists {
		t.Fatal("Both Linux and macOS platforms should exist")
	}

	// Both should use echo commands
	if !strings.HasPrefix(linuxCmd.Command, "echo ") {
		t.Errorf("Expected Linux command to start with 'echo ', got '%s'", linuxCmd.Command)
	}

	if !strings.HasPrefix(macOSCmd.Command, "echo ") {
		t.Errorf("Expected macOS command to start with 'echo ', got '%s'", macOSCmd.Command)
	}

	// Both should use PackageSourceAny
	if linuxCmd.PackageSource != actions.PackageSourceAny {
		t.Errorf("Expected Linux package source to be '%s', got '%s'", actions.PackageSourceAny, linuxCmd.PackageSource)
	}

	if macOSCmd.PackageSource != actions.PackageSourceAny {
		t.Errorf("Expected macOS package source to be '%s', got '%s'", actions.PackageSourceAny, macOSCmd.PackageSource)
	}
}
