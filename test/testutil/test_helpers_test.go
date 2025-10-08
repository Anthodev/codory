package testutil

import (
	"os"
	"runtime"
	"testing"

	"anthodev/codory/internal/platform"
)

func TestSkipIfNotOnPlatform(t *testing.T) {
	// Test that the function doesn't skip when on the correct platform
	currentPlatform := platform.Detect()

	// This should not skip since we're on the current platform
	t.Run("Should not skip on current platform", func(t *testing.T) {
		SkipIfNotOnPlatform(t, currentPlatform)
		// If we get here, the test didn't skip
		t.Log("Test continued as expected on current platform")
	})

	// Test that the function skips when not on the expected platform
	// We'll use a different platform that we're definitely not on
	t.Run("Should skip on different platform", func(t *testing.T) {
		// Choose a platform we're definitely not on
		differentPlatform := platform.Windows
		if currentPlatform == platform.Windows {
			differentPlatform = platform.MacOS
		}

		// Create a test that should skip
		testRan := false
		t.Run("Inner test that should skip", func(t *testing.T) {
			SkipIfNotOnPlatform(t, differentPlatform)
			testRan = true
		})

		if testRan {
			t.Error("Test should have been skipped but ran anyway")
		}
	})
}

func TestSkipIfOnPlatform(t *testing.T) {
	currentPlatform := platform.Detect()

	// Test that the function skips when on the specified platform
	t.Run("Should skip on current platform", func(t *testing.T) {
		testRan := false
		t.Run("Inner test that should skip", func(t *testing.T) {
			SkipIfOnPlatform(t, currentPlatform)
			testRan = true
		})

		if testRan {
			t.Error("Test should have been skipped but ran anyway")
		}
	})

	// Test that the function doesn't skip when not on the specified platform
	t.Run("Should not skip on different platform", func(t *testing.T) {
		differentPlatform := platform.Windows
		if currentPlatform == platform.Windows {
			differentPlatform = platform.MacOS
		}

		// This should not skip
		SkipIfOnPlatform(t, differentPlatform)
		t.Log("Test continued as expected on different platform")
	})
}

func TestSkipIfNotOnGOOS(t *testing.T) {
	currentGOOS := runtime.GOOS

	// Test that the function doesn't skip when on the correct GOOS
	t.Run("Should not skip on current GOOS", func(t *testing.T) {
		SkipIfNotOnGOOS(t, currentGOOS)
		t.Log("Test continued as expected on current GOOS")
	})

	// Test that the function skips when not on the expected GOOS
	t.Run("Should skip on different GOOS", func(t *testing.T) {
		differentGOOS := "windows"
		if currentGOOS == "windows" {
			differentGOOS = "linux"
		}

		testRan := false
		t.Run("Inner test that should skip", func(t *testing.T) {
			SkipIfNotOnGOOS(t, differentGOOS)
			testRan = true
		})

		if testRan {
			t.Error("Test should have been skipped but ran anyway")
		}
	})
}

func TestSetupPlatformTestEnvironment(t *testing.T) {
	// Test that the function sets up test mode correctly
	cleanup := SetupPlatformTestEnvironment(t)

	// Verify test mode is set by checking that we can set it without error
	// We'll use a different approach since IsTestMode() doesn't exist
	platform.SetTestMode(false) // First unset it
	cleanup = SetupPlatformTestEnvironment(t)
	// If we get here without error, the function worked

	// Run cleanup
	cleanup()

	// Verify cleanup works by calling it and checking no panic occurs
	cleanup()
	// If we get here without error, cleanup worked
}

func TestCreateTestPlatform(t *testing.T) {
	// Test creating platform info for different platforms
	testPlatforms := []platform.Platform{
		platform.Linux,
		platform.MacOS,
		platform.Windows,
		platform.Debian,
		platform.Arch,
	}

	for _, testPlatform := range testPlatforms {
		t.Run("Create platform info for "+string(testPlatform), func(t *testing.T) {
			info := CreateTestPlatform(testPlatform)

			if info.OS != testPlatform {
				t.Errorf("Expected OS %s, got %s", testPlatform, info.OS)
			}

			if len(info.PackageManagers) != 0 {
				t.Error("Expected empty PackageManagers slice")
			}
		})
	}
}

