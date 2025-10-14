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
		{
			name: "command with custom success message",
			action: &Action{
				ID:             "test-custom-message",
				Type:           ActionTypeCommand,
				SuccessMessage: "Custom success: Installation completed successfully!",
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {Command: "echo test"},
				},
			},
			mockOS:  platform.Linux,
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
				} else if tt.shouldSkip && result != "Command or files already exist, skipping installation" {
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
				if tt.shouldSkip && result != "Command or files already exist, skipping installation" {
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

// TestCommandExists tests the commandExists function with various scenarios
func TestCommandExists(t *testing.T) {
	// Temporarily disable test mode to test actual command checking
	os.Unsetenv(testEnvVar)
	defer os.Setenv(testEnvVar, testEnvValue)

	tests := []struct {
		name     string
		command  string
		expected bool // Note: This will vary based on the system
	}{
		{
			name:     "simple command that should exist - echo",
			command:  "echo",
			expected: true,
		},
		{
			name:     "command with arguments - echo test",
			command:  "echo test",
			expected: true,
		},
		{
			name:     "command with complex arguments - echo hello world",
			command:  "echo hello world",
			expected: true,
		},
		{
			name:     "nonexistent command",
			command:  "thiscommanddoesnotexist12345",
			expected: false,
		},
		{
			name:     "nonexistent command with arguments",
			command:  "thiscommanddoesnotexist12345 version",
			expected: false,
		},
		{
			name:     "complex shell command with test",
			command:  "test -d /tmp",
			expected: true,
		},
		{
			name:     "complex shell command with test on nonexistent directory",
			command:  "test -d /thisdirectorydoesnotexist12345",
			expected: false,
		},
		{
			name:     "docker compose version (WSL scenario)",
			command:  "docker compose version",
			expected: true, // Will be true if Docker is available, false otherwise - we just verify it doesn't crash
		},
		{
			name:     "complex command with operators",
			command:  "which echo && echo found",
			expected: true,
		},
		{
			name:     "command with OR operator",
			command:  "false || true",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := commandExists(tt.command)

			// For commands that should exist, we expect them to actually exist
			// For commands that shouldn't exist, we expect them to not exist
			// For Docker commands, we just verify the function doesn't crash and handles them properly
			if tt.name == "docker compose version (WSL scenario)" {
				// Just log the result for Docker commands - don't fail the test
				t.Logf("Docker compose version check result: %t (environment dependent)", result)
			} else if tt.expected && !result {
				t.Errorf("expected command %q to exist, but it didn't", tt.command)
			} else if !tt.expected && result {
				t.Errorf("expected command %q to not exist, but it did", tt.command)
			}
		})
	}
}

// TestExecutor_CommandExists tests the public CommandExists method
func TestExecutor_CommandExists(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	executor := NewExecutor()

	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "simple command that should exist - echo",
			command:  "echo",
			expected: true,
		},
		{
			name:     "command with arguments - echo test",
			command:  "echo test",
			expected: true,
		},
		{
			name:     "any command in test mode returns true",
			command:  "nonexistentcommand12345",
			expected: true, // In test mode, all commands return true
		},
		{
			name:     "complex shell command with test",
			command:  "test -d /tmp",
			expected: true,
		},
		{
			name:     "any test command in test mode returns true",
			command:  "test -d /nonexistentdirectory12345",
			expected: true, // In test mode, all commands return true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.CommandExists(tt.command)
			if result != tt.expected {
				t.Errorf("CommandExists(%q) = %v, want %v", tt.command, result, tt.expected)
			}
		})
	}
}

// TestExecutor_DockerComposeWSLHandling tests Docker Compose WSL integration handling
func TestExecutor_DockerComposeWSLHandling(t *testing.T) {
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
			name: "docker compose install with check command",
			action: &Action{
				ID:   "test-docker-compose",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command:      "echo installing docker-compose",
						CheckCommand: "docker compose version",
					},
				},
			},
			mockOS:  platform.Linux,
			wantErr: false,
			// Should not skip in test environment since docker won't be available
			shouldSkip: false,
		},
		{
			name: "docker compose install with fallback",
			action: &Action{
				ID:   "test-docker-compose-fallback",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformDebian: {
						Command:      "echo installing docker-compose on debian",
						CheckCommand: "docker compose version",
					},
				},
			},
			mockOS:  platform.Debian,
			wantErr: false,
			// Should not skip in test environment since docker won't be available
			shouldSkip: false,
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
				if result == "" {
					t.Error("expected non-empty result")
				} else if tt.shouldSkip && result != "Command or files already exist, skipping installation" {
					t.Errorf("expected skip message, got %q", result)
				}
			}
		})
	}
}

