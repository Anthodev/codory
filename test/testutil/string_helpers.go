// Package testutil provides common testing utilities and helpers for string operations and platform testing
package testutil

import (
	"os"
	"runtime"
	"testing"

	"anthodev/codory/internal/platform"
	"anthodev/codory/pkg/utils"
)

// AssertContains checks if a string contains a substring and reports an error if not found
func AssertContains(t *testing.T, s, substr string, message string) {
	t.Helper()
	if !utils.ContainsString(s, substr) {
		if message != "" {
			t.Errorf("%s: string '%s' does not contain '%s'", message, s, substr)
		} else {
			t.Errorf("string '%s' does not contain '%s'", s, substr)
		}
	}
}

// AssertContainsSubstring checks if a string contains a substring using the efficient implementation
func AssertContainsSubstring(t *testing.T, s, substr string, message string) {
	t.Helper()
	if !utils.ContainsSubstring(s, substr) {
		if message != "" {
			t.Errorf("%s: string '%s' does not contain '%s'", message, s, substr)
		} else {
			t.Errorf("string '%s' does not contain '%s'", s, substr)
		}
	}
}

// AssertNotContains checks if a string does not contain a substring
func AssertNotContains(t *testing.T, s, substr string, message string) {
	t.Helper()
	if utils.ContainsString(s, substr) {
		if message != "" {
			t.Errorf("%s: string '%s' unexpectedly contains '%s'", message, s, substr)
		} else {
			t.Errorf("string '%s' unexpectedly contains '%s'", s, substr)
		}
	}
}

// AssertStringEquals checks if two strings are equal
func AssertStringEquals(t *testing.T, got, want string, message string) {
	t.Helper()
	if got != want {
		if message != "" {
			t.Errorf("%s: got '%s', want '%s'", message, got, want)
		} else {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	}
}

// AssertError checks if an error is not nil
func AssertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		if message != "" {
			t.Errorf("%s: expected error, got nil", message)
		} else {
			t.Error("expected error, got nil")
		}
	}
}

// AssertNoError checks if an error is nil
func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		if message != "" {
			t.Errorf("%s: unexpected error: %v", message, err)
		} else {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

// PlatformTestCase represents a test case that varies by platform
type PlatformTestCase struct {
	Name        string
	MockOS      platform.Platform
	WantErr     bool
	ErrContains string
	SkipFunc    func() bool
	SetupFunc   func() // Optional setup function
	CleanupFunc func() // Optional cleanup function
}

// SkipIfNotOnPlatform skips the test if not running on the expected platform
func SkipIfNotOnPlatform(t *testing.T, expectedPlatform platform.Platform) {
	t.Helper()
	currentOS := platform.Detect()
	if currentOS != expectedPlatform {
		t.Skipf("Skipping test: expected OS %s, but running on %s", expectedPlatform, currentOS)
	}
}

// SkipIfOnPlatform skips the test if running on the specified platform
func SkipIfOnPlatform(t *testing.T, platformToSkip platform.Platform) {
	t.Helper()
	currentOS := platform.Detect()
	if currentOS == platformToSkip {
		t.Skipf("Skipping test on %s platform", platformToSkip)
	}
}

// SkipIfNotOnGOOS skips the test if not running on the expected GOOS
func SkipIfNotOnGOOS(t *testing.T, expectedGOOS string) {
	t.Helper()
	if runtime.GOOS != expectedGOOS {
		t.Skipf("Skipping test: expected GOOS %s, but running on %s", expectedGOOS, runtime.GOOS)
	}
}

// SetupPlatformTestEnvironment sets up a consistent test environment for platform-dependent tests
func SetupPlatformTestEnvironment(t *testing.T) (cleanup func()) {
	t.Helper()

	// Set test mode
	platform.SetTestMode(true)

	// Return cleanup function
	return func() {
		platform.SetTestMode(false)
	}
}

// CreateTestPlatform creates platform info for testing with specified OS
func CreateTestPlatform(os platform.Platform) platform.Info {
	return platform.Info{
		OS:              os,
		PackageManagers: []platform.PackageManager{},
	}
}

// RunPlatformTestCases runs a set of platform-based test cases
func RunPlatformTestCases(t *testing.T, testCases []PlatformTestCase, testFunc func(t *testing.T, tt PlatformTestCase)) {
	t.Helper()

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			// Check skip function
			if tt.SkipFunc != nil && tt.SkipFunc() {
				t.Skip("Skipping test case based on skip function")
			}

			// Setup if needed
			if tt.SetupFunc != nil {
				tt.SetupFunc()
			}

			// Cleanup if needed
			if tt.CleanupFunc != nil {
				defer tt.CleanupFunc()
			}

			// Run the actual test
			testFunc(t, tt)
		})
	}
}

// AssertPlatformSpecific skips test if not on expected platform and runs assertion
func AssertPlatformSpecific(t *testing.T, expectedPlatform platform.Platform, assertion func()) {
	t.Helper()
	SkipIfNotOnPlatform(t, expectedPlatform)
	assertion()
}

// SkipWindows skips test on Windows platform
func SkipWindows(t *testing.T) {
	t.Helper()
	SkipIfOnPlatform(t, platform.Windows)
}

// SkipNonLinux skips test on non-Linux platforms
func SkipNonLinux(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skipf("Skipping test on non-Linux platform: %s", runtime.GOOS)
	}
}

// SetupTestModeWithCleanup sets up test mode and returns cleanup function
func SetupTestModeWithCleanup(t *testing.T, testEnvVar string) (cleanup func()) {
	t.Helper()

	// Set test mode
	platform.SetTestMode(true)
	if testEnvVar != "" {
		os.Setenv(testEnvVar, "1")
	}

	// Return cleanup function
	return func() {
		platform.SetTestMode(false)
		if testEnvVar != "" {
			os.Unsetenv(testEnvVar)
		}
	}
}
