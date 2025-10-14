package actions

import (
	"testing"
)

func TestCategory_IsLeaf(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		want     bool
	}{
		{
			name: "category with no subcategories",
			category: Category{
				SubCategories: []*Category{},
			},
			want: true,
		},
		{
			name: "category with subcategories",
			category: Category{
				SubCategories: []*Category{
					{ID: "sub1"},
					{ID: "sub2"},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.IsLeaf(); got != tt.want {
				t.Errorf("Category.IsLeaf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategory_HasActions(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		want     bool
	}{
		{
			name: "category with no actions",
			category: Category{
				Actions: []*Action{},
			},
			want: false,
		},
		{
			name: "category with actions",
			category: Category{
				Actions: []*Action{
					{ID: "action1"},
					{ID: "action2"},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.HasActions(); got != tt.want {
				t.Errorf("Category.HasActions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategory_IsVisibleOnPlatform(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		platform Platform
		want     bool
	}{
		{
			name: "category visible on platform",
			category: Category{
				HiddenOnPlatforms: []Platform{PlatformWindows, PlatformMacOS},
			},
			platform: PlatformLinux,
			want:     true,
		},
		{
			name: "category hidden on platform",
			category: Category{
				HiddenOnPlatforms: []Platform{PlatformLinux, PlatformWindows},
			},
			platform: PlatformLinux,
			want:     false,
		},
		{
			name: "category with no hidden platforms",
			category: Category{
				HiddenOnPlatforms: []Platform{},
			},
			platform: PlatformLinux,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.IsVisibleOnPlatform(tt.platform); got != tt.want {
				t.Errorf("Category.IsVisibleOnPlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAction_IsVisibleOnPlatform(t *testing.T) {
	tests := []struct {
		name     string
		action   Action
		platform Platform
		want     bool
	}{
		{
			name: "action visible on platform",
			action: Action{
				HiddenOnPlatforms: []Platform{PlatformWindows, PlatformMacOS},
			},
			platform: PlatformLinux,
			want:     true,
		},
		{
			name: "action hidden on platform",
			action: Action{
				HiddenOnPlatforms: []Platform{PlatformLinux, PlatformWindows},
			},
			platform: PlatformLinux,
			want:     false,
		},
		{
			name: "action with no hidden platforms",
			action: Action{
				HiddenOnPlatforms: []Platform{},
			},
			platform: PlatformLinux,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.IsVisibleOnPlatform(tt.platform); got != tt.want {
				t.Errorf("Action.IsVisibleOnPlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategory_IsHiddenOnPlatform(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		platform Platform
		want     bool
	}{
		{
			name: "category hidden on platform",
			category: Category{
				HiddenOnPlatforms: []Platform{PlatformLinux, PlatformWindows},
			},
			platform: PlatformLinux,
			want:     true,
		},
		{
			name: "category not hidden on platform",
			category: Category{
				HiddenOnPlatforms: []Platform{PlatformWindows, PlatformMacOS},
			},
			platform: PlatformLinux,
			want:     false,
		},
		{
			name: "category with no hidden platforms",
			category: Category{
				HiddenOnPlatforms: []Platform{},
			},
			platform: PlatformLinux,
			want:     false,
		},
		{
			name: "category with nil hidden platforms",
			category: Category{
				HiddenOnPlatforms: nil,
			},
			platform: PlatformLinux,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.IsHiddenOnPlatform(tt.platform); got != tt.want {
				t.Errorf("Category.IsHiddenOnPlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAction_IsHiddenOnPlatform(t *testing.T) {
	tests := []struct {
		name     string
		action   Action
		platform Platform
		want     bool
	}{
		{
			name: "action hidden on platform",
			action: Action{
				HiddenOnPlatforms: []Platform{PlatformLinux, PlatformWindows},
			},
			platform: PlatformLinux,
			want:     true,
		},
		{
			name: "action not hidden on platform",
			action: Action{
				HiddenOnPlatforms: []Platform{PlatformWindows, PlatformMacOS},
			},
			platform: PlatformLinux,
			want:     false,
		},
		{
			name: "action with no hidden platforms",
			action: Action{
				HiddenOnPlatforms: []Platform{},
			},
			platform: PlatformLinux,
			want:     false,
		},
		{
			name: "action with nil hidden platforms",
			action: Action{
				HiddenOnPlatforms: nil,
			},
			platform: PlatformLinux,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.IsHiddenOnPlatform(tt.platform); got != tt.want {
				t.Errorf("Action.IsHiddenOnPlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAction_GetPlatformCommand(t *testing.T) {
	tests := []struct {
		name     string
		action   Action
		platform Platform
		want     PlatformCommand
		want1    bool
	}{
		{
			name: "command exists for platform",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "apt update", CheckCommand: "which apt"},
					PlatformMacOS: {Command: "brew update", CheckCommand: "which brew"},
				},
			},
			platform: PlatformLinux,
			want:     PlatformCommand{Command: "apt update", CheckCommand: "which apt"},
			want1:    true,
		},
		{
			name: "command does not exist for platform",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "apt update", CheckCommand: "which apt"},
				},
			},
			platform: PlatformMacOS,
			want:     PlatformCommand{Command: "", PackageSource: "", CheckCommand: ""},
			want1:    false,
		},
		{
			name: "empty platform commands",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{},
			},
			platform: PlatformLinux,
			want:     PlatformCommand{Command: "", PackageSource: "", CheckCommand: ""},
			want1:    false,
		},
		{
			name: "interactive command",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "sudo pacman -S helix", CheckCommand: "hx", Interactive: true},
				},
			},
			platform: PlatformLinux,
			want:     PlatformCommand{Command: "sudo pacman -S helix", CheckCommand: "hx", Interactive: true},
			want1:    true,
		},
		{
			name: "non-interactive command (default)",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "echo test", CheckCommand: "which echo", Interactive: false},
				},
			},
			platform: PlatformLinux,
			want:     PlatformCommand{Command: "echo test", CheckCommand: "which echo", Interactive: false},
			want1:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.action.GetPlatformCommand(tt.platform)
			if got != tt.want {
				t.Errorf("Action.GetPlatformCommand() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Action.GetPlatformCommand() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAction_HasCommandForPlatform(t *testing.T) {
	tests := []struct {
		name     string
		action   Action
		platform Platform
		want     bool
	}{
		{
			name: "direct platform match",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "apt update", CheckCommand: "which apt"},
					PlatformMacOS: {Command: "brew update", CheckCommand: "which brew"},
				},
			},
			platform: PlatformLinux,
			want:     true,
		},
		{
			name: "debian falls back to linux",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "apt update", CheckCommand: "which apt"},
				},
			},
			platform: PlatformDebian,
			want:     true,
		},
		{
			name: "arch falls back to linux",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "pacman -Syu", CheckCommand: "which pacman"},
				},
			},
			platform: PlatformArch,
			want:     true,
		},
		{
			name: "fallback to PlatformAny",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformAny: {Command: "echo hello", CheckCommand: "which echo"},
				},
			},
			platform: PlatformWindows,
			want:     true,
		},
		{
			name: "no command for platform",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformMacOS: {Command: "brew update", CheckCommand: "which brew"},
				},
			},
			platform: PlatformWindows,
			want:     false,
		},
		{
			name: "empty platform commands",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{},
			},
			platform: PlatformLinux,
			want:     false,
		},
		{
			name: "nil platform commands",
			action: Action{
				PlatformCommands: nil,
			},
			platform: PlatformLinux,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.HasCommandForPlatform(tt.platform); got != tt.want {
				t.Errorf("Action.HasCommandForPlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionArgument(t *testing.T) {
	tests := []struct {
		name     string
		argument ActionArgument
		want     struct {
			name        string
			description string
			required    bool
		}
	}{
		{
			name: "required argument",
			argument: ActionArgument{
				Name:        "package",
				Description: "Package name to install",
				Required:    true,
			},
			want: struct {
				name        string
				description string
				required    bool
			}{
				name:        "package",
				description: "Package name to install",
				required:    true,
			},
		},
		{
			name: "optional argument",
			argument: ActionArgument{
				Name:        "verbose",
				Description: "Enable verbose output",
				Required:    false,
			},
			want: struct {
				name        string
				description string
				required    bool
			}{
				name:        "verbose",
				description: "Enable verbose output",
				required:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.argument.Name; got != tt.want.name {
				t.Errorf("ActionArgument.Name = %v, want %v", got, tt.want.name)
			}
			if got := tt.argument.Description; got != tt.want.description {
				t.Errorf("ActionArgument.Description = %v, want %v", got, tt.want.description)
			}
			if got := tt.argument.Required; got != tt.want.required {
				t.Errorf("ActionArgument.Required = %v, want %v", got, tt.want.required)
			}
		})
	}
}

func TestCategory_EdgeCases(t *testing.T) {
	t.Run("nil subcategories", func(t *testing.T) {
		category := Category{
			SubCategories: nil,
		}
		if got := category.IsLeaf(); got != true {
			t.Errorf("Category.IsLeaf() with nil SubCategories = %v, want %v", got, true)
		}
	})

	t.Run("nil actions", func(t *testing.T) {
		category := Category{
			Actions: nil,
		}
		if got := category.HasActions(); got != false {
			t.Errorf("Category.HasActions() with nil Actions = %v, want %v", got, false)
		}
	})
}

func TestAction_HasCheckCommand(t *testing.T) {
	tests := []struct {
		name     string
		action   Action
		platform Platform
		want     bool
	}{
		{
			name: "action with check command",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command:      "apt update",
						CheckCommand: "which apt",
					},
				},
			},
			platform: PlatformLinux,
			want:     true,
		},
		{
			name: "action without check command",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command: "apt update",
					},
				},
			},
			platform: PlatformLinux,
			want:     false,
		},
		{
			name: "no platform command",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformMacOS: {
						Command:      "brew update",
						CheckCommand: "which brew",
					},
				},
			},
			platform: PlatformLinux,
			want:     false,
		},
		{
			name: "nil platform commands",
			action: Action{
				PlatformCommands: nil,
			},
			platform: PlatformLinux,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.HasCheckCommand(tt.platform); got != tt.want {
				t.Errorf("Action.HasCheckCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAction_GetCheckCommand(t *testing.T) {
	tests := []struct {
		name     string
		action   Action
		platform Platform
		want     string
		want1    bool
	}{
		{
			name: "action with check command",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command:      "apt update",
						CheckCommand: "which apt",
					},
				},
			},
			platform: PlatformLinux,
			want:     "which apt",
			want1:    true,
		},
		{
			name: "action without check command",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command: "apt update",
					},
				},
			},
			platform: PlatformLinux,
			want:     "",
			want1:    false,
		},
		{
			name: "no platform command",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformMacOS: {
						Command:      "brew update",
						CheckCommand: "which brew",
					},
				},
			},
			platform: PlatformLinux,
			want:     "",
			want1:    false,
		},
		{
			name: "nil platform commands",
			action: Action{
				PlatformCommands: nil,
			},
			platform: PlatformLinux,
			want:     "",
			want1:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.action.GetCheckCommand(tt.platform)
			if got != tt.want {
				t.Errorf("Action.GetCheckCommand() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Action.GetCheckCommand() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAction_EdgeCases(t *testing.T) {
	t.Run("nil hidden platforms", func(t *testing.T) {
		action := Action{
			HiddenOnPlatforms: nil,
		}
		if got := action.IsVisibleOnPlatform(PlatformLinux); got != true {
			t.Errorf("Action.IsVisibleOnPlatform() with nil HiddenOnPlatforms = %v, want %v", got, true)
		}
		if got := action.IsHiddenOnPlatform(PlatformLinux); got != false {
			t.Errorf("Action.IsHiddenOnPlatform() with nil HiddenOnPlatforms = %v, want %v", got, false)
		}
	})

	t.Run("nil platform commands", func(t *testing.T) {
		action := Action{
			PlatformCommands: nil,
		}
		got, ok := action.GetPlatformCommand(PlatformLinux)
		if ok != false {
			t.Errorf("Action.GetPlatformCommand() with nil PlatformCommands ok = %v, want %v", ok, false)
		}
		if got != (PlatformCommand{Command: "", PackageSource: "", CheckCommand: ""}) {
			t.Errorf("Action.GetPlatformCommand() with nil PlatformCommands got = %v, want %v", got, PlatformCommand{Command: "", PackageSource: "", CheckCommand: ""})
		}
	})
}

func TestPlatformCommand_Interactive(t *testing.T) {
	tests := []struct {
		name    string
		command PlatformCommand
		want    bool
	}{
		{
			name: "interactive command",
			command: PlatformCommand{
				Command:       "sudo pacman -S helix",
				CheckCommand:  "hx",
				Interactive:   true,
				PackageSource: PackageSourceOfficial,
			},
			want: true,
		},
		{
			name: "non-interactive command",
			command: PlatformCommand{
				Command:       "echo test",
				CheckCommand:  "which echo",
				Interactive:   false,
				PackageSource: PackageSourceOfficial,
			},
			want: false,
		},
		{
			name: "command with no interactive flag set (defaults to false)",
			command: PlatformCommand{
				Command:      "curl -fsSL https://example.com",
				CheckCommand: "which curl",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.command.Interactive; got != tt.want {
				t.Errorf("PlatformCommand.Interactive = %v, want %v", got, tt.want)
			}
		})
	}
}
