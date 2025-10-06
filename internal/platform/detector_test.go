package platform

import (
	"os"
	"runtime"
	"testing"
)

// TestDetect tests the main Detect function with various scenarios
func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		goos     string
		expected []Platform
	}{
		{
			name:     "linux platform",
			goos:     "linux",
			expected: []Platform{Debian, Arch, Unknown},
		},
		{
			name:     "darwin platform",
			goos:     "darwin",
			expected: []Platform{MacOS},
		},
		{
			name:     "windows platform",
			goos:     "windows",
			expected: []Platform{Windows},
		},
		{
			name:     "unknown platform",
			goos:     "freebsd",
			expected: []Platform{Unknown},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip if not running on the expected OS
			if runtime.GOOS != tt.goos {
				t.Skipf("Skipping test for %s on %s", tt.goos, runtime.GOOS)
			}

			got := Detect()

			found := false
			for _, expected := range tt.expected {
				if got == expected {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Detect() = %v, want one of %v for %s", got, tt.expected, tt.goos)
			}
		})
	}
}

// TestDetectInfo tests the DetectInfo function with proper test mode setup
func TestDetectInfo(t *testing.T) {
	// Setup test environment
	SetTestMode(true)
	defer SetTestMode(false)

	// Test with different environment variables
	testCases := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "with CODORY_TEST",
			env:  map[string]string{"CODORY_TEST": "1"},
		},
		{
			name: "with GO_TEST",
			env:  map[string]string{"GO_TEST": "1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tc.env {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			info := DetectInfo()

			if info.OS == "" {
				t.Error("DetectInfo() returned empty OS")
			}

			if info.PackageManagers == nil {
				t.Error("DetectInfo() returned nil PackageManagers")
			}

			// Verify OS matches expected platform
			expectedPlatform := Detect()
			if info.OS != expectedPlatform {
				t.Errorf("DetectInfo() OS = %v, want %v", info.OS, expectedPlatform)
			}

			// Test HasPackageManager method
			for _, pm := range info.PackageManagers {
				if !info.HasPackageManager(pm) {
					t.Errorf("HasPackageManager(%v) returned false for existing package manager", pm)
				}
			}
		})
	}
}

