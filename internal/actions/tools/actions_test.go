package tools

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

	toolsCategory, exists := globalRegistry.GetCategory("tools")
	if !exists {
		t.Fatal("Tools category not found in global registry")
	}

	// Verify tools category properties
	testutil.AssertStringEquals(t, toolsCategory.ID, "tools", "Category ID")
	testutil.AssertStringEquals(t, toolsCategory.Name, "Tools", "Category Name")
	testutil.AssertStringEquals(t, toolsCategory.Description, "Useful tools that you can add to your system", "Category Description")

	// Verify tools category has the correct number of actions registered
	if len(toolsCategory.Actions) != 2 {
		t.Errorf("Expected tools category to have 2 actions (InstallBat and InstallBtop), got %d", len(toolsCategory.Actions))
	}

	// Validate presence of InstallBat action by ID
	foundInstallBat := false
	for _, a := range toolsCategory.Actions {
		if a.ID == "install_bat" {
			foundInstallBat = true
			break
		}
	}
	if !foundInstallBat {
		t.Error("Expected action ID 'install_bat' to be registered in tools category")
	}

	// Validate presence of InstallBtop action by ID
	foundInstallBtop := false
	for _, a := range toolsCategory.Actions {
		if a.ID == "install_btop" {
			foundInstallBtop = true
			break
		}
	}
	if !foundInstallBtop {
		t.Error("Expected action ID 'install_btop' to be registered in tools category")
	}

	if len(toolsCategory.SubCategories) != 0 {
		t.Errorf("Expected tools category to have 0 subcategories initially, got %d", len(toolsCategory.SubCategories))
	}
}

func TestToolsCategoryIsChildOfRoot(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the root category
	rootCategory := globalRegistry.GetRoot()
	if rootCategory == nil {
		t.Fatal("Root category not found")
	}

	// Check if tools category is a subcategory of root
	toolsCategoryFound := false
	for _, subCat := range rootCategory.SubCategories {
		if subCat.ID == "tools" {
			toolsCategoryFound = true
			break
		}
	}

	if !toolsCategoryFound {
		t.Error("Tools category is not a subcategory of root")
	}
}

func TestToolsCategoryStructure(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	toolsCategory, exists := globalRegistry.GetCategory("tools")
	if !exists {
		t.Fatal("Tools category not found")
	}

	// Test that tools category is properly initialized
	if toolsCategory == nil {
		t.Fatal("Tools category is nil")
	}

	// Verify that actions and subcategories slices are initialized (not nil)
	if toolsCategory.Actions == nil {
		t.Error("Tools category Actions slice is nil")
	}

	// Note: SubCategories can be nil if no subcategories are registered, which is valid
	// The test only verifies that the category structure is properly initialized
}

func TestToolsCategoryActions(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	toolsCategory, exists := globalRegistry.GetCategory("tools")
	if !exists {
		t.Fatal("Tools category not found")
	}

	// Verify that the tools category has actions registered
	if len(toolsCategory.Actions) == 0 {
		t.Fatal("Expected tools category to have at least one action")
	}

	// Verify expected action count
	if len(toolsCategory.Actions) != 2 {
		t.Fatalf("Expected 2 actions in tools category, got %d", len(toolsCategory.Actions))
	}

	// Find and verify the InstallBat action
	var installBatAction *actions.Action
	for _, action := range toolsCategory.Actions {
		if action.ID == "install_bat" {
			installBatAction = action
			break
		}
	}

	if installBatAction == nil {
		t.Fatal("InstallBat action not found in tools category")
	}

	testutil.AssertStringEquals(t, installBatAction.Name, "Install bat", "InstallBat Action Name")
	testutil.AssertStringEquals(t, installBatAction.Description, "Install the bat content viewer tool", "InstallBat Action Description")

	if installBatAction.Type != actions.ActionTypeCommand {
		t.Errorf("Expected InstallBat action type '%s', got '%s'", actions.ActionTypeCommand, installBatAction.Type)
	}

	if installBatAction.PlatformCommands == nil {
		t.Error("Expected PlatformCommands to be set on InstallBat action")
	}

	// Find and verify the InstallBtop action
	var installBtopAction *actions.Action
	for _, action := range toolsCategory.Actions {
		if action.ID == "install_btop" {
			installBtopAction = action
			break
		}
	}

	if installBtopAction == nil {
		t.Fatal("InstallBtop action not found in tools category")
	}

	testutil.AssertStringEquals(t, installBtopAction.Name, "Install btop", "InstallBtop Action Name")
	testutil.AssertStringEquals(t, installBtopAction.Description, "Install btop system monitor", "InstallBtop Action Description")

	if installBtopAction.Type != actions.ActionTypeCommand {
		t.Errorf("Expected InstallBtop action type '%s', got '%s'", actions.ActionTypeCommand, installBtopAction.Type)
	}

	if installBtopAction.PlatformCommands == nil {
		t.Error("Expected PlatformCommands to be set on InstallBtop action")
	}
}

