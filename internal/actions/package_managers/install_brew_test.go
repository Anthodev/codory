package package_managers

import (
	"context"
	"os"
	"testing"

	"anthodev/codory/internal/platform"
)

func TestNewInstallBrewAction(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	action := NewInstallBrewAction()

	if action == nil {
		t.Fatal("NewInstallBrewAction() returned nil")
	}

	if action.ID != "install_brew" {
		t.Errorf("Expected action ID to be 'install_brew', got '%s'", action.ID)
	}

	if action.Name != "Install Homebrew" {
		t.Errorf("Expected action name to be 'Install Homebrew', got '%s'", action.Name)
	}

	if action.Description != "Install Homebrew on the system" {
		t.Errorf("Expected action description to be 'Install Homebrew on the system', got '%s'", action.Description)
	}

	if action.Type != "function" {
		t.Errorf("Expected action type to be 'function', got '%s'", action.Type)
	}

	if action.Handler == nil {
		t.Error("Expected action handler to be set")
	}

	if len(action.HiddenOnPlatforms) != 1 {
		t.Fatalf("Expected 1 hidden platform, got %d", len(action.HiddenOnPlatforms))
	}

	if action.HiddenOnPlatforms[0] != "windows" {
		t.Errorf("Expected action to be hidden on 'windows' platform, got '%s'", action.HiddenOnPlatforms[0])
	}
}

func TestInstallBrew_ValidationOnly(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	tests := []struct {
		name          string
		mockOS        platform.Platform
		brewInstalled bool
		wantErr       bool
		errContains   string
		wantResult    string
	}{
		{
			name:          "install on Windows platform",
			mockOS:        platform.Windows,
			brewInstalled: false,
			wantErr:       true,
			errContains:   "brew can not be installed on Windows",
			wantResult:    "",
		},
		{
			name:          "brew already installed",
			mockOS:        platform.MacOS,
			brewInstalled: true,
			wantErr:       false,
			wantResult:    "Homebrew is already installed",
		},
		{
			name:          "brew not installed - blocked in test environment",
			mockOS:        platform.MacOS,
			brewInstalled: false,
			wantErr:       true,
			errContains:   "skipping brew installation in test environment",
			wantResult:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS := platform.Detect()
			if currentOS != tt.mockOS {
				t.Skipf("Skipping test: expected OS %s, but running on %s", tt.mockOS, currentOS)
			}

			// For the "brew already installed" test, skip if brew is not actually installed
			if tt.name == "brew already installed" && !platform.IsBrewInstalled() {
				t.Skip("Skipping 'brew already installed' test because brew is not installed")
			}

			// For the "brew not installed" test, skip if brew is already installed
			if tt.name == "brew not installed - blocked in test environment" && platform.IsBrewInstalled() {
				t.Skip("Skipping 'brew not installed' test because brew is already installed")
			}

			result, err := installBrew(context.Background())

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

func TestInstallBrew_PlatformValidationOnly(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// This test only runs on supported platforms to validate the early validation logic
	if platform.Detect() == platform.Windows {
		t.Skip("Skipping validation test on Windows system")
	}

	// Test the validation logic without actually attempting installation
	if platform.IsBrewInstalled() {
		result, err := installBrew(context.Background())
		if err != nil {
			t.Errorf("Expected no error when brew is already installed, got: %v", err)
		}
		expectedResult := "Homebrew is already installed"
		if result != expectedResult {
			t.Errorf("Expected result '%s', got '%s'", expectedResult, result)
		}
	} else {
		// If brew is not installed, verify that installation is blocked in test environment
		result, err := installBrew(context.Background())
		if err == nil {
			t.Error("Expected error when trying to install brew in test environment, got none")
		}
		if !contains(err.Error(), "skipping brew installation in test environment") {
			t.Errorf("Expected error to contain 'skipping brew installation in test environment', got: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty result when installation is blocked, got: %s", result)
		}
	}
}

func TestValidateInstallBrewRequirements(t *testing.T) {
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
			name:        "validation on Windows platform",
			mockOS:      platform.Windows,
			wantErr:     true,
			errContains: "brew can not be installed on Windows",
		},
		{
			name:        "validation on supported platform",
			mockOS:      platform.MacOS,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "validation on Linux platform",
			mockOS:      platform.Debian,
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

			err := validateInstallBrewRequirements()

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
