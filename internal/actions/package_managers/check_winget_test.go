package package_managers

import (
	"context"
	"os"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/internal/platform"
)

func TestNewCheckWingetAction(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	action := NewCheckWingetAction()

	if action == nil {
		t.Fatal("NewCheckWingetAction() returned nil")
	}

	if action.ID != "check_winget" {
		t.Errorf("Expected action ID to be 'check_winget', got '%s'", action.ID)
	}

	if action.Name != "Check Winget" {
		t.Errorf("Expected action name to be 'Check Winget', got '%s'", action.Name)
	}

	if action.Description != "Check if Winget is installed and show installation instructions if not installed" {
		t.Errorf("Expected action description to be 'Check if Winget is installed and show installation instructions if not installed', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeFunction {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeFunction, action.Type)
	}

	if action.Handler == nil {
		t.Error("Expected action handler to be set")
	}

	if len(action.VisibleOnPlatforms) != 1 {
		t.Fatalf("Expected 1 visible platform, got %d", len(action.VisibleOnPlatforms))
	}

	if action.VisibleOnPlatforms[0] != actions.PlatformWindows {
		t.Errorf("Expected action to be visible on '%s' platform, got '%s'", actions.PlatformWindows, action.VisibleOnPlatforms[0])
	}
}

func TestCheckWinget_ValidationOnly(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	tests := []struct {
		name        string
		mockOS      platform.Platform
		wantErr     bool
		errContains string
		wantResult  string
	}{
		{
			name:        "check on non-Windows platform",
			mockOS:      platform.Arch,
			wantErr:     true,
			errContains: "winget is not supported on this platform",
			wantResult:  "",
		},
		{
			name:        "check on macOS platform",
			mockOS:      platform.MacOS,
			wantErr:     true,
			errContains: "winget is not supported on this platform",
			wantResult:  "",
		},
		{
			name:        "check on Windows platform - validation only",
			mockOS:      platform.Windows,
			wantErr:     false,
			errContains: "",
			wantResult:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS := platform.Detect()
			if currentOS != tt.mockOS {
				t.Skipf("Skipping test: expected OS %s, but running on %s", tt.mockOS, currentOS)
			}

			result, err := checkWinget(context.Background())

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
				// For Windows platform, we don't check the exact result since it depends on winget installation status
				if tt.wantResult != "" && result != tt.wantResult {
					t.Errorf("Expected result '%s', got '%s'", tt.wantResult, result)
				}
			}
		})
	}
}

func TestCheckWinget_PlatformValidationOnly(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// This test only runs on non-Windows platforms to validate the early validation logic
	if platform.Detect() == platform.Windows {
		t.Skip("Skipping validation test on Windows system")
	}

	// Test the validation logic - should fail on non-Windows platforms
	result, err := checkWinget(context.Background())
	if err == nil {
		t.Error("Expected error when checking winget on non-Windows platform, got none")
	}
	if !contains(err.Error(), "winget is not supported on this platform") {
		t.Errorf("Expected error to contain 'winget is not supported on this platform', got: %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty result when validation fails, got: %s", result)
	}
}

func TestValidateCheckWinget(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	tests := []struct {
		name        string
		mockOS      platform.Platform
		wantErr     bool
		errContains string
	}{
		{
			name:        "validation on Linux platform",
			mockOS:      platform.Arch,
			wantErr:     true,
			errContains: "winget is not supported on this platform",
		},
		{
			name:        "validation on macOS platform",
			mockOS:      platform.MacOS,
			wantErr:     true,
			errContains: "winget is not supported on this platform",
		},
		{
			name:        "validation on Debian platform",
			mockOS:      platform.Debian,
			wantErr:     true,
			errContains: "winget is not supported on this platform",
		},
		{
			name:        "validation on Arch platform",
			mockOS:      platform.Arch,
			wantErr:     true,
			errContains: "winget is not supported on this platform",
		},
		{
			name:        "validation on Windows platform",
			mockOS:      platform.Windows,
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

			err := validateCheckWinget(context.Background())

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

func TestCheckWinget_HandlerIntegration(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	action := NewCheckWingetAction()

	if action.Handler == nil {
		t.Fatal("Action handler is nil")
	}

	// Test that the handler can be called
	ctx := context.Background()
	result, err := action.Handler(ctx)

	// On non-Windows platforms, should return an error
	if platform.Detect() != platform.Windows {
		if err == nil {
			t.Error("Expected error when running winget check on non-Windows platform")
		}
		if result != "" {
			t.Error("Expected empty result when validation fails")
		}
		return
	}

	// On Windows platform, we just verify the handler runs without panicking
	// The actual result depends on whether winget is installed or not
	if err != nil && result == "" {
		t.Error("Expected result to contain installation instructions when winget is not installed")
	}
	if err == nil && result == "" {
		t.Error("Expected result when winget check succeeds")
	}
}

func TestCheckWinget_ActionProperties(t *testing.T) {
	action := NewCheckWingetAction()

	tests := []struct {
		name     string
		property string
		expected interface{}
		actual   interface{}
	}{
		{
			name:     "Action ID",
			property: "ID",
			expected: "check_winget",
			actual:   action.ID,
		},
		{
			name:     "Action Name",
			property: "Name",
			expected: "Check Winget",
			actual:   action.Name,
		},
		{
			name:     "Action Description",
			property: "Description",
			expected: "Check if Winget is installed and show installation instructions if not installed",
			actual:   action.Description,
		},
		{
			name:     "Action Type",
			property: "Type",
			expected: actions.ActionTypeFunction,
			actual:   action.Type,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("Expected %s to be '%v', got '%v'", tt.property, tt.expected, tt.actual)
			}
		})
	}

	// Test platform visibility
	if len(action.VisibleOnPlatforms) != 1 {
		t.Errorf("Expected 1 visible platform, got %d", len(action.VisibleOnPlatforms))
	} else if action.VisibleOnPlatforms[0] != actions.PlatformWindows {
		t.Errorf("Expected visible platform to be '%s', got '%s'", actions.PlatformWindows, action.VisibleOnPlatforms[0])
	}

	// Test that there are no hidden platforms
	if len(action.HiddenOnPlatforms) != 0 {
		t.Errorf("Expected 0 hidden platforms, got %d", len(action.HiddenOnPlatforms))
	}
}

func TestCheckWinget_ContextCancellation(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Only test on Windows since other platforms fail validation first
	if platform.Detect() != platform.Windows {
		t.Skip("Skipping context cancellation test on non-Windows platform")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result, err := checkWinget(ctx)

	// Should handle cancelled context gracefully
	if err != nil && result == "" {
		t.Log("Context cancellation handled appropriately")
	}
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
