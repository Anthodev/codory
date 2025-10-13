package shells

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestInit(t *testing.T) {
	// Set test mode to prevent actual system operations
	SetTestMode(true)
	defer SetTestMode(false)

	globalRegistry := actions.GlobalRegistry()

	shellsCategory, exists := globalRegistry.GetCategory("shells")
	if !exists {
		t.Fatal("Shells category not found in global registry")
	}

	// Verify shells category properties
	tests := []struct {
		name     string
		expected string
		actual   string
	}{
		{
			name:     "category ID",
			expected: "shells",
			actual:   shellsCategory.ID,
		},
		{
			name:     "category name",
			expected: "Shells",
			actual:   shellsCategory.Name,
		},
		{
			name:     "category description",
			expected: "Install and manage different shells",
			actual:   shellsCategory.Description,
		},
	}

	for _, tt := range tests {
		t.Run("Category_"+tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("Expected %s to be '%s', got '%s'", tt.name, tt.expected, tt.actual)
			}
		})
	}

	// Verify shells category has the InstallAtuin action registered
	if len(shellsCategory.Actions) != 1 {
		t.Errorf("Expected shells category to have 1 action (InstallAtuin), got %d", len(shellsCategory.Actions))
	}

	// Verify the action is InstallAtuin
	if len(shellsCategory.Actions) > 0 && shellsCategory.Actions[0].ID != "install_atuin" {
		t.Errorf("Expected first action to be 'install_atuin', got '%s'", shellsCategory.Actions[0].ID)
	}

	if len(shellsCategory.SubCategories) != 0 {
		t.Errorf("Expected shells category to have 0 subcategories initially, got %d", len(shellsCategory.SubCategories))
	}
}

func TestShellsCategoryIsChildOfRoot(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the root category
	rootCategory := globalRegistry.GetRoot()
	if rootCategory == nil {
		t.Fatal("Root category not found")
	}

	// Check if shells category is a subcategory of root
	shellsCategoryFound := false
	for _, subCat := range rootCategory.SubCategories {
		if subCat.ID == "shells" {
			shellsCategoryFound = true
			break
		}
	}

	if !shellsCategoryFound {
		t.Error("Shells category is not a subcategory of root")
	}
}

func TestShellsCategoryStructure(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	shellsCategory, exists := globalRegistry.GetCategory("shells")
	if !exists {
		t.Fatal("Shells category not found")
	}

	// Test that shells category is properly initialized
	if shellsCategory == nil {
		t.Fatal("Shells category is nil")
	}

	// Verify that actions and subcategories slices are initialized (not nil)
	if shellsCategory.Actions == nil {
		t.Error("Shells category Actions slice is nil")
	}

	if shellsCategory.SubCategories == nil {
		t.Error("Shells category SubCategories slice is nil")
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

// Test mode variable and function (similar to package_managers)
var testMode = false

func SetTestMode(enabled bool) {
	testMode = enabled
}
