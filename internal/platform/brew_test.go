package platform

import (
	"context"
	"os"
	"runtime"
	"testing"
)

func TestNewBrewInstaller(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	testMode := func() bool { return true }
	installer := NewBrewInstaller(testMode)

	if installer == nil {
		t.Fatal("NewBrewInstaller() returned nil")
	}

	if installer.isTestMode == nil {
		t.Error("Expected isTestMode function to be set")
	}
}

func TestBrewInstaller_ValidateRequirements(t *testing.T) {
	tests := []struct {
		name        string
		mockOS      Platform
		brewExist   bool
		bashExist   bool
		curlExist   bool
		gitExist    bool
		wantErr     bool
		errContains string
		skipFunc    func() bool
	}{
		{
			name:        "brew already installed",
			mockOS:      MacOS,
			brewExist:   true,
			bashExist:   true,
			curlExist:   true,
			gitExist:    true,
			wantErr:     true,
			errContains: "brew is already installed",
		},
		{
			name:        "install on Windows platform",
			mockOS:      Windows,
			brewExist:   false,
			bashExist:   true,
			curlExist:   true,
			gitExist:    true,
			wantErr:     true,
			errContains: "brew cannot be installed on Windows",
		},
		{
			name:        "install on supported platform without bash",
			mockOS:      MacOS,
			brewExist:   false,
			bashExist:   false,
			curlExist:   true,
			gitExist:    true,
			wantErr:     true,
			errContains: "bash is required to install brew",
			skipFunc: func() bool {
				// Skip this test case if bash is actually installed on the system
				return commandExists("bash")
			},
		},
		{
			name:        "install on supported platform without curl",
			mockOS:      MacOS,
			brewExist:   false,
			bashExist:   true,
			curlExist:   false,
			gitExist:    true,
			wantErr:     true,
			errContains: "curl is required to install brew",
			skipFunc: func() bool {
				// Skip this test case if curl is actually installed on the system
				return commandExists("curl")
			},
		},
		{
			name:        "install on supported platform without git",
			mockOS:      MacOS,
			brewExist:   false,
			bashExist:   true,
			curlExist:   true,
			gitExist:    false,
			wantErr:     true,
			errContains: "git is required to install brew",
			skipFunc: func() bool {
				// Skip this test case if git is actually installed on the system
				return commandExists("git")
			},
		},
		{
			name:      "valid setup for installation",
			mockOS:    MacOS,
			brewExist: false,
			bashExist: true,
			curlExist: true,
			gitExist:  true,
			wantErr:   false,
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
			if tt.bashExist && !commandExists("bash") {
				t.Skip("Bash not found on this system, skipping test")
			}

			if tt.curlExist && !commandExists("curl") {
				t.Skip("Curl not found on this system, skipping test")
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
			installer := NewBrewInstaller(testMode)
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

func TestBrewInstaller_addBrewToPath(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Skip this test entirely if not on Linux
	if runtime.GOOS != "linux" {
		t.Skip("Skipping addBrewToPath test on non-Linux system")
	}

	// Skip this test entirely if bash is not available
	if !commandExists("bash") {
		t.Skip("Bash not found on this system, skipping test")
	}

	testMode := func() bool { return true }
	installer := NewBrewInstaller(testMode)

	// This test only validates that the function can be called without panicking
	// We don't actually test the path modification since that would modify the system
	err := installer.addBrewToPath(context.Background())

	// We expect the function to either succeed or fail gracefully
	// The important thing is that it doesn't panic or crash
	if err != nil {
		t.Logf("addBrewToPath returned error (expected in test environment): %v", err)
	}
}

func TestBrewInstaller_Install(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Skip this test entirely on Windows to avoid any installation attempts
	if Detect() == Windows {
		t.Skip("Skipping Install test on Windows system")
	}

	// Skip this test entirely if bash is not available
	if !commandExists("bash") {
		t.Skip("Bash not found on this system, skipping test")
	}

	// Skip this test entirely if curl is not available
	if !commandExists("curl") {
		t.Skip("Curl not found on this system, skipping test")
	}

	// Skip this test entirely if git is not available
	if !commandExists("git") {
		t.Skip("Git not found on this system, skipping test")
	}

	// If brew is already installed, the validation will fail with "brew is already installed"
	// before we get to the test mode check. This is expected behavior.
	if IsBrewInstalled() {
		t.Skip("Skipping Install test because brew is already installed")
	}

	testMode := func() bool { return true }
	installer := NewBrewInstaller(testMode)

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

func TestGetBrewPath(t *testing.T) {
	// Test that GetBrewPath returns the expected path based on runtime.GOOS and runtime.GOARCH
	result := GetBrewPath()

	// We can't mock runtime.GOOS and runtime.GOARCH directly, so we test the actual behavior
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			if result != "/opt/homebrew/bin/brew" {
				t.Errorf("GetBrewPath() = %v, want /opt/homebrew/bin/brew for darwin arm64", result)
			}
		} else {
			if result != "/usr/local/bin/brew" {
				t.Errorf("GetBrewPath() = %v, want /usr/local/bin/brew for darwin amd64", result)
			}
		}
	case "linux":
		if result != "/home/linuxbrew/.linuxbrew/bin/brew" {
			t.Errorf("GetBrewPath() = %v, want /home/linuxbrew/.linuxbrew/bin/brew for linux", result)
		}
	default:
		// For other platforms, we expect the Linux path as default
		if result != "/home/linuxbrew/.linuxbrew/bin/brew" {
			t.Errorf("GetBrewPath() = %v, want /home/linuxbrew/.linuxbrew/bin/brew for default", result)
		}
	}
}
