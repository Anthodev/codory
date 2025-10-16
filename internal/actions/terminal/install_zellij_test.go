package terminal

import (
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
)

// TestInstallZellij verifies the basic fields of the InstallZellij action.
func TestInstallZellij(t *testing.T) {
	action := InstallZellij()

	if action == nil {
		t.Fatal("InstallZellij() returned nil")
	}
	if action.ID != "install_zellij" {
		t.Errorf("expected ID 'install_zellij', got '%s'", action.ID)
	}
	if action.Name != "Install Zellij" {
		t.Errorf("expected Name 'Install Zellij', got '%s'", action.Name)
	}
	if action.Description != "Install Zellij, a terminal multiplexer, on your system" {
		t.Errorf("unexpected description: %s", action.Description)
	}
	if action.Type != actions.ActionTypeCommand {
		t.Errorf("expected Type %s, got %s", actions.ActionTypeCommand, action.Type)
	}
	if action.PlatformCommands == nil {
		t.Fatal("PlatformCommands map is nil")
	}
}

// TestInstallZellij_PlatformCommands validates command configuration per platform.
func TestInstallZellij_PlatformCommands(t *testing.T) {
	action := InstallZellij()
	tests := []struct {
		platform            actions.Platform
		expectedCommand     string
		expectedSource      actions.PackageSource
		expectedCheck       string
		expectedInteractive bool
	}{
		{
			platform:            actions.PlatformArch,
			expectedCommand:     "sudo pacman -S zellij",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "zellij",
			expectedInteractive: true,
		},
		{
			platform:            actions.PlatformLinux,
			expectedCommand:     "brew installl zellij",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "zellij",
			expectedInteractive: false,
		},
		{
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install zellij",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "zellij",
			expectedInteractive: false,
		},
	}
	for _, tt := range tests {
		t.Run(string(tt.platform), func(t *testing.T) {
			cmd, ok := action.PlatformCommands[tt.platform]
			if !ok {
				t.Fatalf("platform %s missing", tt.platform)
			}
			if cmd.Command != tt.expectedCommand {
				t.Errorf("command mismatch for %s: expected %q, got %q", tt.platform, tt.expectedCommand, cmd.Command)
			}
			if cmd.PackageSource != tt.expectedSource {
				t.Errorf("package source mismatch for %s: expected %s, got %s", tt.platform, tt.expectedSource, cmd.PackageSource)
			}
			if cmd.CheckCommand != tt.expectedCheck {
				t.Errorf("check command mismatch for %s: expected %s, got %s", tt.platform, tt.expectedCheck, cmd.CheckCommand)
			}
			if cmd.Interactive != tt.expectedInteractive {
				t.Errorf("interactive flag mismatch for %s: expected %v, got %v", tt.platform, tt.expectedInteractive, cmd.Interactive)
			}
		})
	}
}

// TestInstallZellij_UnsupportedPlatforms ensures unsupported platforms have no entry.
func TestInstallZellij_UnsupportedPlatforms(t *testing.T) {
	action := InstallZellij()
	unsupported := []actions.Platform{actions.PlatformWindows, actions.PlatformDebian}
	for _, p := range unsupported {
		if _, ok := action.PlatformCommands[p]; ok {
			t.Errorf("unexpected command for unsupported platform %s", p)
		}
	}
}

// TestInstallZellij_ActionConsistency checks that multiple calls return identical actions.
func TestInstallZellij_ActionConsistency(t *testing.T) {
	a1 := InstallZellij()
	a2 := InstallZellij()

	if a1.ID != a2.ID || a1.Name != a2.Name || a1.Description != a2.Description || a1.Type != a2.Type {
		t.Error("basic fields differ between calls")
	}
	if len(a1.PlatformCommands) != len(a2.PlatformCommands) {
		t.Error("platform command count differs")
	}
	for p, c1 := range a1.PlatformCommands {
		c2, ok := a2.PlatformCommands[p]
		if !ok {
			t.Errorf("platform %s missing in second call", p)
			continue
		}
		if c1 != c2 {
			t.Errorf("platform %s command differs between calls", p)
		}
	}
}

// TestInstallZellij_MultipleCalls tests that multiple calls return equivalent actions.
func TestInstallZellij_MultipleCalls(t *testing.T) {
	actions := make([]*actions.Action, 5)
	for i := 0; i < 5; i++ {
		actions[i] = InstallZellij()
	}

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

// TestInstallZellij_ExpectedPlatforms verifies all expected platforms are present.
func TestInstallZellij_ExpectedPlatforms(t *testing.T) {
	action := InstallZellij()
	expected := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}
	for _, p := range expected {
		if _, ok := action.PlatformCommands[p]; !ok {
			t.Errorf("expected platform %s to be present", p)
		}
	}
}

