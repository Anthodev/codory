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
					PlatformLinux: {Command: "apt update"},
					PlatformMacOS: {Command: "brew update"},
				},
			},
			platform: PlatformLinux,
			want:     PlatformCommand{Command: "apt update"},
			want1:    true,
		},
		{
			name: "command does not exist for platform",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "apt update"},
				},
			},
			platform: PlatformMacOS,
			want:     PlatformCommand{},
			want1:    false,
		},
		{
			name: "empty platform commands",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{},
			},
			platform: PlatformLinux,
			want:     PlatformCommand{},
			want1:    false,
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
					PlatformLinux: {Command: "apt update"},
					PlatformMacOS: {Command: "brew update"},
				},
			},
			platform: PlatformLinux,
			want:     true,
		},
		{
			name: "debian falls back to linux",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "apt update"},
				},
			},
			platform: PlatformDebian,
			want:     true,
		},
		{
			name: "arch falls back to linux",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "pacman -Syu"},
				},
			},
			platform: PlatformArch,
			want:     true,
		},
		{
			name: "fallback to PlatformAny",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformAny: {Command: "echo hello"},
				},
			},
			platform: PlatformWindows,
			want:     true,
		},
		{
			name: "no command for platform",
			action: Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformMacOS: {Command: "brew update"},
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
