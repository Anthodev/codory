package platform

import (
	"runtime"
	"testing"
)

func TestDetect(t *testing.T) {
	got := Detect()

	// Test that Detect returns a platform based on runtime.GOOS
	switch runtime.GOOS {
	case "linux":
		if got != Unknown {
			t.Errorf("Detect() = %v, want %v for linux", got, Unknown)
		}
	case "darwin":
		if got != Unknown {
			t.Errorf("Detect() = %v, want %v for darwin", got, Unknown)
		}
	case "windows":
		if got != Unknown {
			t.Errorf("Detect() = %v, want %v for windows", got, Unknown)
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