func TestExecutor_GetCommandString(t *testing.T) {
	executor := NewExecutor()

	tests := []struct {
		name        string
		action      *Action
		wantCmd     string
		wantFound   bool
		description string
	}{
		{
			name: "action with command for current platform",
			action: &Action{
				ID:   "test_action",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					Platform(executor.GetPlatformInfo().OS): {
						Command: "echo test",
					},
				},
			},
			wantCmd:     "echo test",
			wantFound:   true,
			description: "should return command for current platform",
		},
		{
			name: "action with PlatformLinux fallback for Debian",
			action: &Action{
				ID:   "test_action",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command: "echo linux",
					},
				},
			},
			wantCmd:     "echo linux",
			wantFound:   true,
			description: "should fallback to PlatformLinux for Debian/Arch",
		},
		{
			name: "action without command for platform",
			action: &Action{
				ID:               "test_action",
				Type:             ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{},
			},
			wantCmd:     "",
			wantFound:   false,
			description: "should return not found for missing platform",
		},
		{
			name: "action with complex shell command",
			action: &Action{
				ID:   "test_action",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					Platform(executor.GetPlatformInfo().OS): {
						Command: "curl -fsSL https://example.com | sh",
					},
				},
			},
			wantCmd:     "curl -fsSL https://example.com | sh",
			wantFound:   true,
			description: "should return complex shell command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotFound := executor.GetCommandString(tt.action)

			if gotFound != tt.wantFound {
				t.Errorf("GetCommandString() found = %v, want %v", gotFound, tt.wantFound)
			}

			if gotCmd != tt.wantCmd {
				t.Errorf("GetCommandString() cmd = %v, want %v", gotCmd, tt.wantCmd)
			}
		})
	}
}

func TestExecutor_IsInteractiveCommand(t *testing.T) {
	executor := NewExecutor()

	tests := []struct {
		name        string
		action      *Action
		want        bool
		description string
	}{
		{
			name: "interactive command",
			action: &Action{
				ID:   "install_helix",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					Platform(executor.GetPlatformInfo().OS): {
						Command:     "sudo pacman -S helix",
						Interactive: true,
					},
				},
			},
			want:        true,
			description: "should detect interactive command",
		},
		{
			name: "non-interactive command",
			action: &Action{
				ID:   "install_something",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					Platform(executor.GetPlatformInfo().OS): {
						Command:     "echo test",
						Interactive: false,
					},
				},
			},
			want:        false,
			description: "should detect non-interactive command",
		},
		{
			name: "command with no Interactive flag (defaults to false)",
			action: &Action{
				ID:   "install_something",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					Platform(executor.GetPlatformInfo().OS): {
						Command: "echo test",
					},
				},
			},
			want:        false,
			description: "should default to non-interactive",
		},
		{
			name: "interactive command with PlatformLinux fallback",
			action: &Action{
				ID:   "install_interactive",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					PlatformLinux: {
						Command:     "sudo apt install something",
						Interactive: true,
					},
				},
			},
			want:        true,
			description: "should detect interactive with fallback platform",
		},
		{
			name: "action without platform command",
			action: &Action{
				ID:               "test_action",
				Type:             ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{},
			},
			want:        false,
			description: "should return false for missing platform",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executor.IsInteractiveCommand(tt.action)

			if got != tt.want {
				t.Errorf("IsInteractiveCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecutor_InteractiveCommandExecution(t *testing.T) {
	executor := NewExecutor()
	ctx := context.Background()

	tests := []struct {
		name           string
		action         *Action
		wantErrContain string
		description    string
	}{
		{
			name: "interactive command returns special error",
			action: &Action{
				ID:   "install_helix",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					Platform(executor.GetPlatformInfo().OS): {
						Command:     "sudo pacman -S helix",
						Interactive: true,
					},
				},
			},
			wantErrContain: "INTERACTIVE_COMMAND:",
			description:    "should return INTERACTIVE_COMMAND error for interactive actions",
		},
		{
			name: "interactive command includes command string in error",
			action: &Action{
				ID:   "install_something",
				Type: ActionTypeCommand,
				PlatformCommands: map[Platform]PlatformCommand{
					Platform(executor.GetPlatformInfo().OS): {
						Command:     "sudo apt install something",
						Interactive: true,
					},
				},
			},
			wantErrContain: "sudo apt install something",
			description:    "should include command string in error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executor.Execute(ctx, tt.action)

			if err == nil {
				t.Errorf("Execute() expected error for interactive command, got nil")
				return
			}

			if !strings.Contains(err.Error(), tt.wantErrContain) {
				t.Errorf("Execute() error = %v, should contain %v", err.Error(), tt.wantErrContain)
			}
		})
	}
}

// TestCommandExists_WSLDockerIntegration tests WSL Docker integration specific scenarios
func TestCommandExists_WSLDockerIntegration(t *testing.T) {
	// Temporarily disable test mode to test actual command checking
	os.Unsetenv(testEnvVar)
	defer os.Setenv(testEnvVar, testEnvValue)

	// Test WSL Docker integration error handling
	t.Run("WSL Docker integration error detection", func(t *testing.T) {
		// This test simulates the WSL Docker integration issue
		// In a real WSL environment without Docker Desktop integration,
		// "docker compose version" would return the WSL error message

		// Test that the function properly handles commands with arguments
		result := commandExists("docker compose version")

		// In most test environments, docker won't be available, so this should return false
		// The important thing is that it doesn't panic and handles the command parsing correctly
		if result {
			t.Log("Docker appears to be available in this environment")
		} else {
			t.Log("Docker is not available in this environment (expected in CI/test environments)")
		}
	})

	t.Run("command parsing for docker compose", func(t *testing.T) {
		// Test that commands with multiple arguments are parsed correctly
		testCases := []struct {
			name    string
			command string
		}{
			{"docker compose version", "docker compose version"},
			{"docker version", "docker version"},
			{"docker info", "docker info"},
			{"complex docker command", "docker run --rm hello-world"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Should not panic and should handle the command properly
				result := commandExists(tc.command)
				t.Logf("Command: %s, Result: %t", tc.command, result)
			})
		}
	})
}
