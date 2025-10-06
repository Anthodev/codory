package actions

import (
	"context"
	"os"
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

// SetTestMode sets the test mode for the platform package
func setPlatformTestMode(enabled bool) {
	platform.SetTestMode(enabled)
}

// Helper function to create test platform info
func createTestPlatformInfo(os platform.Platform) platform.Info {
	return platform.Info{
		OS:              os,
		PackageManagers: []platform.PackageManager{},
	}
}

func TestNewExecutor(t *testing.T) {
	// Ensure we're in test mode
	setPlatformTestMode(true)
	defer setPlatformTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	executor := NewExecutor()
	if executor == nil {
		t.Fatal("NewExecutor() returned nil")
	}
	if executor.platformInfo.OS == "" {
		t.Error("Expected platform info to be set")
	}
	if executor.yayInstaller == nil {
		t.Error("Expected yayInstaller to be set")
	}
	if executor.brewInstaller == nil {
		t.Error("Expected brewInstaller to be set")
	}
}

func TestExecutor_Execute(t *testing.T) {
	// Ensure we're in test mode
	setPlatformTestMode(true)
	defer setPlatformTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	tests := []struct {
		name        string
		action      *Action
		wantErr     bool
		errContains string
		wantResult  string
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
			wantErr:    false,
			wantResult: "function result",
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
		{
			name: "function action with error",
			action: &Action{
				ID:   "test-function-error",
				Type: ActionTypeFunction,
				Handler: func(ctx context.Context) (string, error) {
					return "", context.DeadlineExceeded
				},
			},
			wantErr:     true,
			errContains: "context deadline exceeded",
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
				if tt.wantResult != "" && result != tt.wantResult {
					t.Errorf("Expected '%s', got '%s'", tt.wantResult, result)
				}
			}
		})
	}
}

func TestExecutor_executeCommand(t *testing.T) {
	// Ensure we're in test mode
	setPlatformTestMode(true)
	defer setPlatformTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

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
		{
			name: "command with package source",
			action: &Action{
				ID:   "test-command-package-source",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command:       "echo test",
						PackageSource: PackageSourceOfficial,
					},
				},
			},
			mockOS:  PlatformLinux,
			wantErr: false,
		},
		{
			name: "command that fails",
			action: &Action{
				ID:   "test-command-fails",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "false"},
				},
			},
			mockOS:      PlatformLinux,
			wantErr:     true,
			errContains: "command failed",
		},
		{
			name: "command with check command - exists",
			action: &Action{
				ID:   "test-check-command-exists",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command:      "echo should-not-run",
						CheckCommand: "echo",
					},
				},
			},
			mockOS:  PlatformLinux,
			wantErr: false,
		},
		{
			name: "command with check command - not exists",
			action: &Action{
				ID:   "test-check-command-not-exists",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command:      "echo test",
						CheckCommand: "nonexistentcommand12345",
					},
				},
			},
			mockOS:  PlatformLinux,
			wantErr: false,
		},
		{
			name: "command with PlatformLinux fallback for Debian",
			action: &Action{
				ID:   "test-command-linux-fallback",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "echo linux fallback"},
				},
			},
			mockOS:  PlatformDebian,
			wantErr: false,
		},
		{
			name: "command with PlatformLinux fallback for Arch",
			action: &Action{
				ID:   "test-command-linux-fallback-arch",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "echo linux fallback arch"},
				},
			},
			mockOS:  PlatformArch,
			wantErr: false,
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
				} else if tt.name == "command with check command - exists" && result != "Command already exists, skipping installation" {
					t.Errorf("Expected skip message, got '%s'", result)
				}
			}
		})
	}
}

func TestExecutor_GetPlatformInfo(t *testing.T) {
	// Ensure we're in test mode
	setPlatformTestMode(true)
	defer setPlatformTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	executor := NewExecutor()
	info := executor.GetPlatformInfo()

	if info.OS == "" {
		t.Error("Expected platform info to have OS set")
	}
	if info.PackageManagers == nil {
		t.Error("Expected platform info to have PackageManagers set")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

func TestExecutor_checkPackageManagerDependency(t *testing.T) {
	// Ensure we're in test mode
	setPlatformTestMode(true)
	defer setPlatformTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

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
			name: "AUR on Arch without yay installed",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.Arch),
			},
			source:      PackageSourceAUR,
			wantErr:     true,
			errContains: "yay is not installed",
			skipFunc: func() bool {
				// Skip this test case if yay is actually installed on the system
				return platform.IsYayInstalled()
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
		{
			name: "unknown package source",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.Debian),
			},
			source:  PackageSource("unknown"),
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
	// Ensure we're in test mode
	setPlatformTestMode(true)
	defer setPlatformTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

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
		{
			name: "AUR action with PlatformLinux fallback for Arch",
			executor: &Executor{
				platformInfo: createMockPlatformInfo(platform.Arch),
			},
			action: &Action{
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
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
