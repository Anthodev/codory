package platform

import (
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
)

var testMode = false

func SetTestMode(enabled bool) {
	testMode = enabled
}

// Platform represents a system platform
type Platform string

const (
	Linux   Platform = "linux"
	Debian  Platform = "debian"
	Arch    Platform = "arch"
	MacOS   Platform = "macos"
	Windows Platform = "windows"
	Unknown Platform = "unknown"
)

// PackageManager represents a package manager
type PackageManager string

const (
	PackageManagerAPT    PackageManager = "apt"
	PackageManagerPacman PackageManager = "pacman"
	PackageManagerYay    PackageManager = "yay"
	PackageManagerBrew   PackageManager = "brew"
	PackageManagerWinget PackageManager = "winget"
	PackageManagerNone   PackageManager = "none"
)

// Info contains platform information
type Info struct {
	OS              Platform
	PackageManagers []PackageManager
}

// Detect detects the current platform
func Detect() Platform {
	switch runtime.GOOS {
	case "darwin":
		return MacOS
	case "windows":
		return Windows
	case "linux":
		return detectLinuxDistro()
	default:
		return Unknown
	}
}

func detectLinuxDistro() Platform {
	if fileExists("/etc/debian_version") {
		return Debian
	}

	if fileExists("/etc/arch-release") {
		return Arch
	}

	if commandExists("apt") || commandExists("apt-get") {
		return Debian
	}

	if commandExists("pacman") {
		return Arch
	}

	if commandExists("uname") {
		return Linux
	}

	return Unknown
}

// DetectInfo detects the current platform and available package managers
func DetectInfo() Info {
	info := Info{
		OS:              Detect(),
		PackageManagers: make([]PackageManager, 0),
	}

	// Détecter les gestionnaires de paquets disponibles
	switch info.OS {
	case Debian:
		if commandExists("apt") || commandExists("apt-get") {
			info.PackageManagers = append(info.PackageManagers, PackageManagerAPT)
		}

	case Arch:
		if commandExists("pacman") {
			info.PackageManagers = append(info.PackageManagers, PackageManagerPacman)
		}
		if commandExists("yay") {
			info.PackageManagers = append(info.PackageManagers, PackageManagerYay)
		}

	case MacOS:
		if commandExists("brew") {
			info.PackageManagers = append(info.PackageManagers, PackageManagerBrew)
		}

	case Windows:
		if commandExists("winget") {
			info.PackageManagers = append(info.PackageManagers, PackageManagerWinget)
		}
	}

	if info.OS == Linux || info.OS == Debian || info.OS == Arch {
		if commandExists("brew") {
			info.PackageManagers = append(info.PackageManagers, PackageManagerBrew)
		}
	}

	return info
}

func (i Info) HasPackageManager(pm PackageManager) bool {
	return slices.Contains(i.PackageManagers, pm)
}

func IsYayInstalled() bool {
	return commandExists("yay")
}

func IsPacmanInstalled() bool {
	return commandExists("pacman")
}

func IsBrewInstalled() bool {
	return commandExists("brew")
}

// IsWingetInstalled is a variable holding the function to check if winget is installed.
// This allows for easier testing by mocking this function.
var IsWingetInstalled = func() bool {
	return commandExists("winget")
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
	case Linux:
		return "Linux"
	case Debian:
		return "Debian/Ubuntu and derivatives"
	case Arch:
		return "Arch Linux and derivatives"
	case MacOS:
		return "macOS"
	case Windows:
		return "Windows"
	default:
		return "Unknown"
	}
}

func isTestEnvironment() bool {
	// 1. Check package-level test mode flag (set by tests)
	if testMode {
		return true
	}

	// 2. Check CODORY_TEST environment variable
	if os.Getenv("CODORY_TEST") == "1" {
		return true
	}

	// 3. Check GO_TEST environment variable
	if os.Getenv("GO_TEST") == "1" {
		return true
	}

	// 4. Check if test binary is running (test binaries have .test suffix or contain .test. in name)
	if exePath, err := os.Executable(); err == nil {
		if len(exePath) > 0 {
			// Check if it's a test binary
			return strings.Contains(exePath, ".test") || strings.Contains(exePath, "_test")
		}
	}

	return false
}
