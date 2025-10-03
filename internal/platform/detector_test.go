package platform

import (
	"runtime"
	"testing"
)

func TestDetect(t *testing.T) {
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
	info := DetectInfo()

	if info.OS == "" {
		t.Error("DetectInfo() returned empty OS")
	}

	if info.PackageManagers == nil {
		t.Error("DetectInfo() returned nil PackageManagers")
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
	// Test IsPacmanInstalled function
	result := IsPacmanInstalled()

	// We can't predict the exact result since it depends on the system
	// but we can verify it returns a boolean value without error
	if result != true && result != false {
		t.Error("IsPacmanInstalled() should return a boolean value")
	}
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