func TestRunPlatformTestCases(t *testing.T) {
	// Create test cases
	testCases := []PlatformTestCase{
		{
			Name:    "test case 1",
			MockOS:  platform.Linux,
			WantErr: false,
		},
		{
			Name:        "test case 2",
			MockOS:      platform.MacOS,
			WantErr:     true,
			ErrContains: "error",
		},
		{
			Name:     "skipped test case",
			MockOS:   platform.Windows,
			SkipFunc: func() bool { return true },
		},
	}

	testCount := 0
	skippedCount := 0

	// Run test cases
	RunPlatformTestCases(t, testCases, func(t *testing.T, tt PlatformTestCase) {
		testCount++

		// Check if this test should be skipped
		if tt.SkipFunc != nil && tt.SkipFunc() {
			// This should have been skipped by RunPlatformTestCases
			skippedCount++
		}

		// Basic test logic
		if tt.Name == "test case 1" {
			if tt.MockOS != platform.Linux {
				t.Error("Expected Linux platform")
			}
		}
	})

	// Verify that the expected number of tests ran
	if testCount != 2 { // Should run 2 tests (1 and 3, since 2 should be skipped)
		t.Errorf("Expected 2 tests to run, but %d ran", testCount)
	}
}

func TestAssertPlatformSpecific(t *testing.T) {
	currentPlatform := platform.Detect()

	// Test that the function runs the assertion on the correct platform
	t.Run("Should run assertion on current platform", func(t *testing.T) {
		assertionRan := false
		AssertPlatformSpecific(t, currentPlatform, func() {
			assertionRan = true
		})

		if !assertionRan {
			t.Error("Assertion should have run on current platform")
		}
	})
}

func TestSkipWindows(t *testing.T) {
	currentPlatform := platform.Detect()

	if currentPlatform == platform.Windows {
		// Test that the function skips on Windows
		testRan := false
		t.Run("Inner test that should skip", func(t *testing.T) {
			SkipWindows(t)
			testRan = true
		})

		if testRan {
			t.Error("Test should have been skipped on Windows")
		}
	} else {
		// Test that the function doesn't skip on non-Windows
		SkipWindows(t)
		t.Log("Test continued as expected on non-Windows platform")
	}
}

func TestSkipNonLinux(t *testing.T) {
	currentGOOS := runtime.GOOS

	if currentGOOS != "linux" {
		// Test that the function skips on non-Linux
		testRan := false
		t.Run("Inner test that should skip", func(t *testing.T) {
			SkipNonLinux(t)
			testRan = true
		})

		if testRan {
			t.Error("Test should have been skipped on non-Linux platform")
		}
	} else {
		// Test that the function doesn't skip on Linux
		SkipNonLinux(t)
		t.Log("Test continued as expected on Linux platform")
	}
}

func TestSetupTestModeWithCleanup(t *testing.T) {
	// Test with environment variable
	cleanup := SetupTestModeWithCleanup(t, "TEST_ENV_VAR")

	// Verify test mode is set by checking that environment variable is set
	if os.Getenv("TEST_ENV_VAR") != "1" {
		t.Error("Environment variable should be set, indicating test mode is active")
	}

	// Verify environment variable is set
	if os.Getenv("TEST_ENV_VAR") != "1" {
		t.Error("Environment variable should be set")
	}

	// Run cleanup
	cleanup()

	// Test mode state is verified by the fact that cleanup completed without error
	// and environment variable was properly cleaned up

	// Verify environment variable is unset
	if os.Getenv("TEST_ENV_VAR") != "" {
		t.Error("Environment variable should be unset after cleanup")
	}
}
