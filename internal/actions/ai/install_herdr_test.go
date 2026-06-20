package ai

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestInstallHerdr(t *testing.T) {
	action := InstallHerdr()

	if action == nil {
		t.Fatal("InstallHerdr() returned nil")
	}
	if action.ID != "install_herdr" {
		t.Errorf("Expected action ID %q, got %q", "install_herdr", action.ID)
	}
	if action.Name != "Install Herdr" {
		t.Errorf("Expected action name %q, got %q", "Install Herdr", action.Name)
	}
	if action.Description != "Install Herdr agent multiplexer" {
		t.Errorf("Expected description %q, got %q", "Install Herdr agent multiplexer", action.Description)
	}
	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected type %q, got %q", actions.ActionTypeCommand, action.Type)
	}
	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands")
	}
}

func TestInstallHerdr_PlatformCommands(t *testing.T) {
	action := InstallHerdr()

	tests := []struct {
		name        string
		platform    actions.Platform
		command     string
		source      actions.PackageSource
		check       string
		interactive bool
	}{
		{
			name:        "Arch",
			platform:    actions.PlatformArch,
			command:     "yay -S herdr-bin",
			source:      actions.PackageSourceAUR,
			check:       "herdr",
			interactive: true,
		},
		{
			name:        "Linux",
			platform:    actions.PlatformLinux,
			command:     "brew install herdr",
			source:      actions.PackageSourceBrew,
			check:       "herdr",
			interactive: false,
		},
		{
			name:        "macOS",
			platform:    actions.PlatformMacOS,
			command:     "brew install herdr",
			source:      actions.PackageSourceBrew,
			check:       "herdr",
			interactive: false,
		},
		{
			name:        "Windows preview",
			platform:    actions.PlatformWindows,
			command:     "powershell -ExecutionPolicy Bypass -c \"irm https://herdr.dev/install.ps1 | iex\"",
			source:      actions.PackageSourceAny,
			check:       "herdr",
			interactive: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, ok := action.PlatformCommands[tt.platform]
			if !ok {
				t.Fatalf("Expected command for %s", tt.platform)
			}
			if cmd.Command != tt.command {
				t.Errorf("Expected command %q, got %q", tt.command, cmd.Command)
			}
			if cmd.PackageSource != tt.source {
				t.Errorf("Expected source %q, got %q", tt.source, cmd.PackageSource)
			}
			if cmd.CheckCommand != tt.check {
				t.Errorf("Expected check %q, got %q", tt.check, cmd.CheckCommand)
			}
			if cmd.Interactive != tt.interactive {
				t.Errorf("Expected interactive %v, got %v", tt.interactive, cmd.Interactive)
			}
		})
	}
}

func TestInstallHerdr_UnsupportedPlatforms(t *testing.T) {
	action := InstallHerdr()

	unsupported := []actions.Platform{
		actions.PlatformDebian,
	}

	for _, platform := range unsupported {
		t.Run(string(platform), func(t *testing.T) {
			if _, ok := action.PlatformCommands[platform]; ok {
				t.Errorf("Expected no command for unsupported platform %s", platform)
			}
		})
	}
}

func TestInstallHerdr_ActionConsistency(t *testing.T) {
	action1 := InstallHerdr()
	action2 := InstallHerdr()

	if action1.ID != action2.ID || action1.Name != action2.Name || action1.Description != action2.Description || action1.Type != action2.Type {
		t.Fatal("Expected deterministic action metadata")
	}
	if len(action1.PlatformCommands) != len(action2.PlatformCommands) {
		t.Fatalf("Expected deterministic platform command count")
	}
	for platform, cmd1 := range action1.PlatformCommands {
		cmd2, ok := action2.PlatformCommands[platform]
		if !ok {
			t.Fatalf("Platform %s missing in second action", platform)
		}
		if cmd1 != cmd2 {
			t.Errorf("Expected deterministic command for %s", platform)
		}
	}
}
