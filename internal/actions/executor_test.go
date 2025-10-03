package actions

import (
	"context"
	"testing"

	"anthodev/codory/internal/platform"
)

// Mock platform info for testing
func createMockPlatformInfo(os platform.Platform) platform.Info {
	return platform.Info{
		OS:              os,
		PackageManagers: []platform.PackageManager{},
	}
}

// Helper function to create test platform info
func createTestPlatformInfo(os platform.Platform) platform.Info {
	return platform.Info{
		OS:              os,
		PackageManagers: []platform.PackageManager{},
	}
}

func TestNewExecutor(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("NewExecutor() returned nil")
	}
	if executor.platformInfo.OS == "" {
		t.Error("Expected platform info to be set")
	}
}

func TestExecutor_Execute(t *testing.T) {
	tests := []struct {
		name        string
		action      *Action
		wantErr     bool
		errContains string
	}{
		{
			name: "function action with handler",
			action: &Action{
				ID:   "test-function",
				Type: ActionTypeFunction,
				Handler: func(ctx context.Context) (string, error) {
					return "function result", nil
				},
			},
			wantErr: false,
		},
		{
			name: "function action without handler",
			action: &Action{
				ID:   "test-function-no-handler",
				Type: ActionTypeFunction,
			},
			wantErr:     true,
			errContains: "no handler defined",
		},
		{
			name: "unknown action type",
			action: &Action{
				ID:   "test-unknown",
				Type: ActionType("unknown"),
			},
			wantErr:     true,
			errContains: "unknown action type",
		},
	}

	executor := NewExecutor()
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := executor.Execute(ctx, tt.action)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.action.Type == ActionTypeFunction && result != "function result" {
					t.Errorf("Expected 'function result', got '%s'", result)
				}
			}
		})
	}
}

func TestExecutor_executeCommand(t *testing.T) {
	tests := []struct {
		name        string
		action      *Action
		mockOS      Platform
		wantErr     bool
		errContains string
	}{
		{
			name: "command for specific platform",
			action: &Action{
				ID:   "test-command",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "echo test"},
				},
			},
			mockOS:  PlatformLinux,
			wantErr: false,
		},
		{
			name: "command with PlatformAny fallback",
			action: &Action{
				ID:   "test-command-any",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformAny: {Command: "echo test"},
				},
			},
			mockOS:  PlatformLinux,
			wantErr: false,
		},
		{
			name: "no command for platform",
			action: &Action{
				ID:   "test-no-command",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformMacOS: {Command: "echo test"},
				},
			},
			mockOS:      PlatformLinux,
			wantErr:     true,
			errContains: "no command defined for platform",
		},
		{
			name: "empty command",
			action: &Action{
				ID:   "test-empty-command",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: ""},
				},
			},
			mockOS:      PlatformLinux,
			wantErr:     true,
			errContains: "empty command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &Executor{
				platformInfo: createTestPlatformInfo(platform.Platform(tt.mockOS)),
			}

			result, err := executor.executeCommand(context.Background(), tt.action)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == "" {
					t.Error("Expected non-empty result")
				}
			}
		})
	}
}

func TestExecutor_GetPlatformInfo(t *testing.T) {
	executor := NewExecutor()
	info := executor.GetPlatformInfo()

	if info.OS == "" {
		t.Error("Expected platform info to have OS set")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

func TestExecutor_checkPackageManagerDependency(t *testing.T) {
	tests := []struct {
		name        string
		executor    *Executor
		source      PackageSource
		wantErr     bool
		errContains string
		skipFunc    func() bool
	}{
		{
			name: "AUR on non-Arch platform",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.Debian),
			},
			source:      PackageSourceAUR,
			wantErr:     true,
			errContains: "AUR packages are only available on Arch Linux",
		},
		{
			name: "AUR on Arch with yay installed",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.Arch),
			},
			source:  PackageSourceAUR,
			wantErr: false,
			skipFunc: func() bool {
				// Skip this test case if yay is not actually installed on the system
				return !platform.IsYayInstalled()
			},
		},
		{
			name: "official package source",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.Debian),
			},
			source:  PackageSourceOfficial,
			wantErr: false,
		},
		{
			name: "brew package source",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.MacOS),
			},
			source:  PackageSourceBrew,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipFunc != nil && tt.skipFunc() {
				t.Skip("Skipping test case based on skip function")
			}
			err := tt.executor.checkPackageManagerDependency(context.Background(), tt.source)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestExecutor_NeedsPackageManagerInstallation(t *testing.T) {
	tests := []struct {
		name       string
		executor   *Executor
		action     *Action
		wantNeeds  bool
		wantSource PackageSource
		skipFunc   func() bool
	}{
		{
			name: "AUR action on Arch without yay",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.Arch),
			},
			action: &Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformArch: {
						Command:       "yay -S package",
						PackageSource: PackageSourceAUR,
					},
				},
			},
			wantNeeds:  true,
			wantSource: PackageSourceAUR,
			skipFunc: func() bool {
				// Skip this test case if yay is actually installed on the system
				return platform.IsYayInstalled()
			},
		},
		{
			name: "AUR action on non-Arch platform",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.Debian),
			},
			action: &Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformAny: {
						Command:       "yay -S package",
						PackageSource: PackageSourceAUR,
					},
				},
			},
			wantNeeds:  false,
			wantSource: "",
		},
		{
			name: "official package action",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.Debian),
			},
			action: &Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformDebian: {
						Command:       "apt install package",
						PackageSource: PackageSourceOfficial,
					},
				},
			},
			wantNeeds:  false,
			wantSource: "",
		},
		{
			name: "no platform command",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.Debian),
			},
			action: &Action{
				PlatformCommands: map[Platform]PlatformCommand{},
			},
			wantNeeds:  false,
			wantSource: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipFunc != nil && tt.skipFunc() {
				t.Skip("Skipping test case based on skip function")
			}
			needs, source := tt.executor.NeedsPackageManagerInstallation(tt.action)

			if needs != tt.wantNeeds {
				t.Errorf("NeedsPackageManagerInstallation() needs = %v, want %v", needs, tt.wantNeeds)
			}
			if source != tt.wantSource {
				t.Errorf("NeedsPackageManagerInstallation() source = %v, want %v", source, tt.wantSource)
			}
		})
	}
}
