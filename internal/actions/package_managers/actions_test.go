package package_managers

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestInit(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	packageManagersCategory, exists := globalRegistry.GetCategory("package_managers")
	if !exists {
		t.Fatal("Package Managers category not found in global registry")
	}

	// Verify package_managers category properties
	if packageManagersCategory.ID != "package_managers" {
		t.Errorf("Expected package_managers category ID to be 'package_managers', got '%s'", packageManagersCategory.ID)
	}

	if packageManagersCategory.Name != "Package Managers" {
		t.Errorf("Expected package_managers category name to be 'Package Managers', got '%s'", packageManagersCategory.Name)
	}

	if packageManagersCategory.Description != "Tools and utilities specific to Linux" {
		t.Errorf("Expected package_managers category description to be 'Tools and utilities specific to Linux', got '%s'", packageManagersCategory.Description)
	}

	// Check if both actions are registered
	if len(packageManagersCategory.Actions) != 2 {
		t.Fatalf("Expected 2 actions in package_managers category, got %d", len(packageManagersCategory.Actions))
	}

	var installYayAction *actions.Action
	var installBrewAction *actions.Action
	for _, action := range packageManagersCategory.Actions {
		if action.ID == "install_yay" {
			installYayAction = action
		} else if action.ID == "install_brew" {
			installBrewAction = action
		}
	}

	if installYayAction == nil {
		t.Fatal("InstallYay action not found in package_managers category")
	}

	if installBrewAction == nil {
		t.Fatal("InstallBrew action not found in package_managers category")
	}

	// Test InstallYay action properties
	if installYayAction.Name != "Install Yay (AUR Helper)" {
		t.Errorf("Expected InstallYay action name to be 'Install Yay (AUR Helper)', got '%s'", installYayAction.Name)
	}

	if installYayAction.Description != "Install Yay AUR helper for Arch Linux" {
		t.Errorf("Expected InstallYay action description to be 'Install Yay AUR helper for Arch Linux', got '%s'", installYayAction.Description)
	}

	if installYayAction.Type != actions.ActionTypeFunction {
		t.Errorf("Expected InstallYay action type to be 'function', got '%s'", installYayAction.Type)
	}

	if installYayAction.Handler == nil {
		t.Error("InstallYay action has no handler")
	}

	if len(installYayAction.VisibleOnPlatforms) != 1 {
		t.Fatalf("Expected 1 visible platform for InstallYay, got %d", len(installYayAction.VisibleOnPlatforms))
	}

	if installYayAction.VisibleOnPlatforms[0] != actions.PlatformArch {
		t.Errorf("Expected InstallYay action to be visible on 'arch' platform, got '%s'", installYayAction.VisibleOnPlatforms[0])
	}

	// Test InstallBrew action properties
	if installBrewAction.Name != "Install Homebrew" {
		t.Errorf("Expected InstallBrew action name to be 'Install Homebrew', got '%s'", installBrewAction.Name)
	}

	if installBrewAction.Description != "Install Homebrew on the system" {
		t.Errorf("Expected InstallBrew action description to be 'Install Homebrew on the system', got '%s'", installBrewAction.Description)
	}

	if installBrewAction.Type != actions.ActionTypeFunction {
		t.Errorf("Expected InstallBrew action type to be 'function', got '%s'", installBrewAction.Type)
	}

	if installBrewAction.Handler == nil {
		t.Error("InstallBrew action has no handler")
	}

	if len(installBrewAction.HiddenOnPlatforms) != 1 {
		t.Fatalf("Expected 1 hidden platform for InstallBrew, got %d", len(installBrewAction.HiddenOnPlatforms))
	}

	if installBrewAction.HiddenOnPlatforms[0] != actions.PlatformWindows {
		t.Errorf("Expected InstallBrew action to be hidden on 'windows' platform, got '%s'", installBrewAction.HiddenOnPlatforms[0])
	}
}

func TestPackageManagersCategoryIsChildOfRoot(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the root category
	rootCategory := globalRegistry.GetRoot()
	if rootCategory == nil {
		t.Fatal("Root category not found")
	}

	// Check if package_managers category is a subcategory of root
	packageManagersCategoryFound := false
	for _, subCat := range rootCategory.SubCategories {
		if subCat.ID == "package_managers" {
			packageManagersCategoryFound = true
			break
		}
	}

	if !packageManagersCategoryFound {
		t.Error("Package Managers category is not a subcategory of root")
	}
}

func TestPackageManagersCategoryActionsAreProperlyRegistered(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	packageManagersCategory, exists := globalRegistry.GetCategory("package_managers")
	if !exists {
		t.Fatal("Package Managers category not found")
	}

	for _, action := range packageManagersCategory.Actions {
		if action.ID == "" {
			t.Error("Action ID is empty")
		}
		if action.Name == "" {
			t.Errorf("Action name is empty for ID: %s", action.ID)
		}
		if action.Type == "" {
			t.Errorf("Action type is empty for ID: %s", action.ID)
		}
		if action.Type == actions.ActionTypeFunction && action.Handler == nil {
			t.Errorf("Function action has no handler for ID: %s", action.ID)
		}
	}
}
