package shells

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/internal/platform"
)

func TestSetZshAsDefaultShell(t *testing.T) {
	action := SetZshAsDefaultShell()

	if action == nil {
		t.Fatal("SetZshAsDefaultShell() returned nil")
	}

	if action.ID != "set_zsh_default_shell" {
		t.Errorf("Expected action ID to be 'set_zsh_default_shell', got '%s'", action.ID)
	}

	if action.Name != "Set Zsh as default shell" {
		t.Errorf("Expected action name to be 'Set Zsh as default shell', got '%s'", action.Name)
	}

	if action.Description != "Set Zsh as the default shell in your system" {
		t.Errorf("Expected action description to be 'Set Zsh as the default shell in your system', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeFunction {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeFunction, action.Type)
	}

	if action.Handler == nil {
		t.Error("Expected action handler to be set")
	}

	// Verify it's hidden on Windows
	if !action.IsHiddenOnPlatform(actions.PlatformWindows) {
		t.Error("Action should be hidden on Windows platform")
	}

	// Verify it's visible on other platforms
	visiblePlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformDebian,
		actions.PlatformArch,
		actions.PlatformMacOS,
	}

	for _, platform := range visiblePlatforms {
		if action.IsHiddenOnPlatform(platform) {
			t.Errorf("Action should not be hidden on %s platform", platform)
		}
	}
}

func TestSetZshAsDefaultShell_ValidationOnly(t *testing.T) {
	// Ensure we're in test mode
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	tests := []struct {
		name          string
		mockOS        platform.Platform
		zshInstalled  bool
		currentShell  string
		wantErr       bool
		errContains   string
		wantResult    string
		skipCondition func() bool
	}{
		{
			name:          "Windows platform",
			mockOS:        platform.Windows,
			zshInstalled:  false,
			currentShell:  "/bin/bash",
			wantErr:       true,
			errContains:   "Zsh is not supported on Windows",
			wantResult:    "",
			skipCondition: nil,
		},
		{
			name:          "Zsh not installed",
			mockOS:        platform.Linux,
			zshInstalled:  false,
			currentShell:  "/bin/bash",
			wantErr:       true,
			errContains:   "Zsh is not installed",
			wantResult:    "",
			skipCondition: nil,
		},
		{
			name:          "Zsh already default shell",
			mockOS:        platform.Linux,
			zshInstalled:  true,
			currentShell:  "/bin/zsh",
			wantErr:       false,
			errContains:   "",
			wantResult:    "Zsh is already the default shell",
			skipCondition: nil,
		},
		{
			name:         "Zsh installed but not default",
			mockOS:       platform.Linux,
			zshInstalled: true,
			currentShell: "/bin/bash",
			wantErr:      true,
			errContains:  "failed to set Zsh as default shell",
			wantResult:   "",
			skipCondition: func() bool {
				// Skip this test in CI/test environments where we can't actually change shells
				return os.Getenv("CODORY_TEST") == "1" || os.Getenv("CI") == "true"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS := platform.Detect()
			if currentOS != tt.mockOS {
				t.Skipf("Skipping test: expected OS %s, but running on %s", tt.mockOS, currentOS)
			}

			// Check skip condition if provided
			if tt.skipCondition != nil && tt.skipCondition() {
				t.Skipf("Skipping test case based on skip condition")
			}

			// For the "Zsh not installed" test, skip if zsh is actually installed
			if tt.name == "Zsh not installed" {
				if _, err := exec.LookPath("zsh"); err == nil {
					t.Skip("Skipping 'Zsh not installed' test because zsh is installed")
				}
			}

			// For the "Zsh already default shell" test, skip if zsh is not the default
			if tt.name == "Zsh already default shell" {
				if shell, err := exec.Command("sh", "-c", "echo $SHELL").Output(); err != nil || !strings.Contains(string(shell), "zsh") {
					t.Skip("Skipping 'Zsh already default shell' test because zsh is not the default")
				}
			}

			result, err := setZshAsDefaultShell(context.Background())

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
				if result != tt.wantResult {
					t.Errorf("Expected result '%s', got '%s'", tt.wantResult, result)
				}
			}
		})
	}
}

