package actions

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"anthodev/codory/internal/platform"
)

// Test constants
const (
	testEnvVar   = "CODORY_TEST"
	testEnvValue = "1"
)

// setupTestEnvironment configures the test environment with proper cleanup
func setupTestEnvironment(t *testing.T) func() {
	t.Helper()

	// Set test mode
	platform.SetTestMode(true)
	os.Setenv(testEnvVar, testEnvValue)

	// Return cleanup function
	return func() {
		platform.SetTestMode(false)
		os.Unsetenv(testEnvVar)
	}
}

// createTestPlatform creates platform info for testing
func createTestPlatform(os platform.Platform) platform.Info {
	return platform.Info{
		OS:              os,
		PackageManagers: []platform.PackageManager{},
	}
}

// TestNewExecutor verifies executor creation
func TestNewExecutor(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	executor := NewExecutor()

	if executor == nil {
		t.Fatal("NewExecutor() returned nil")
	}
	if executor.platformInfo.OS == "" {
		t.Error("expected platform info to be set")
	}
	if executor.yayInstaller == nil {
		t.Error("expected yayInstaller to be set")
	}
	if executor.brewInstaller == nil {
		t.Error("expected brewInstaller to be set")
	}
}

// TestExecutor_Execute tests the main Execute method
func TestExecutor_Execute(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

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
			name: "function action with error",
			action: &Action{
				ID:   "test-function-error",
				Type: ActionTypeFunction,
				Handler: func(ctx context.Context) (string, error) {
					return "", errors.New("handler error")
				},
			},
			wantErr:     true,
			errContains: "handler error",
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
					t.Errorf("expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.wantResult {
					t.Errorf("expected result %q, got %q", tt.wantResult, result)
				}
			}
		})
	}
}

// TestExecutor_ExecuteCommand tests command execution
func TestExecutor_ExecuteCommand(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	tests := []struct {
		name        string
		action      *Action
		mockOS      platform.Platform
		wantErr     bool
		errContains string
		shouldSkip  bool // indicates if command should be skipped due to checkCommand
	}{
		{
			name: "simple command execution",
			action: &Action{
				ID:   "test-simple-command",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "echo test"},
				},
			},
			mockOS:  platform.Linux,
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
			mockOS:  platform.Linux,
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
			mockOS:      platform.Linux,
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
			mockOS:      platform.Linux,
			wantErr:     true,
			errContains: "empty command",
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
			mockOS:      platform.Linux,
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
			mockOS:     platform.Linux,
			wantErr:    false,
			shouldSkip: true,
		},
		{
			name: "complex shell command with pipe",
			action: &Action{
				ID:   "test-complex-pipe",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command: "echo hello | grep hello",
					},
				},
			},
			mockOS:  platform.Linux,
			wantErr: false,
		},
		{
			name: "complex shell command with variable expansion",
			action: &Action{
				ID:   "test-complex-variable",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command: "echo $HOME",
					},
				},
			},
			mockOS:  platform.Linux,
			wantErr: false,
		},
		{
			name: "command with PlatformLinux fallback for Debian",
			action: &Action{
				ID:   "test-command-linux-fallback-debian",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "echo linux fallback"},
				},
			},
			mockOS:  platform.Debian,
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
			mockOS:  platform.Arch,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &Executor{
				platformInfo: createTestPlatform(tt.mockOS),
			}

			result, err := executor.executeCommand(context.Background(), tt.action)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == "" {
					t.Error("expected non-empty result")
				} else if tt.shouldSkip && result != "Command already exists, skipping installation" {
					t.Errorf("expected skip message, got %q", result)
				}
			}
		})
	}
}

