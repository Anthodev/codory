package platform

import (
	"os/exec"
	"runtime"
	"slices"
)

// Platform represents a system platform
type Platform string

const (
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
	// Vérifier si c'est une distribution basée sur Debian
	if fileExists("/etc/debian_version") {
		return Debian
	}

	// Vérifier si c'est Arch Linux
	if fileExists("/etc/arch-release") {
		return Arch
	}

	// Vérifier via les gestionnaires de paquets
	if commandExists("apt") || commandExists("apt-get") {
		return Debian
	}

	if commandExists("pacman") {
		return Arch
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

	// Brew peut aussi être disponible sur Linux
	if info.OS == Debian || info.OS == Arch {
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
	case Debian:
		return "Debian/Ubuntu"
	case Arch:
		return "Arch Linux"
	case MacOS:
		return "macOS"
	case Windows:
		return "Windows"
	default:
		return "Unknown"
	}
}
