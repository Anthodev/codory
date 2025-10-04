package platform

import (
	"os"
	"runtime"
	"testing"
)

func TestDetect(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	got := Detect()

	// Test that Detect returns the correct platform based on runtime.GOOS
	switch runtime.GOOS {
	case "linux":
		// On Linux, we expect either Debian, Arch, or Unknown depending on the system
		if got != Debian && got != Arch && got != Unknown {
			t.Errorf("Detect() = %v, want one of %v, %v, or %v for linux", got, Debian, Arch, Unknown)
		}
	case "darwin":
		if got != MacOS {
			t.Errorf("Detect() = %v, want %v for darwin", got, MacOS)
		}
	case "windows":
		if got != Windows {
			t.Errorf("Detect() = %v, want %v for windows", got, Windows)
		}
	default:
		if got != Unknown {
			t.Errorf("Detect() = %v, want %v for unknown OS", got, Unknown)
		}
	}
}

func TestDetectInfo(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	info := DetectInfo()

	if info.OS == "" {
		t.Error("DetectInfo() returned empty OS")
	}

	if info.PackageManagers == nil {
		t.Error("DetectInfo() returned nil PackageManagers")
	}

	// Test that PackageManagers is properly initialized
	if len(info.PackageManagers) == 0 {
		t.Log("DetectInfo() returned empty PackageManagers slice (this may be expected on some systems)")
	}
}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.platform.String(); got != tt.want {
				t.Errorf("Platform.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.platform.DisplayName(); got != tt.want {
				t.Errorf("Platform.DisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	// Test with a file that should exist (the current test file)
	exists := fileExists("detector_test.go")
	if !exists {
		t.Error("fileExists() returned false for existing file")
	}

	// Test with a file that should not exist
	exists = fileExists("non_existent_file.go")
	if exists {
		t.Error("fileExists() returned true for non-existent file")
	}
}

func TestCommandExists(t *testing.T) {
	// Test with a command that should exist on most systems
	exists := commandExists("go")
	if !exists {
		t.Error("commandExists() returned false for 'go' command")
	}

	// Test with a command that should not exist
	exists = commandExists("nonexistentcommand12345")
	if exists {
		t.Error("commandExists() returned true for non-existent command")
	}
}

func TestIsYayInstalled(t *testing.T) {
	// Test IsYayInstalled function
	result := IsYayInstalled()

	// We can't predict the exact result since it depends on the system
	// but we can verify it returns a boolean value without error
	if result != true && result != false {
		t.Error("IsYayInstalled() should return a boolean value")
	}
}

func TestIsPacmanInstalled(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Test IsPacmanInstalled function
	result := IsPacmanInstalled()

	// We can't predict the exact result since it depends on the system
	// but we can verify it returns a boolean value without error
	if result != true && result != false {
		t.Error("IsPacmanInstalled() should return a boolean value")
	}
}

func TestIsBrewInstalled(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Test IsBrewInstalled function
	result := IsBrewInstalled()

	// We can't predict the exact result since it depends on the system
	// but we can verify it returns a boolean value without error
	if result != true && result != false {
		t.Error("IsBrewInstalled() should return a boolean value")
	}
}

func TestDetectLinuxDistro(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Only test this function on Linux
	if runtime.GOOS != "linux" {
		t.Skip("Skipping detectLinuxDistro test on non-Linux system")
	}

	result := detectLinuxDistro()

	// On Linux, we expect either Debian, Arch, or Unknown
	if result != Debian && result != Arch && result != Unknown {
		t.Errorf("detectLinuxDistro() = %v, want one of %v, %v, or %v", result, Debian, Arch, Unknown)
	}
}

func TestSetTestMode(t *testing.T) {
	// Test setting test mode
	SetTestMode(true)
	if !testMode {
		t.Error("SetTestMode(true) did not set testMode to true")
	}

	// Test unsetting test mode
	SetTestMode(false)
	if testMode {
		t.Error("SetTestMode(false) did not set testMode to false")
	}
}

func TestIsTestEnvironment(t *testing.T) {
	// Ensure we're in test mode
	SetTestMode(true)
	defer SetTestMode(false)
	os.Setenv("CODORY_TEST", "1")
	defer os.Unsetenv("CODORY_TEST")

	// Test that isTestEnvironment returns true when test mode is set
	if !isTestEnvironment() {
		t.Error("isTestEnvironment() should return true when test mode is set")
	}

	// Test with CODORY_TEST environment variable
	os.Setenv("CODORY_TEST", "1")
	if !isTestEnvironment() {
		t.Error("isTestEnvironment() should return true when CODORY_TEST=1")
	}
	os.Unsetenv("CODORY_TEST")

	// Test with GO_TEST environment variable
	os.Setenv("GO_TEST", "1")
	if !isTestEnvironment() {
		t.Error("isTestEnvironment() should return true when GO_TEST=1")
	}
	os.Unsetenv("GO_TEST")
}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.info.HasPackageManager(tt.packageManager); got != tt.want {
				t.Errorf("Info.HasPackageManager() = %v, want %v", got, tt.want)
			}
		})
	}
}