// TestExecutor_ComplexShellCommands tests complex shell command execution
func TestExecutor_ComplexShellCommands(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	tests := []struct {
		name       string
		action     *Action
		mockOS     platform.Platform
		wantErr    bool
		shouldSkip bool
	}{
		{
			name: "Oh My Zsh install command simulation",
			action: &Action{
				ID:   "test-omz-install",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command:       `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`,
						CheckCommand:  "test -d $HOME/.oh-my-zsh || which omz",
						PackageSource: PackageSourceAny,
					},
				},
			},
			mockOS:  platform.Linux,
			wantErr: false,
		},
		{
			name: "command with directory check using test",
			action: &Action{
				ID:   "test-dir-check",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command:      "echo creating directory",
						CheckCommand: "test -d /tmp",
					},
				},
			},
			mockOS:     platform.Linux,
			wantErr:    false,
			shouldSkip: true,
		},
		{
			name: "command with complex check that requires shell",
			action: &Action{
				ID:   "test-complex-check",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command:      "echo should skip",
						CheckCommand: "test -d /tmp && echo exists",
					},
				},
			},
			mockOS:     platform.Linux,
			wantErr:    false,
			shouldSkip: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &Executor{
				platformInfo: createTestPlatform(tt.mockOS),
			}

			result, err := executor.executeCommand(context.Background(), tt.action)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.shouldSkip && result != "Command already exists, skipping installation" {
					t.Errorf("expected skip message, got %q", result)
				}
			}
		})
	}
}

// TestExecutor_CheckPackageManagerDependency tests package manager dependency checks
func TestExecutor_CheckPackageManagerDependency(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

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
				platformInfo: createTestPlatform(platform.Debian),
			},
			source:      PackageSourceAUR,
			wantErr:     true,
			errContains: "AUR packages are only available on Arch Linux",
		},
		{
			name: "AUR on Arch with yay installed",
			executor: &Executor{
				platformInfo: createTestPlatform(platform.Arch),
			},
			source:  PackageSourceAUR,
			wantErr: false,
			skipFunc: func() bool {
				return !platform.IsYayInstalled()
			},
		},
		{
			name: "AUR on Arch without yay installed",
			executor: &Executor{
				platformInfo: createTestPlatform(platform.Arch),
			},
			source:      PackageSourceAUR,
			wantErr:     true,
			errContains: "yay is not installed",
			skipFunc: func() bool {
				return platform.IsYayInstalled()
			},
		},
		{
			name: "official package source",
			executor: &Executor{
				platformInfo: createTestPlatform(platform.Debian),
			},
			source:  PackageSourceOfficial,
			wantErr: false,
		},
		{
			name: "brew package source",
			executor: &Executor{
				platformInfo: createTestPlatform(platform.MacOS),
			},
			source:  PackageSourceBrew,
			wantErr: false,
		},
		{
			name: "unknown package source",
			executor: &Executor{
				platformInfo: createTestPlatform(platform.Debian),
			},
			source:  PackageSource("unknown"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipFunc != nil && tt.skipFunc() {
				t.Skip("skipping test based on system configuration")
			}

			err := tt.executor.checkPackageManagerDependency(context.Background(), tt.source)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestExecutor_NeedsPackageManagerInstallation tests package manager installation needs
func TestExecutor_NeedsPackageManagerInstallation(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

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
				platformInfo: createTestPlatform(platform.Arch),
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
				return platform.IsYayInstalled()
			},
		},
		{
			name: "AUR action on non-Arch platform",
			executor: &Executor{
				platformInfo: createTestPlatform(platform.Debian),
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
				platformInfo: createTestPlatform(platform.Debian),
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
				platformInfo: createTestPlatform(platform.Debian),
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
				platformInfo: createTestPlatform(platform.Arch),
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
				return platform.IsYayInstalled()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipFunc != nil && tt.skipFunc() {
				t.Skip("skipping test based on system configuration")
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

// TestExecutor_GetPlatformInfo tests platform info retrieval
func TestExecutor_GetPlatformInfo(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	executor := NewExecutor()
	info := executor.GetPlatformInfo()

	if info.OS == "" {
		t.Error("expected platform info to have OS set")
	}
	if info.PackageManagers == nil {
		t.Error("expected platform info to have PackageManagers set")
	}
}