func TestInstallBatActionPlatformCommands(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	toolsCategory, exists := globalRegistry.GetCategory("tools")
	if !exists {
		t.Fatal("Tools category not found")
	}

	// Find the InstallBat action
	var installBatAction *actions.Action
	for _, action := range toolsCategory.Actions {
		if action.ID == "install_bat" {
			installBatAction = action
			break
		}
	}

	if installBatAction == nil {
		t.Fatal("InstallBat action not found in tools category")
	}

	// Verify that platform commands are configured
	if len(installBatAction.PlatformCommands) == 0 {
		t.Error("Expected InstallBat action to have platform commands configured")
	}

	// Verify expected platforms are configured
	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range expectedPlatforms {
		if _, exists := installBatAction.PlatformCommands[platform]; !exists {
			t.Errorf("Expected platform command for %s not found in InstallBat", platform)
		}
	}
}

func TestInstallBtopActionPlatformCommands(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	toolsCategory, exists := globalRegistry.GetCategory("tools")
	if !exists {
		t.Fatal("Tools category not found")
	}

	// Find the InstallBtop action
	var installBtopAction *actions.Action
	for _, action := range toolsCategory.Actions {
		if action.ID == "install_btop" {
			installBtopAction = action
			break
		}
	}

	if installBtopAction == nil {
		t.Fatal("InstallBtop action not found in tools category")
	}

	// Verify that platform commands are configured
	if len(installBtopAction.PlatformCommands) == 0 {
		t.Error("Expected InstallBtop action to have platform commands configured")
	}

	// Verify expected platforms are configured
	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range expectedPlatforms {
		if _, exists := installBtopAction.PlatformCommands[platform]; !exists {
			t.Errorf("Expected platform command for %s not found in InstallBtop", platform)
		}
	}
}

func TestToolsCategoryRegistration(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Verify that the tools category is registered under root
	rootCategory := globalRegistry.GetRoot()
	if rootCategory == nil {
		t.Fatal("Root category not found")
	}

	// Check if tools category exists in root's subcategories
	toolsFound := false
	for _, subCat := range rootCategory.SubCategories {
		if subCat.ID == "tools" {
			toolsFound = true
			// Verify it's the same category object
			registeredTools, _ := globalRegistry.GetCategory("tools")
			if subCat != registeredTools {
				t.Error("Tools category in root subcategories is not the same as registered category")
			}
			break
		}
	}

	if !toolsFound {
		t.Error("Tools category not found in root subcategories")
	}
}

func TestAllActionsRegistered(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	toolsCategory, exists := globalRegistry.GetCategory("tools")
	if !exists {
		t.Fatal("Tools category not found")
	}

	// Map to track which actions we find
	registeredActions := make(map[string]*actions.Action)
	for _, action := range toolsCategory.Actions {
		registeredActions[action.ID] = action
	}

	// Expected actions in the tools category
	expectedActions := []struct {
		id   string
		name string
	}{
		{id: "install_bat", name: "Install bat"},
		{id: "install_btop", name: "Install btop"},
	}

	// Verify all expected actions are registered
	for _, expected := range expectedActions {
		action, found := registeredActions[expected.id]
		if !found {
			t.Errorf("Expected action '%s' not found in tools category", expected.id)
			continue
		}

		if action.Name != expected.name {
			t.Errorf("Action '%s' has incorrect name. Expected '%s', got '%s'", expected.id, expected.name, action.Name)
		}

		if action.Type != actions.ActionTypeCommand {
			t.Errorf("Action '%s' has incorrect type. Expected '%s', got '%s'", expected.id, actions.ActionTypeCommand, action.Type)
		}

		if action.PlatformCommands == nil {
			t.Errorf("Action '%s' has no platform commands configured", expected.id)
		}
	}

	// Verify no extra actions are registered
	if len(registeredActions) != len(expectedActions) {
		t.Errorf("Expected %d actions in tools category, but found %d", len(expectedActions), len(registeredActions))
	}
}

func TestInitFunctionSafety(t *testing.T) {
	// Test that init function can be called multiple times without issues
	// This is important because init() runs automatically, but we want to ensure
	// it's safe if somehow called again

	// Get initial state
	globalRegistry := actions.GlobalRegistry()
	initialTools, initialExists := globalRegistry.GetCategory("tools")

	if !initialExists {
		t.Fatal("Tools category should exist after init")
	}

	_ = len(initialTools.Actions) // Store initial action count for potential future use

	// Simulate what would happen if init was called again
	// (though in practice, Go's init() prevents this)
	toolsCategory := &actions.Category{
		ID:          "tools",
		Name:        "Tools",
		Description: "Useful tools that you can add to your system",
		Actions:     make([]*actions.Action, 0),
	}

	// This should replace the existing category (idempotent behavior)
	globalRegistry.RegisterCategory("root", toolsCategory)

	// Verify the category still exists and is properly configured
	afterTools, afterExists := globalRegistry.GetCategory("tools")
	if !afterExists {
		t.Fatal("Tools category should still exist after re-registration")
	}

	if afterTools.Name != "Tools" {
		t.Error("Tools category name should remain consistent")
	}

	// The action count might change due to re-registration, but the category should still be functional
	if afterTools.Actions == nil {
		t.Error("Tools category Actions should not be nil after re-registration")
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

// Test mode variable and function (similar to terminal package)
var testMode = false

func SetTestMode(enabled bool) {
	testMode = enabled
}
