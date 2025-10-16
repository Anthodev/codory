package terminal

import (
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
)

// TestInstallTmux verifies the basic fields of the InstallTmux action.
func TestInstallTmux(t *testing.T) {
	action := InstallTmux()

	if action == nil {
		t.Fatal("InstallTmux() returned nil")
	}
	if action.ID != "install_tmux" {
		t.Errorf("expected ID 'install_tmux', got '%s'", action.ID)
	}
	if action.Name != "Install Tmux" {
		t.Errorf("expected Name 'Install Tmux', got '%s'", action.Name)
	}
	if action.Description != "Install Tmux, the terminal multiplexer, on your system" {
		t.Errorf("unexpected description: %s", action.Description)
	}
	if action.Type != actions.ActionTypeCommand {
		t.Errorf("expected Type %s, got %s", actions.ActionTypeCommand, action.Type)
	}
	if action.PlatformCommands == nil {
		t.Fatal("PlatformCommands map is nil")
	}
}

// TestInstallTmux_PlatformCommands validates command configuration per platform.
func TestInstallTmux_PlatformCommands(t *testing.T) {
	action := InstallTmux()
	tests := []struct {
		platform            actions.Platform
		expectedCommand     string
		expectedSource      actions.PackageSource
		expectedCheck       string
		expectedInteractive bool
	}{
		{
			platform:            actions.PlatformArch,
			expectedCommand:     "sudo pacman -S tmux",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "tmux",
			expectedInteractive: true,
		},
		{
			platform:            actions.PlatformDebian,
			expectedCommand:     "sudo apt-get install tmux",
			expectedSource:      actions.PackageSourceOfficial,
			expectedCheck:       "tmux",
			expectedInteractive: true,
		},
		{
			platform:            actions.PlatformLinux,
			expectedCommand:     "brew install tmux",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "tmux",
			expectedInteractive: false,
		},
		{
			platform:            actions.PlatformMacOS,
			expectedCommand:     "brew install tmux",
			expectedSource:      actions.PackageSourceBrew,
			expectedCheck:       "tmux",
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

// TestInstallTmux_UnsupportedPlatforms ensures unsupported platforms have no entry.
func TestInstallTmux_UnsupportedPlatforms(t *testing.T) {
	action := InstallTmux()
	unsupported := []actions.Platform{actions.PlatformWindows}
	for _, p := range unsupported {
		if _, ok := action.PlatformCommands[p]; ok {
			t.Errorf("unexpected command for unsupported platform %s", p)
		}
	}
}

// TestInstallTmux_ActionConsistency checks that multiple calls return identical actions.
func TestInstallTmux_ActionConsistency(t *testing.T) {
	a1 := InstallTmux()
	a2 := InstallTmux()

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

// TestInstallTmux_ExpectedPlatforms verifies all expected platforms are present.
func TestInstallTmux_ExpectedPlatforms(t *testing.T) {
	action := InstallTmux()
	expected := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}
	for _, p := range expected {
		if _, ok := action.PlatformCommands[p]; !ok {
			t.Errorf("expected platform %s to be present", p)
		}
	}
}

// TestInstallTmux_ArchCommandStructure validates the Arch command specifics.
func TestInstallTmux_ArchCommandStructure(t *testing.T) {
	action := InstallTmux()
	cmd, ok := action.PlatformCommands[actions.PlatformArch]
	if !ok {
		t.Fatal("Arch platform missing")
	}
	if cmd.Command != "sudo pacman -S tmux" {
		t.Errorf("unexpected Arch command: %s", cmd.Command)
	}
	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("unexpected Arch package source: %s", cmd.PackageSource)
	}
	if cmd.CheckCommand != "tmux" {
		t.Errorf("unexpected Arch check command: %s", cmd.CheckCommand)
	}
	if !cmd.Interactive {
		t.Error("Arch command should be interactive")
	}
}

// TestInstallTmux_DebianCommandStructure validates the Debian command specifics.
func TestInstallTmux_DebianCommandStructure(t *testing.T) {
	action := InstallTmux()
	cmd, ok := action.PlatformCommands[actions.PlatformDebian]
	if !ok {
		t.Fatal("Debian platform missing")
	}
	if cmd.Command != "sudo apt-get install tmux" {
		t.Errorf("unexpected Debian command: %s", cmd.Command)
	}
	if cmd.PackageSource != actions.PackageSourceOfficial {
		t.Errorf("unexpected Debian package source: %s", cmd.PackageSource)
	}
	if cmd.CheckCommand != "tmux" {
		t.Errorf("unexpected Debian check command: %s", cmd.CheckCommand)
	}
	if !cmd.Interactive {
		t.Error("Debian command should be interactive")
	}
}

// TestInstallTmux_LinuxCommandStructure validates the generic Linux command specifics.
func TestInstallTmux_LinuxCommandStructure(t *testing.T) {
	action := InstallTmux()
	cmd, ok := action.PlatformCommands[actions.PlatformLinux]
	if !ok {
		t.Fatal("Linux platform missing")
	}
	if cmd.Command != "brew install tmux" {
		t.Errorf("unexpected Linux command: %s", cmd.Command)
	}
	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("unexpected Linux package source: %s", cmd.PackageSource)
	}
	if cmd.CheckCommand != "tmux" {
		t.Errorf("unexpected Linux check command: %s", cmd.CheckCommand)
	}
	if cmd.Interactive {
		t.Error("Linux command should be non‑interactive")
	}
}

// TestInstallTmux_MacOSCommandStructure validates the macOS command specifics.
func TestInstallTmux_MacOSCommandStructure(t *testing.T) {
	action := InstallTmux()
	cmd, ok := action.PlatformCommands[actions.PlatformMacOS]
	if !ok {
		t.Fatal("macOS platform missing")
	}
	if cmd.Command != "brew install tmux" {
		t.Errorf("unexpected macOS command: %s", cmd.Command)
	}
	if cmd.PackageSource != actions.PackageSourceBrew {
		t.Errorf("unexpected macOS package source: %s", cmd.PackageSource)
	}
	if cmd.CheckCommand != "tmux" {
		t.Errorf("unexpected macOS check command: %s", cmd.CheckCommand)
	}
	if cmd.Interactive {
		t.Error("macOS command should be non‑interactive")
	}
}

// TestInstallTmux_CheckCommandStructure ensures all platforms share the same check command.
func TestInstallTmux_CheckCommandStructure(t *testing.T) {
	action := InstallTmux()
	for p, cmd := range action.PlatformCommands {
		if cmd.CheckCommand != "tmux" {
			t.Errorf("platform %s has unexpected check command %s", p, cmd.CheckCommand)
		}
	}
}

// TestInstallTmux_NoCommandExecution guarantees no side‑effects during testing.
func TestInstallTmux_NoCommandExecution(t *testing.T) {
	action := InstallTmux()
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

// TestInstallTmux_SafeForCI ensures the action can be instantiated in CI environments.
func TestInstallTmux_SafeForCI(t *testing.T) {
	action := InstallTmux()
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

// TestInstallTmux_BrewConsistency validates brew‑based platforms use the correct source.
func TestInstallTmux_BrewConsistency(t *testing.T) {
	action := InstallTmux()
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

// TestInstallTmux_PackageSourceValidation checks that all package sources are among allowed values.
func TestInstallTmux_PackageSourceValidation(t *testing.T) {
	action := InstallTmux()
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

// TestInstallTmux_InteractiveFlagConsistency verifies interactive flags per platform.
func TestInstallTmux_InteractiveFlagConsistency(t *testing.T) {
	action := InstallTmux()
	interactive := []actions.Platform{actions.PlatformArch, actions.PlatformDebian}
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
			t.Errorf("%s should be non‑interactive", p)
		}
	}
}

// TestInstallTmux_CommandComplexity ensures commands are non‑trivial and contain expected keywords.
func TestInstallTmux_CommandComplexity(t *testing.T) {
	action := InstallTmux()
	for p, cmd := range action.PlatformCommands {
		if len(cmd.Command) < 5 {
			t.Errorf("command for %s too short: %s", p, cmd.Command)
		}
		switch p {
		case actions.PlatformArch:
			if !utils.ContainsString(cmd.Command, "pacman") {
				t.Errorf("Arch command missing 'pacman': %s", cmd.Command)
			}
		case actions.PlatformDebian:
			if !utils.ContainsString(cmd.Command, "apt-get") {
				t.Errorf("Debian command missing 'apt-get': %s", cmd.Command)
			}
		case actions.PlatformLinux, actions.PlatformMacOS:
			if !utils.ContainsString(cmd.Command, "brew") {
				t.Errorf("%s command missing 'brew': %s", p, cmd.Command)
			}
		}
	}
}
