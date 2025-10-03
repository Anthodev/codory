package platform

import (
	"context"
	"os"
	"testing"
)

func TestNewYayInstaller(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	testMode := func() bool { return true }
	installer := NewYayInstaller(testMode)

	if installer == nil {
		t.Fatal("NewYayInstaller() returned nil")
	}
}

func TestYayInstaller_ValidateRequirements(t *testing.T) {
	tests := []struct {
		name        string
		mockOS      Platform
		pacmanExist bool
		gitExist    bool
		wantErr     bool
		errContains string
		skipFunc    func() bool
	}{
		{
			name:        "install on non-Arch platform",
			mockOS:      Debian,
			pacmanExist: true,
			gitExist:    true,
			wantErr:     true,
			errContains: "yay can only be installed on Arch Linux",
		},
		{
			name:        "install on Arch without pacman",
			mockOS:      Arch,
			pacmanExist: false,
			gitExist:    true,
			wantErr:     true,
			errContains: "pacman is not installed",
			skipFunc: func() bool {
				// Skip this test case if pacman is actually installed on the system
				return IsPacmanInstalled()
			},
		},
		{
			name:        "install on Arch without git",
			mockOS:      Arch,
			pacmanExist: true,
			gitExist:    false,
			wantErr:     true,
			errContains: "git is required to install yay",
			skipFunc: func() bool {
				// Skip this test case if git is actually installed on the system
				return commandExists("git")
			},
		},
		{
			name:        "valid Arch setup",
			mockOS:      Arch,
			pacmanExist: true,
			gitExist:    true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS := Detect()
			if currentOS != tt.mockOS {
				t.Skipf("Skipping test: expected OS %s, but running on %s", tt.mockOS, currentOS)
			}

			// Check if we should skip this test case
			if tt.skipFunc != nil && tt.skipFunc() {
				t.Skip("Skipping test case based on skip function")
			}

			// For test cases that expect tools to exist, skip if they're missing
			// (skipFunc handles cases where we expect tools to NOT exist)
			if tt.pacmanExist && !IsPacmanInstalled() {
				t.Skip("Pacman not found on this system, skipping test")
			}

			if tt.gitExist && !commandExists("git") {
				t.Skip("Git not found on this system, skipping test")
			}

			// Ensure we're in test mode
			SetTestMode(true)
			defer SetTestMode(false)
			os.Setenv("CODORY_TEST", "1")
			defer os.Unsetenv("CODORY_TEST")

			testMode := func() bool { return true }
			installer := NewYayInstaller(testMode)
			err := installer.validateRequirements()

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

func TestYayInstaller_ensureBaseDevel(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Skip this test entirely if not on Arch Linux
	if Detect() != Arch {
		t.Skip("Skipping base-devel test on non-Arch system")
	}

	// Skip this test entirely if pacman is not available
	if !IsPacmanInstalled() {
		t.Skip("Pacman not found on this system, skipping test")
	}

	testMode := func() bool { return true }
	installer := NewYayInstaller(testMode)

	// This test only validates that the function can be called without panicking
	// We don't actually test the installation since that would modify the system
	err := installer.ensureBaseDevel(context.Background())

	// We expect the function to either succeed (if base-devel is already installed)
	// or fail gracefully (if it needs to be installed but can't in test environment)
	// The important thing is that it doesn't panic or crash
	if err != nil {
		t.Logf("ensureBaseDevel returned error (expected in test environment): %v", err)
	}
}

func TestYayInstaller_Install(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Skip this test entirely on non-Arch systems to avoid any installation attempts
	if Detect() != Arch {
		t.Skip("Skipping Install test on non-Arch system")
	}

	// Skip this test entirely if pacman is not available
	if !IsPacmanInstalled() {
		t.Skip("Pacman not found on this system, skipping test")
	}

	// Skip this test entirely if git is not available
	if !commandExists("git") {
		t.Skip("Git not found on this system, skipping test")
	}

	testMode := func() bool { return true }
	installer := NewYayInstaller(testMode)

	// This test only validates that the Install method can be called without panicking
	// We expect it to fail at some point during the actual installation process
	// since we're in a test environment, but it should fail gracefully
	err := installer.Install(context.Background())

	// We expect an error since we're in a test environment
	if err == nil {
		t.Error("Expected Install to fail in test environment, but it succeeded")
	} else {
		// Verify it's the test mode block error
		if !contains(err.Error(), "installation blocked in test mode") {
			t.Errorf("Expected 'installation blocked in test mode' error, got: %v", err)
		}
		t.Logf("Install blocked as expected in test environment: %v", err)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