// TestInstallZellij_ArchCommandStructure validates the Arch command specifics.
func TestInstallZellij_ArchCommandStructure(t *testing.T) {
	action := InstallZellij()
	cmd, ok := action.PlatformCommands[actions.PlatformArch]
	if !ok {
		t.Fatal("Arch platform missing")
	}
	if cmd.Command != "sudo pacman -S zellij" {
		t.Errorf("unexpected Arch command: %s", cmd.Command)
	}
	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("unexpected Arch package source: %s", cmd.PackageSource)
	}
	if cmd.CheckCommand != "zellij" {
		t.Errorf("unexpected Arch check command: %s", cmd.CheckCommand)
	}
	if !cmd.Interactive {
		t.Error("Arch command should be interactive")
	}
}

// TestInstallZellij_LinuxCommandStructure validates the generic Linux command specifics.
func TestInstallZellij_LinuxCommandStructure(t *testing.T) {
	action := InstallZellij()
	cmd, ok := action.PlatformCommands[actions.PlatformLinux]
	if !ok {
		t.Fatal("Linux platform missing")
	}
	if cmd.Command != "brew installl zellij" {
		t.Errorf("unexpected Linux command: %s", cmd.Command)
	}
	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("unexpected Linux package source: %s", cmd.PackageSource)
	}
	if cmd.CheckCommand != "zellij" {
		t.Errorf("unexpected Linux check command: %s", cmd.CheckCommand)
	}
	if cmd.Interactive {
		t.Error("Linux command should be non-interactive")
	}
}

// TestInstallZellij_MacOSCommandStructure validates the macOS command specifics.
func TestInstallZellij_MacOSCommandStructure(t *testing.T) {
	action := InstallZellij()
	cmd, ok := action.PlatformCommands[actions.PlatformMacOS]
	if !ok {
		t.Fatal("macOS platform missing")
	}
	if cmd.Command != "brew install zellij" {
		t.Errorf("unexpected macOS command: %s", cmd.Command)
	}
	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("unexpected macOS package source: %s", cmd.PackageSource)
	}
	if cmd.CheckCommand != "zellij" {
		t.Errorf("unexpected macOS check command: %s", cmd.CheckCommand)
	}
	if cmd.Interactive {
		t.Error("macOS command should be non-interactive")
	}
}

// TestInstallZellij_CheckCommandStructure ensures all platforms share the same check command.
func TestInstallZellij_CheckCommandStructure(t *testing.T) {
	action := InstallZellij()
	for p, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != "zellij" {
			t.Errorf("platform %s has unexpected check command %s", p, cmd.CheckCommand)
		}
	}
}

// TestInstallZellij_NoCommandExecution guarantees no side-effects during testing.
func TestInstallZellij_NoCommandExecution(t *testing.T) {
	action := InstallZellij()
	if action == nil {
		t.Fatal("action is nil")
	}
	if len(action.PlatformCommands) == 0 {
		t.Fatal("no platform commands configured")
	}
	for p, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("platform %s has empty command", p)
		}
	}
}

// TestInstallZellij_SafeForCI ensures the action can be instantiated in CI environments.
func TestInstallZellij_SafeForCI(t *testing.T) {
	action := InstallZellij()
	if action == nil {
		t.Fatal("action is nil")
	}
	for p, cmd := range action.PlatformCommands {
		if cmd.Command == "" {
			t.Errorf("platform %s missing command", p)
		}
		if cmd.CheckCommand == "" {
			t.Errorf("platform %s missing check command", p)
		}
	}
}

// TestInstallZellij_BrewConsistency validates brew-based platforms use the correct source.
func TestInstallZellij_BrewConsistency(t *testing.T) {
	action := InstallZellij()
	for _, p := range []actions.Platform{actions.PlatformLinux, actions.PlatformMacOS} {
		cmd, ok := action.PlatformCommands[p]
		if !ok {
			t.Fatalf("%s platform missing", p)
		}
		if cmd.PackageSource != actions.PackageSourceBrew {
			t.Errorf("%s should use Brew package source, got %s", p, cmd.PackageSource)
		}
	}
}

// TestInstallZellij_PackageSourceValidation checks that all package sources are among allowed values.
func TestInstallZellij_PackageSourceValidation(t *testing.T) {
	action := InstallZellij()
	valid := map[actions.PackageSource]bool{
		actions.PackageSourceOfficial: true,
		actions.PackageSourceBrew:     true,
		actions.PackageSourceAny:      true,
	}
	for p, cmd := range action.PlatformCommands {
		if !valid[cmd.PackageSource] {
			t.Errorf("platform %s has invalid package source %s", p, cmd.PackageSource)
		}
	}
}

// TestInstallZellij_InteractiveFlagConsistency verifies interactive flags per platform.
func TestInstallZellij_InteractiveFlagConsistency(t *testing.T) {
	action := InstallZellij()
	interactive := []actions.Platform{actions.PlatformArch}
	for _, p := range interactive {
		cmd, ok := action.PlatformCommands[p]
		if !ok {
			t.Fatalf("%s platform missing", p)
		}
		if !cmd.Interactive {
			t.Errorf("%s should be interactive", p)
		}
	}
	nonInteractive := []actions.Platform{actions.PlatformLinux, actions.PlatformMacOS}
	for _, p := range nonInteractive {
		cmd, ok := action.PlatformCommands[p]
		if !ok {
			t.Fatalf("%s platform missing", p)
		}
		if cmd.Interactive {
			t.Errorf("%s should be non-interactive", p)
		}
	}
}