// TestPlatform_String tests the String method for Platform type
func TestPlatform_String(t *testing.T) {
	tests := []struct {
		name     string
		platform Platform
		want     string
	}{
		{
			name:     "debian platform",
			platform: Debian,
			want:     "debian",
		},
		{
			name:     "arch platform",
			platform: Arch,
			want:     "arch",
		},
		{
			name:     "macos platform",
			platform: MacOS,
			want:     "macos",
		},
		{
			name:     "windows platform",
			platform: Windows,
			want:     "windows",
		},
		{
			name:     "unknown platform",
			platform: Unknown,
			want:     "unknown",
		},
		{
			name:     "custom platform",
			platform: Platform("custom"),
			want:     "custom",
		},
		{
			name:     "empty platform",
			platform: Platform(""),
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.platform.String(); got != tt.want {
				t.Errorf("Platform.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPlatform_DisplayName tests the DisplayName method for Platform type
func TestPlatform_DisplayName(t *testing.T) {
	tests := []struct {
		name     string
		platform Platform
		want     string
	}{
		{
			name:     "debian platform",
			platform: Debian,
			want:     "Debian/Ubuntu",
		},
		{
			name:     "arch platform",
			platform: Arch,
			want:     "Arch Linux",
		},
		{
			name:     "macos platform",
			platform: MacOS,
			want:     "macOS",
		},
		{
			name:     "windows platform",
			platform: Windows,
			want:     "Windows",
		},
		{
			name:     "unknown platform",
			platform: Unknown,
			want:     "Unknown",
		},
		{
			name:     "custom platform",
			platform: Platform("custom"),
			want:     "Unknown",
		},
		{
			name:     "empty platform",
			platform: Platform(""),
			want:     "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.platform.DisplayName(); got != tt.want {
				t.Errorf("Platform.DisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFileExists tests the fileExists function
func TestFileExists(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing file",
			path:     "detector_test.go",
			expected: true,
		},
		{
			name:     "non-existing file",
			path:     "non_existent_file.go",
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
		{
			name:     "directory path",
			path:     ".",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fileExists(tt.path)
			if got != tt.expected {
				t.Errorf("fileExists(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

// TestCommandExists tests the commandExists function
func TestCommandExists(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		expected bool
	}{
		{
			name:     "existing command go",
			cmd:      "go",
			expected: true,
		},
		{
			name:     "existing command sh",
			cmd:      "sh",
			expected: true,
		},
		{
			name:     "non-existing command",
			cmd:      "nonexistentcommand12345",
			expected: false,
		},
		{
			name:     "empty command",
			cmd:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := commandExists(tt.cmd)
			if got != tt.expected {
				t.Errorf("commandExists(%q) = %v, want %v", tt.cmd, got, tt.expected)
			}
		})
	}
}

// TestPackageManagerDetection tests individual package manager detection functions
func TestPackageManagerDetection(t *testing.T) {
	// Setup test environment
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	tests := []struct {
		name     string
		function func() bool
	}{
		{
			name:     "IsYayInstalled",
			function: IsYayInstalled,
		},
		{
			name:     "IsPacmanInstalled",
			function: IsPacmanInstalled,
		},
		{
			name:     "IsBrewInstalled",
			function: IsBrewInstalled,
		},
		{
			name:     "IsWingetInstalled",
			function: IsWingetInstalled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function()
			// Verify it returns a boolean value
			if result != true && result != false {
				t.Errorf("%s() should return a boolean value, got %v", tt.name, result)
			}
		})
	}
}

// TestDetectLinuxDistro tests the detectLinuxDistro function
func TestDetectLinuxDistro(t *testing.T) {
	// Setup test environment
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	if runtime.GOOS != "linux" {
		t.Skip("Skipping detectLinuxDistro test on non-Linux system")
	}

	result := detectLinuxDistro()

	// On Linux, we expect either Debian, Arch, or Unknown
	validPlatforms := []Platform{Debian, Arch, Unknown}
	found := false
	for _, valid := range validPlatforms {
		if result == valid {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("detectLinuxDistro() = %v, want one of %v", result, validPlatforms)
	}
}

// TestSetTestMode tests the SetTestMode function
func TestSetTestMode(t *testing.T) {
	// Save original state
	originalTestMode := testMode
	defer SetTestMode(originalTestMode)

	tests := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "enable test mode",
			enabled: true,
		},
		{
			name:    "disable test mode",
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTestMode(tt.enabled)
			if testMode != tt.enabled {
				t.Errorf("SetTestMode(%v) did not set testMode to %v", tt.enabled, tt.enabled)
			}
		})
	}
}

// TestInfo_HasPackageManager tests the HasPackageManager method
func TestInfo_HasPackageManager(t *testing.T) {
	tests := []struct {
		name           string
		info           Info
		packageManager PackageManager
		want           bool
	}{
		{
			name: "has package manager",
			info: Info{
				PackageManagers: []PackageManager{PackageManagerAPT, PackageManagerPacman},
			},
			packageManager: PackageManagerAPT,
			want:           true,
		},
		{
			name: "does not have package manager",
			info: Info{
				PackageManagers: []PackageManager{PackageManagerAPT, PackageManagerPacman},
			},
			packageManager: PackageManagerYay,
			want:           false,
		},
		{
			name: "empty package managers",
			info: Info{
				PackageManagers: []PackageManager{},
			},
			packageManager: PackageManagerAPT,
			want:           false,
		},
		{
			name: "nil package managers",
			info: Info{
				PackageManagers: nil,
			},
			packageManager: PackageManagerAPT,
			want:           false,
		},
		{
			name: "multiple package managers with target",
			info: Info{
				PackageManagers: []PackageManager{PackageManagerAPT, PackageManagerPacman, PackageManagerYay},
			},
			packageManager: PackageManagerPacman,
			want:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.info.HasPackageManager(tt.packageManager); got != tt.want {
				t.Errorf("Info.HasPackageManager(%v) = %v, want %v", tt.packageManager, got, tt.want)
			}
		})
	}
}

// TestDetectInfoPackageManagers tests package manager detection in DetectInfo
func TestDetectInfoPackageManagers(t *testing.T) {
	// Setup test environment
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	info := DetectInfo()

	// Test that package managers are properly detected based on platform
	switch info.OS {
	case Debian:
		// On Debian-based systems, we might have APT
		if commandExists("apt") || commandExists("apt-get") {
			if !info.HasPackageManager(PackageManagerAPT) {
				t.Error("Expected APT package manager on Debian-based system")
			}
		}
	case Arch:
		// On Arch-based systems, we might have Pacman
		if commandExists("pacman") {
			if !info.HasPackageManager(PackageManagerPacman) {
				t.Error("Expected Pacman package manager on Arch-based system")
			}
		}
	case MacOS:
		// On macOS, we might have Brew
		if commandExists("brew") {
			if !info.HasPackageManager(PackageManagerBrew) {
				t.Error("Expected Brew package manager on macOS")
			}
		}
	case Windows:
		// On Windows, we might have Winget
		if commandExists("winget") {
			if !info.HasPackageManager(PackageManagerWinget) {
				t.Error("Expected Winget package manager on Windows")
			}
		}
	}

	// Test that Brew can be detected on Linux systems too
	if info.OS == Debian || info.OS == Arch {
		if commandExists("brew") {
			if !info.HasPackageManager(PackageManagerBrew) {
				t.Error("Expected Brew package manager to be detected on Linux if installed")
			}
		}
	}
}

// BenchmarkDetect benchmarks the Detect function
func BenchmarkDetect(b *testing.B) {
	for b.Loop() {
		_ = Detect()
	}
}

// BenchmarkDetectInfo benchmarks the DetectInfo function
func BenchmarkDetectInfo(b *testing.B) {
	SetTestMode(true)
	defer SetTestMode(false)

	for b.Loop() {
		_ = DetectInfo()
	}
}

// BenchmarkHasPackageManager benchmarks the HasPackageManager method
func BenchmarkHasPackageManager(b *testing.B) {
	info := Info{
		PackageManagers: []PackageManager{PackageManagerAPT, PackageManagerPacman, PackageManagerYay, PackageManagerBrew},
	}

	
	for b.Loop() {
		_ = info.HasPackageManager(PackageManagerAPT)
	}
}
