package platform

import (
	"os/exec"
	"runtime"
)

// Platform represents a system platform
type Platform string

const (
	Unknown Platform = "unknown"
)

// PackageManager represents a package manager
type PackageManager string

const (
	PackageManagerNone PackageManager = "none"
)

// Info contains platform information
type Info struct {
	OS              Platform
	PackageManagers []PackageManager
}

// Detect detects the current platform
func Detect() Platform {
	switch runtime.GOOS {
	default:
		return Unknown
	}
}

// DetectInfo detects the current platform and available package managers
func DetectInfo() Info {
	info := Info{
		OS:              Detect(),
		PackageManagers: make([]PackageManager, 0),
	}

	return info
}

func fileExists(path string) bool {
	cmd := exec.Command("test", "-f", path)
	return cmd.Run() == nil
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func (p Platform) String() string {
	return string(p)
}

func (p Platform) DisplayName() string {
	switch p {
	default:
		return "Unknown"
	}
}
