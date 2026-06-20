# Codory Action Test Template

Create one test file beside the action.

Path:

```text
internal/actions/<category>/install_<app_snake>_test.go
```

Before using this template, inspect nearby `install_*_test.go` files and adapt expected platforms, commands, sources, and unsupported-platform checks to current project constants and patterns.

Template. Example test rows are illustrative; replace with actual supported platform commands and unsupported platforms:

```go
package <category_package>

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestInstall<AppPascal>(t *testing.T) {
	action := Install<AppPascal>()

	if action == nil {
		t.Fatal("Install<AppPascal>() returned nil")
	}
	if action.ID != "install_<app_snake>" {
		t.Errorf("Expected action ID %q, got %q", "install_<app_snake>", action.ID)
	}
	if action.Name != "Install <AppName>" {
		t.Errorf("Expected action name %q, got %q", "Install <AppName>", action.Name)
	}
	if action.Description != "Install <short description>" {
		t.Errorf("Expected description %q, got %q", "Install <short description>", action.Description)
	}
	if action.Type != actions.ActionTypeCommand {
		t.Errorf("Expected type %q, got %q", actions.ActionTypeCommand, action.Type)
	}
	if action.PlatformCommands == nil {
		t.Fatal("Expected PlatformCommands")
	}
}

func TestInstall<AppPascal>_PlatformCommands(t *testing.T) {
	action := Install<AppPascal>()

	tests := []struct {
		name        string
		platform    actions.Platform
		command     string
		source      actions.PackageSource
		check       string
		interactive bool
	}{
		{
			name:        "example platform",
			platform:    actions.PlatformMacOS,
			command:     "brew install <formula>",
			source:      actions.PackageSourceBrew,
			check:       "<binary>",
			interactive: false,
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

func TestInstall<AppPascal>_UnsupportedPlatforms(t *testing.T) {
	action := Install<AppPascal>()

	unsupported := []actions.Platform{
		actions.PlatformWindows,
	}

	for _, platform := range unsupported {
		t.Run(string(platform), func(t *testing.T) {
			if _, ok := action.PlatformCommands[platform]; ok {
				t.Errorf("Expected no command for unsupported platform %s", platform)
			}
		})
	}
}

func TestInstall<AppPascal>_ActionConsistency(t *testing.T) {
	action1 := Install<AppPascal>()
	action2 := Install<AppPascal>()

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
```

Keep tests as data checks. Never call executor or shell commands.
