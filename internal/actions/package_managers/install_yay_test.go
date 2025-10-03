package package_managers

import (
	"context"
	"os"
	"testing"

	"anthodev/codory/internal/platform"
)

func TestNewInstallYayAction(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	action := NewInstallYayAction()

	if action == nil {
		t.Fatal("NewInstallYayAction() returned nil")
	}

	if action.ID != "install_yay" {
		t.Errorf("Expected action ID to be 'install_yay', got '%s'", action.ID)
	}

	if action.Name != "Install Yay (AUR Helper)" {
		t.Errorf("Expected action name to be 'Install Yay (AUR Helper)', got '%s'", action.Name)
	}

	if action.Description != "Install Yay AUR helper for Arch Linux" {
		t.Errorf("Expected action description to be 'Install Yay AUR helper for Arch Linux', got '%s'", action.Description)
	}

	if action.Type != "function" {
		t.Errorf("Expected action type to be 'function', got '%s'", action.Type)
	}

	if action.Handler == nil {
		t.Error("Expected action handler to be set")
	}

	if len(action.VisibleOnPlatforms) != 1 {
		t.Fatalf("Expected 1 visible platform, got %d", len(action.VisibleOnPlatforms))
	}

	if action.VisibleOnPlatforms[0] != "arch" {
		t.Errorf("Expected action to be visible on 'arch' platform, got '%s'", action.VisibleOnPlatforms[0])
	}
}

func TestInstallYay_ValidationOnly(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	tests := []struct {
		name         string
		mockOS       platform.Platform
		yayInstalled bool
		wantErr      bool
		errContains  string
		wantResult   string
	}{
		{
			name:         "install on non-Arch platform",
			mockOS:       platform.Debian,
			yayInstalled: false,
			wantErr:      true,
			errContains:  "yay can only be installed on Arch Linux",
			wantResult:   "",
		},
		{
			name:         "yay already installed",
			mockOS:       platform.Arch,
			yayInstalled: true,
			wantErr:      false,
			wantResult:   "Yay is already installed!",
		},
		{
			name:         "yay not installed - blocked in test environment",
			mockOS:       platform.Arch,
			yayInstalled: false,
			wantErr:      true,
			errContains:  "skipping yay installation in test environment",
			wantResult:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS := platform.Detect()
			if currentOS != tt.mockOS {
				t.Skipf("Skipping test: expected OS %s, but running on %s", tt.mockOS, currentOS)
			}

			// For the "yay already installed" test, skip if yay is not actually installed
			if tt.name == "yay already installed" && !platform.IsYayInstalled() {
				t.Skip("Skipping 'yay already installed' test because yay is not installed")
			}

			// For the "yay not installed" test, skip if yay is already installed
			if tt.name == "yay not installed - blocked in test environment" && platform.IsYayInstalled() {
				t.Skip("Skipping 'yay not installed' test because yay is already installed")
			}

			result, err := installYay(context.Background())

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

func TestInstallYay_ArchValidationOnly(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// This test only runs on Arch Linux to validate the early validation logic
	if platform.Detect() != platform.Arch {
		t.Skip("Skipping Arch-specific validation test on non-Arch system")
	}

	// Test the validation logic without actually attempting installation
	if platform.IsYayInstalled() {
		result, err := installYay(context.Background())
		if err != nil {
			t.Errorf("Expected no error when yay is already installed, got: %v", err)
		}
		expectedResult := "Yay is already installed!"
		if result != expectedResult {
			t.Errorf("Expected result '%s', got '%s'", expectedResult, result)
		}
	} else {
		// If yay is not installed, verify that installation is blocked in test environment
		result, err := installYay(context.Background())
		if err == nil {
			t.Error("Expected error when trying to install yay in test environment, got none")
		}
		if !contains(err.Error(), "skipping yay installation in test environment") {
			t.Errorf("Expected error to contain 'skipping yay installation in test environment', got: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty result when installation is blocked, got: %s", result)
		}
	}
}

func TestValidateInstallYayRequirements(t *testing.T) {
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
			name:        "validation on non-Arch platform",
			mockOS:      platform.Debian,
			wantErr:     true,
			errContains: "yay can only be installed on Arch Linux",
		},
		{
			name:        "validation on Arch platform",
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

			err := validateInstallYayRequirements()

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

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
