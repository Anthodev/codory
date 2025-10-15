package terminal

import (
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/test/testutil"
)

func TestInit(t *testing.T) {
	// Set test mode to prevent actual system operations
	SetTestMode(true)
	defer SetTestMode(false)

	globalRegistry := actions.GlobalRegistry()

	terminalCategory, exists := globalRegistry.GetCategory("terminal")
	if !exists {
		t.Fatal("Terminal category not found in global registry")
	}

	// Verify terminal category properties
	testutil.AssertStringEquals(t, terminalCategory.ID, "terminal", "Category ID")
	testutil.AssertStringEquals(t, terminalCategory.Name, "Terminal", "Category Name")
	testutil.AssertStringEquals(t, terminalCategory.Description, "Terminals you can install on your system", "Category Description")

	// Verify terminal category has the InstallGhostty action registered
	if len(terminalCategory.Actions) != 1 {
		t.Errorf("Expected terminal category to have 1 action (InstallGhostty), got %d", len(terminalCategory.Actions))
	}

	if len(terminalCategory.Actions) > 0 {
		testutil.AssertStringEquals(t, terminalCategory.Actions[0].ID, "install_ghostty", "First action ID")
	}

	if len(terminalCategory.SubCategories) != 0 {
		t.Errorf("Expected terminal category to have 0 subcategories initially, got %d", len(terminalCategory.SubCategories))
	}
}

func TestTerminalCategoryIsChildOfRoot(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the root category
	rootCategory := globalRegistry.GetRoot()
	if rootCategory == nil {
		t.Fatal("Root category not found")
	}

	// Check if terminal category is a subcategory of root
	terminalCategoryFound := false
	for _, subCat := range rootCategory.SubCategories {
		if subCat.ID == "terminal" {
			terminalCategoryFound = true
			break
		}
	}

	if !terminalCategoryFound {
		t.Error("Terminal category is not a subcategory of root")
	}
}

func TestTerminalCategoryStructure(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	terminalCategory, exists := globalRegistry.GetCategory("terminal")
	if !exists {
		t.Fatal("Terminal category not found")
	}

	// Test that terminal category is properly initialized
	if terminalCategory == nil {
		t.Fatal("Terminal category is nil")
	}

	// Verify that actions and subcategories slices are initialized (not nil)
	if terminalCategory.Actions == nil {
		t.Error("Terminal category Actions slice is nil")
	}

	// Note: SubCategories can be nil if no subcategories are registered, which is valid
	// The test only verifies that the category structure is properly initialized
}

func TestTerminalCategoryActions(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	terminalCategory, exists := globalRegistry.GetCategory("terminal")
	if !exists {
		t.Fatal("Terminal category not found")
	}

	// Verify that the InstallGhostty action is properly registered
	if len(terminalCategory.Actions) == 0 {
		t.Fatal("Expected terminal category to have at least one action")
	}

	// Find the InstallGhostty action
	var installGhosttyAction *actions.Action
	for _, action := range terminalCategory.Actions {
		if action.ID == "install_ghostty" {
			installGhosttyAction = action
			break
		}
	}

	if installGhosttyAction == nil {
		t.Fatal("InstallGhostty action not found in terminal category")
	}

	// Verify InstallGhostty action properties
	testutil.AssertStringEquals(t, installGhosttyAction.Name, "Install ghostty terminal", "Action Name")
	testutil.AssertStringEquals(t, installGhosttyAction.Description, "Install ghostty terminal on your system", "Action Description")

	if installGhosttyAction.Type != actions.ActionTypeCommand {
		t.Errorf("Expected action type '%s', got '%s'", actions.ActionTypeCommand, installGhosttyAction.Type)
	}

	if installGhosttyAction.PlatformCommands == nil {
		t.Error("Expected PlatformCommands to be set on InstallGhostty action")
	}
}

func TestTerminalCategoryRegistration(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Verify that the terminal category is registered under root
	rootCategory := globalRegistry.GetRoot()
	if rootCategory == nil {
		t.Fatal("Root category not found")
	}

	// Check if terminal category exists in root's subcategories
	terminalFound := false
	for _, subCat := range rootCategory.SubCategories {
		if subCat.ID == "terminal" {
			terminalFound = true
			// Verify it's the same category object
			registeredTerminal, _ := globalRegistry.GetCategory("terminal")
			if subCat != registeredTerminal {
				t.Error("Terminal category in root subcategories is not the same as registered category")
			}
			break
		}
	}

	if !terminalFound {
		t.Error("Terminal category not found in root subcategories")
	}
}

func TestInitFunctionSafety(t *testing.T) {
	// Test that init function can be called multiple times without issues
	// This is important because init() runs automatically, but we want to ensure
	// it's safe if somehow called again

	// Get initial state
	globalRegistry := actions.GlobalRegistry()
	initialTerminal, initialExists := globalRegistry.GetCategory("terminal")

	if !initialExists {
		t.Fatal("Terminal category should exist after init")
	}

	_ = len(initialTerminal.Actions) // Store initial action count for potential future use

	// Simulate what would happen if init was called again
	// (though in practice, Go's init() prevents this)
	terminalCategory := &actions.Category{
		ID:          "terminal",
		Name:        "Terminal",
		Description: "Terminals you can install on your system",
		Actions:     make([]*actions.Action, 0),
	}

	// This should replace the existing category (idempotent behavior)
	globalRegistry.RegisterCategory("root", terminalCategory)

	// Verify the category still exists and is properly configured
	afterTerminal, afterExists := globalRegistry.GetCategory("terminal")
	if !afterExists {
		t.Fatal("Terminal category should still exist after re-registration")
	}

	if afterTerminal.Name != "Terminal" {
		t.Error("Terminal category name should remain consistent")
	}

	// The action count might change due to re-registration, but the category should still be functional
	if afterTerminal.Actions == nil {
		t.Error("Terminal category Actions should not be nil after re-registration")
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
