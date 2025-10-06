package package_managers

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestInit(t *testing.T) {
	// Set test mode to prevent actual system operations
	SetTestMode(true)
	defer SetTestMode(false)

	globalRegistry := actions.GlobalRegistry()

	packageManagersCategory, exists := globalRegistry.GetCategory("package_managers")
	if !exists {
		t.Fatal("Package Managers category not found in global registry")
	}

	// Verify package_managers category properties
	tests := []struct {
		name     string
		expected string
		actual   string
	}{
		{
			name:     "category ID",
			expected: "package_managers",
			actual:   packageManagersCategory.ID,
		},
		{
			name:     "category name",
			expected: "Package Managers",
			actual:   packageManagersCategory.Name,
		},
		{
			name:     "category description",
			expected: "Tools and utilities specific to Linux",
			actual:   packageManagersCategory.Description,
		},
	}

	for _, tt := range tests {
		t.Run("Category_"+tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("Expected %s to be '%s', got '%s'", tt.name, tt.expected, tt.actual)
			}
		})
	}

	// Check if all actions are registered
	expectedActions := 3
	if len(packageManagersCategory.Actions) != expectedActions {
		t.Fatalf("Expected %d actions in package_managers category, got %d", expectedActions, len(packageManagersCategory.Actions))
	}

	// Test individual action properties using table-driven approach
	actionTests := []struct {
		name                string
		actionID            string
		expectedName        string
		expectedDescription string
		expectedType        actions.ActionType
		expectedPlatforms   struct {
			visible []actions.Platform
			hidden  []actions.Platform
		}
	}{
		{
			name:                "InstallYay",
			actionID:            "install_yay",
			expectedName:        "Install Yay (AUR Helper)",
			expectedDescription: "Install Yay AUR helper for Arch Linux",
			expectedType:        actions.ActionTypeFunction,
			expectedPlatforms: struct {
				visible []actions.Platform
				hidden  []actions.Platform
			}{
				visible: []actions.Platform{actions.PlatformArch},
				hidden:  nil,
			},
		},
		{
			name:                "InstallBrew",
			actionID:            "install_brew",
			expectedName:        "Install Homebrew",
			expectedDescription: "Install Homebrew on the system",
			expectedType:        actions.ActionTypeFunction,
			expectedPlatforms: struct {
				visible []actions.Platform
				hidden  []actions.Platform
			}{
				visible: nil,
				hidden:  []actions.Platform{actions.PlatformWindows},
			},
		},
		{
			name:                "CheckWinget",
			actionID:            "check_winget",
			expectedName:        "Check Winget",
			expectedDescription: "Check if Winget is installed and show installation instructions if not installed",
			expectedType:        actions.ActionTypeFunction,
			expectedPlatforms: struct {
				visible []actions.Platform
				hidden  []actions.Platform
			}{
				visible: []actions.Platform{actions.PlatformWindows},
				hidden:  nil,
			},
		},
	}

	for _, tt := range actionTests {
		t.Run(tt.name, func(t *testing.T) {
			action := findActionByID(packageManagersCategory.Actions, tt.actionID)
			if action == nil {
				t.Fatalf("%s action not found in package_managers category", tt.name)
			}

			// Test basic properties
			if action.Name != tt.expectedName {
				t.Errorf("Expected %s action name to be '%s', got '%s'", tt.name, tt.expectedName, action.Name)
			}

			if action.Description != tt.expectedDescription {
				t.Errorf("Expected %s action description to be '%s', got '%s'", tt.name, tt.expectedDescription, action.Description)
			}

			if action.Type != tt.expectedType {
				t.Errorf("Expected %s action type to be '%s', got '%s'", tt.name, tt.expectedType, action.Type)
			}

			if action.Handler == nil {
				t.Errorf("%s action has no handler", tt.name)
			}

			// Test platform visibility
			if tt.expectedPlatforms.visible != nil {
				if len(action.VisibleOnPlatforms) != len(tt.expectedPlatforms.visible) {
					t.Errorf("Expected %d visible platforms for %s, got %d", len(tt.expectedPlatforms.visible), tt.name, len(action.VisibleOnPlatforms))
				} else if len(action.VisibleOnPlatforms) > 0 && action.VisibleOnPlatforms[0] != tt.expectedPlatforms.visible[0] {
					t.Errorf("Expected %s action to be visible on '%s' platform, got '%s'", tt.name, tt.expectedPlatforms.visible[0], action.VisibleOnPlatforms[0])
				}
			}

			if tt.expectedPlatforms.hidden != nil {
				if len(action.HiddenOnPlatforms) != len(tt.expectedPlatforms.hidden) {
					t.Errorf("Expected %d hidden platforms for %s, got %d", len(tt.expectedPlatforms.hidden), tt.name, len(action.HiddenOnPlatforms))
				} else if len(action.HiddenOnPlatforms) > 0 && action.HiddenOnPlatforms[0] != tt.expectedPlatforms.hidden[0] {
					t.Errorf("Expected %s action to be hidden on '%s' platform, got '%s'", tt.name, tt.expectedPlatforms.hidden[0], action.HiddenOnPlatforms[0])
				}
			}
		})
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

	// Test all actions are properly registered with required fields
	for _, action := range packageManagersCategory.Actions {
		t.Run("Action_"+action.ID, func(t *testing.T) {
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
		})
	}
}

func TestTestModeFunctionality(t *testing.T) {
	// Test that test mode can be set and unset
	SetTestMode(true)
	if !testMode {
		t.Error("Test mode should be true after setting")
	}

	SetTestMode(false)
	if testMode {
		t.Error("Test mode should be false after unsetting")
	}
}

// Helper function to find an action by ID in a slice of actions
func findActionByID(actions []*actions.Action, id string) *actions.Action {
	for _, action := range actions {
		if action.ID == id {
			return action
		}
	}
	return nil
}
