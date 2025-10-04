package package_managers

import (
	"anthodev/codory/internal/actions"
	"os"
	"strings"
)

var testMode = false

func SetTestMode(enabled bool) {
	testMode = enabled
}

func init() {
	registry := actions.GlobalRegistry()

	packageManagersCategory := &actions.Category{
		ID:          "package_managers",
		Name:        "Package Managers",
		Description: "Tools and utilities specific to Linux",
		Actions:     make([]*actions.Action, 0),
	}

	registry.RegisterCategory("root", packageManagersCategory)

	registry.RegisterAction("package_managers", NewInstallYayAction())
	registry.RegisterAction("package_managers", NewInstallBrewAction())
}

// isTestEnvironment checks if we're running in a test environment
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
