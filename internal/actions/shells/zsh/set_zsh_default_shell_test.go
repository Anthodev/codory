package shells

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"anthodev/codory/internal/actions"
)

func TestSetZshAsDefaultShell(t *testing.T) {
	action := SetZshAsDefaultShell()

	// Test action properties
	if action == nil {
		t.Fatal("SetZshAsDefaultShell() returned nil")
	}

	if action.ID != "set_zsh_default_shell" {
		t.Errorf("Expected action ID to be 'set_zsh_default_shell', got '%s'", action.ID)
	}

	if action.Name != "Set Zsh as default shell" {
		t.Errorf("Expected action name to be 'Set Zsh as default shell', got '%s'", action.Name)
	}

	if action.Type != actions.ActionTypeFunction {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeFunction, action.Type)
	}

	if action.Handler == nil {
		t.Error("Expected action handler to be set")
	}

	// Test platform visibility
	if !action.IsHiddenOnPlatform(actions.PlatformWindows) {
		t.Error("Action should be hidden on Windows platform")
	}

	// Should be visible on Linux platforms and macOS
	for _, platform := range []actions.Platform{actions.PlatformLinux, actions.PlatformDebian, actions.PlatformArch, actions.PlatformMacOS} {
		if action.IsHiddenOnPlatform(platform) {
			t.Errorf("Action should not be hidden on %s platform", platform)
		}
	}
}

func TestSetZshAsDefaultShell_Integration(t *testing.T) {
	// Set test mode
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Test Windows platform behavior
	if runtime.GOOS == "windows" {
		result, err := setZshAsDefaultShell(context.Background())
		if err == nil {
			t.Error("Expected error on Windows")
		}
		if !strings.Contains(err.Error(), "Zsh is not supported on Windows") {
			t.Errorf("Expected Windows-specific error, got: %v", err)
		}
		if result != "" {
			t.Error("Expected empty result on Windows")
		}
		return
	}

	// Test zsh installation check
	if _, err := exec.LookPath("zsh"); err != nil {
		// Zsh not installed - should fail validation
		result, err := setZshAsDefaultShell(context.Background())
		if err == nil {
			t.Error("Expected error when zsh is not installed")
		}
		if !strings.Contains(err.Error(), "Zsh is not installed") {
			t.Errorf("Expected 'zsh not installed' error, got: %v", err)
		}
		if result != "" {
			t.Error("Expected empty result when zsh not installed")
		}
		return
	}

	// Test current shell check
	if shell, err := exec.Command("sh", "-c", "echo $SHELL").Output(); err == nil && strings.Contains(string(shell), "zsh") {
		// Zsh is already default - should return early
		result, err := setZshAsDefaultShell(context.Background())
		if err != nil {
			t.Errorf("Expected no error when zsh is already default, got: %v", err)
		}
		expected := "Zsh is already the default shell"
		if result != expected {
			t.Errorf("Expected result '%s', got '%s'", expected, result)
		}
		return
	}

	// Test actual shell change attempt (should fail in test environment)
	result, err := setZshAsDefaultShell(context.Background())
	if err == nil {
		t.Error("Expected error when trying to change shell in test environment")
	}
	if !strings.Contains(err.Error(), "failed to set Zsh as default shell") {
		t.Errorf("Expected 'failed to set' error, got: %v", err)
	}
	if result != "" {
		t.Error("Expected empty result when shell change fails")
	}
}

func TestValidateZshInstallation(t *testing.T) {
	tests := []struct {
		name        string
		checkZsh    bool
		expectError bool
		errorMsg    string
	}{
		{
			name:        "zsh installed",
			checkZsh:    true,
			expectError: false,
		},
		{
			name:        "zsh not installed",
			checkZsh:    false,
			expectError: true,
			errorMsg:    "Zsh is not installed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.checkZsh {
				// Skip if zsh is not actually installed
				if _, err := exec.LookPath("zsh"); err != nil {
					t.Skip("zsh is not installed on this system")
				}
			} else {
				// Skip if zsh is installed
				if _, err := exec.LookPath("zsh"); err == nil {
					t.Skip("zsh is installed on this system")
				}
			}

			err := validateZshInstallation()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test error message formatting
	err := fmt.Errorf("test error")
	wrappedErr := fmt.Errorf("context: %w", err)

	if !strings.Contains(wrappedErr.Error(), "context:") {
		t.Error("Error should be wrapped with context")
	}
	if !strings.Contains(wrappedErr.Error(), "test error") {
		t.Error("Original error should be preserved")
	}
}

func TestActionHandler(t *testing.T) {
	action := SetZshAsDefaultShell()

	// Test that handler can be called
	result, err := action.Handler(context.Background())

	// Should either succeed or fail gracefully, but not panic
	if err != nil && result != "" {
		t.Error("Expected either error or result, not both")
	}
}