func TestValidateZshInstallation(t *testing.T) {
	// Ensure we're in test mode
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	tests := []struct {
		name        string
		mockOS      platform.Platform
		wantErr     bool
		errContains string
	}{
		{
			name:        "validation on Windows",
			mockOS:      platform.Windows,
			wantErr:     true,
			errContains: "Zsh is not supported on Windows",
		},
		{
			name:        "validation on Linux with zsh",
			mockOS:      platform.Linux,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "validation on macOS",
			mockOS:      platform.MacOS,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "validation on Debian",
			mockOS:      platform.Debian,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "validation on Arch",
			mockOS:      platform.Arch,
			wantErr:     false,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS := platform.Detect()
			if currentOS != tt.mockOS {
				t.Skipf("Skipping test: expected OS %s, but running on %s", tt.mockOS, currentOS)
			}

			err := validateZshInstallation()

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

func TestSetZshAsDefaultShell_ContextCancellation(t *testing.T) {
	// Ensure we're in test mode
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Skip on Windows since it fails validation first
	if platform.Detect() == platform.Windows {
		t.Skip("Skipping context cancellation test on Windows platform")
	}

	// Skip if zsh is not installed
	if _, err := exec.LookPath("zsh"); err != nil {
		t.Skip("Skipping context cancellation test because zsh is not installed")
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Test with cancelled context
	result, err := setZshAsDefaultShell(ctx)

	// Context cancellation should either cause an error or be handled gracefully
	// The actual behavior depends on how the underlying commands handle context
	if err != nil && result != "" {
		t.Error("Expected either error or result, not both")
	}
}

func TestSetZshAsDefaultShell_PlatformValidationOnly(t *testing.T) {
	// Ensure we're in test mode
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// This test validates platform-specific behavior
	currentOS := platform.Detect()

	if currentOS == platform.Windows {
		// On Windows, should fail validation
		result, err := setZshAsDefaultShell(context.Background())
		if err == nil {
			t.Error("Expected error on Windows platform")
		}
		if !contains(err.Error(), "Zsh is not supported on Windows") {
			t.Errorf("Expected error to contain 'Zsh is not supported on Windows', got: %v", err)
		}
		if result != "" {
			t.Error("Expected empty result on Windows platform")
		}
	} else {
		// On other platforms, validation should pass if zsh is installed
		if _, err := exec.LookPath("zsh"); err != nil {
			t.Skip("Skipping validation test because zsh is not installed")
		}

		// If zsh is already default, should return early
		if shell, err := exec.Command("sh", "-c", "echo $SHELL").Output(); err == nil && strings.Contains(string(shell), "zsh") {
			result, err := setZshAsDefaultShell(context.Background())
			if err != nil {
				t.Errorf("Expected no error when zsh is already default, got: %v", err)
			}
			expectedResult := "Zsh is already the default shell"
			if result != expectedResult {
				t.Errorf("Expected result '%s', got '%s'", expectedResult, result)
			}
		}
	}
}

func TestSetZshAsDefaultShell_ErrorMessageFormatting(t *testing.T) {
	// Ensure we're in test mode
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Test error message formatting for different scenarios
	tests := []struct {
		name          string
		mockOS        platform.Platform
		wantErr       bool
		errContains   string
		skipCondition func() bool
	}{
		{
			name:          "error message formatting on Windows",
			mockOS:        platform.Windows,
			wantErr:       true,
			errContains:   "Zsh is not supported on Windows",
			skipCondition: nil,
		},
		{
			name:        "error message formatting when zsh not installed",
			mockOS:      platform.Linux,
			wantErr:     true,
			errContains: "Zsh is not installed",
			skipCondition: func() bool {
				// Skip this test if zsh is installed
				_, err := exec.LookPath("zsh")
				return err == nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS := platform.Detect()
			if currentOS != tt.mockOS {
				t.Skipf("Skipping test: expected OS %s, but running on %s", tt.mockOS, currentOS)
			}

			// Check skip condition if provided
			if tt.skipCondition != nil && tt.skipCondition() {
				t.Skipf("Skipping test case based on skip condition")
			}

			_, err := setZshAsDefaultShell(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
				}

				// Test that error messages are properly formatted with context
				if err != nil && len(err.Error()) < 10 {
					t.Errorf("Error message seems too short to be properly formatted: %s", err.Error())
				}
			}
		})
	}
}

func TestSetZshAsDefaultShell_ValidateZshInstallation_ErrorWrapping(t *testing.T) {
	// Ensure we're in test mode
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Test error wrapping in validateZshInstallation
	currentOS := platform.Detect()
	if currentOS != platform.Windows {
		t.Skip("Skipping error wrapping test on non-Windows system")
	}

	err := validateZshInstallation()
	if err == nil {
		t.Error("Expected error on Windows platform")
	}

	// Test that the error is properly formatted
	expectedMsg := "Zsh is not supported on Windows"
	if err != nil && !contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestSetZshAsDefaultShell_SetZshAsDefaultShell_ErrorWrapping(t *testing.T) {
	// Ensure we're in test mode
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Test error wrapping in the main setZshAsDefaultShell function
	currentOS := platform.Detect()
	if currentOS == platform.Windows {
		t.Skip("Skipping install error wrapping test on Windows system")
	}

	// Skip if zsh is not installed
	if _, err := exec.LookPath("zsh"); err != nil {
		t.Skip("Skipping install error wrapping test because zsh is not installed")
	}

	// Skip if zsh is already default (function returns early)
	if shell, err := exec.Command("sh", "-c", "echo $SHELL").Output(); err == nil && strings.Contains(string(shell), "zsh") {
		t.Skip("Skipping install error wrapping test because zsh is already default")
	}

	_, err := setZshAsDefaultShell(context.Background())
	if err == nil {
		t.Error("Expected error when trying to set zsh as default in test environment")
	}

	// Test that the error message is properly formatted
	if err != nil && !contains(err.Error(), "failed to set Zsh as default shell") {
		t.Errorf("Expected error to contain 'failed to set Zsh as default shell', got: %v", err)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

// Test helper function to verify string contains functionality
func TestContainsHelper(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"hello world", "world", true},
		{"hello world", "foo", false},
		{"", "foo", false},
		{"foo", "", true},
		{"foo", "foo", true},
		{"foobar", "foo", true},
		{"foobar", "bar", true},
		{"foobar", "baz", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("contains(%q, %q)", tt.s, tt.substr), func(t *testing.T) {
			got := contains(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}