// TestInstallZellij_CommandComplexity ensures commands are non-trivial and contain expected keywords.
func TestInstallZellij_CommandComplexity(t *testing.T) {
	action := InstallZellij()
	for p, cmd := range action.PlatformCommands {
		if len(cmd.Command) < 5 {
			t.Errorf("command for %s too short: %s", p, cmd.Command)
		}
		switch p {
		case actions.PlatformArch:
			if !utils.ContainsString(cmd.Command, "pacman") {
				t.Errorf("Arch command missing 'pacman': %s", cmd.Command)
			}
		case actions.PlatformLinux, actions.PlatformMacOS:
			if !utils.ContainsString(cmd.Command, "brew") {
				t.Errorf("%s command missing 'brew': %s", p, cmd.Command)
			}
		}
	}
}

// TestInstallZellij_ArchCommandSudo tests that Arch command properly uses sudo.
func TestInstallZellij_ArchCommandSudo(t *testing.T) {
	action := InstallZellij()
	cmd, ok := action.PlatformCommands[actions.PlatformArch]
	if !ok {
		t.Fatal("Arch platform missing")
	}
	if !utils.ContainsString(cmd.Command, "sudo") {
		t.Errorf("Arch command should use sudo: %s", cmd.Command)
	}
}

// TestInstallZellij_PacmanPackageManager ensures Arch uses pacman.
func TestInstallZellij_PacmanPackageManager(t *testing.T) {
	action := InstallZellij()
	cmd, ok := action.PlatformCommands[actions.PlatformArch]
	if !ok {
		t.Fatal("Arch platform missing")
	}
	if !utils.ContainsString(cmd.Command, "pacman") {
		t.Errorf("Arch should use pacman package manager: %s", cmd.Command)
	}
}

// TestInstallZellij_PlatformCount verifies the correct number of platforms are supported.
func TestInstallZellij_PlatformCount(t *testing.T) {
	action := InstallZellij()
	expectedCount := 3 // Arch, Linux, macOS
	if len(action.PlatformCommands) != expectedCount {
		t.Errorf("expected %d platforms, got %d", expectedCount, len(action.PlatformCommands))
	}
}

// TestInstallZellij_AllCommandsHaveCheckCommands ensures every platform command has a check command.
func TestInstallZellij_AllCommandsHaveCheckCommands(t *testing.T) {
	action := InstallZellij()
	for p, cmd := range action.PlatformCommands {
		if cmd.CheckCommand == "" {
			t.Errorf("platform %s has no check command", p)
		}
	}
}

// TestInstallZellij_AllCommandsHavePackageSources ensures every platform command has a package source.
func TestInstallZellij_AllCommandsHavePackageSources(t *testing.T) {
	action := InstallZellij()
	for p, cmd := range action.PlatformCommands {
		if cmd.PackageSource == "" {
			t.Errorf("platform %s has no package source", p)
		}
	}
}

// TestInstallZellij_ZellijKeywordPresence ensures all commands reference zellij.
func TestInstallZellij_ZellijKeywordPresence(t *testing.T) {
	action := InstallZellij()
	for p, cmd := range action.PlatformCommands {
		if !utils.ContainsString(cmd.Command, "zellij") {
			t.Errorf("platform %s command missing 'zellij': %s", p, cmd.Command)
		}
	}
}

// TestInstallZellij_ActionTypeIsCommand verifies the action type is correctly set.
func TestInstallZellij_ActionTypeIsCommand(t *testing.T) {
	action := InstallZellij()
	if action.Type != actions.ActionTypeCommand {
		t.Errorf("expected action type to be %s, got %s", actions.ActionTypeCommand, action.Type)
	}
}

// TestInstallZellij_IDFormat verifies the action ID follows the expected format.
func TestInstallZellij_IDFormat(t *testing.T) {
	action := InstallZellij()
	if action.ID != "install_zellij" {
		t.Errorf("expected ID format 'install_zellij', got '%s'", action.ID)
	}
}

// TestInstallZellij_NameFormat verifies the action name is properly formatted.
func TestInstallZellij_NameFormat(t *testing.T) {
	action := InstallZellij()
	if action.Name != "Install Zellij" {
		t.Errorf("expected name 'Install Zellij', got '%s'", action.Name)
	}
}

// TestInstallZellij_DescriptionFormat verifies the action description is properly formatted.
func TestInstallZellij_DescriptionFormat(t *testing.T) {
	action := InstallZellij()
	expectedDesc := "Install Zellij, a terminal multiplexer, on your system"
	if action.Description != expectedDesc {
		t.Errorf("expected description '%s', got '%s'", expectedDesc, action.Description)
	}
}
