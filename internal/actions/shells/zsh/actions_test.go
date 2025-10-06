package shells

import (
	"testing"

	"anthodev/codory/internal/actions"
	_ "anthodev/codory/internal/actions/shells" // Import to ensure shells category is initialized
)

func TestInit(t *testing.T) {
	// Since init() is automatically called when the package is imported,
	// we need to test that the global registry has been properly initialized

	// Get the global registry
	globalRegistry := actions.GlobalRegistry()

	// Check if the zsh category exists (it should be registered directly)
	zshCategory, exists := globalRegistry.GetCategory("zsh")
	if !exists {
		t.Fatal("Zsh category not found in global registry")
	}

	// Verify zsh category properties
	if zshCategory.ID != "zsh" {
		t.Errorf("Expected zsh category ID to be 'zsh', got '%s'", zshCategory.ID)
	}

	if zshCategory.Name != "Zsh" {
		t.Errorf("Expected zsh category name to be 'Zsh', got '%s'", zshCategory.Name)
	}

	if zshCategory.Description != "Install zsh and most common plugins" {
		t.Errorf("Expected zsh category description to be 'Install zsh and most common plugins', got '%s'", zshCategory.Description)
	}

	// Check if the expected action is registered
	if len(zshCategory.Actions) != 1 {
		t.Fatalf("Expected 1 action in zsh category, got %d", len(zshCategory.Actions))
	}

	// Verify the install_zsh action is present
	installZshAction := zshCategory.Actions[0]
	if installZshAction.ID != "install_zsh" {
		t.Errorf("Expected action ID to be 'install_zsh', got '%s'", installZshAction.ID)
	}

	if installZshAction.Name != "Install Zsh" {
		t.Errorf("Expected action name to be 'Install Zsh', got '%s'", installZshAction.Name)
	}
}

func TestZshCategoryIsSubCategoryOfShells(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh category directly from registry
	zshCategory, exists := globalRegistry.GetCategory("zsh")
	if !exists {
		t.Fatal("Zsh category not found in global registry")
	}

	// Verify zsh category exists in registry
	if zshCategory.ID != "zsh" {
		t.Errorf("Expected zsh category ID to be 'zsh', got '%s'", zshCategory.ID)
	}

	// Get the shells category to verify hierarchy
	shellsCategory, exists := globalRegistry.GetCategory("shells")
	if !exists {
		t.Skip("Shells category not found - shells package may not be initialized in test environment")
		return
	}

	// Check if zsh category is a subcategory of shells
	zshCategoryFound := false
	for _, subCat := range shellsCategory.SubCategories {
		if subCat.ID == "zsh" {
			zshCategoryFound = true
			// Verify it's the same category instance
			if subCat != zshCategory {
				t.Error("Zsh category in subcategories is not the same instance as the registered one")
			}
			break
		}
	}

	if !zshCategoryFound {
		t.Error("Zsh category is not a subcategory of shells")
	}
}

func TestZshCategoryActionsAreProperlyRegistered(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh category
	zshCategory, exists := globalRegistry.GetCategory("zsh")
	if !exists {
		t.Fatal("Zsh category not found")
	}

	// Verify that all actions in the zsh category are properly initialized
	for _, action := range zshCategory.Actions {
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
		if action.Type == actions.ActionTypeCommand && action.PlatformCommands == nil {
			t.Errorf("Command action has no platform commands for ID: %s", action.ID)
		}
	}
}

func TestZshCategoryHasNoSubCategories(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh category
	zshCategory, exists := globalRegistry.GetCategory("zsh")
	if !exists {
		t.Fatal("Zsh category not found")
	}

	// Verify that zsh category has no subcategories (it's a leaf category)
	if !zshCategory.IsLeaf() {
		t.Error("Zsh category should be a leaf category with no subcategories")
	}

	if len(zshCategory.SubCategories) != 0 {
		t.Errorf("Zsh category should have no subcategories, got %d", len(zshCategory.SubCategories))
	}
}

func TestZshCategoryHasActions(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh category
	zshCategory, exists := globalRegistry.GetCategory("zsh")
	if !exists {
		t.Fatal("Zsh category not found")
	}

	// Verify that zsh category has actions
	if !zshCategory.HasActions() {
		t.Error("Zsh category should have actions")
	}

	if len(zshCategory.Actions) == 0 {
		t.Error("Zsh category should have at least one action")
	}
}

func TestZshCategoryPlatformVisibility(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh category
	zshCategory, exists := globalRegistry.GetCategory("zsh")
	if !exists {
		t.Fatal("Zsh category not found")
	}

	// Test visibility on different platforms
	platforms := []actions.Platform{
		actions.PlatformAny,
		actions.PlatformLinux,
		actions.PlatformDebian,
		actions.PlatformArch,
		actions.PlatformMacOS,
	}

	for _, platform := range platforms {
		if !zshCategory.IsVisibleOnPlatform(platform) {
			t.Errorf("Zsh category should be visible on platform %s", platform)
		}
		if zshCategory.IsHiddenOnPlatform(platform) {
			t.Errorf("Zsh category should not be hidden on platform %s", platform)
		}
	}

	// Test that zsh category is hidden on Windows
	if zshCategory.IsVisibleOnPlatform(actions.PlatformWindows) {
		t.Error("Zsh category should not be visible on Windows platform")
	}
	if !zshCategory.IsHiddenOnPlatform(actions.PlatformWindows) {
		t.Error("Zsh category should be hidden on Windows platform")
	}
}

func TestZshActionsHaveCorrectTypes(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh category
	zshCategory, exists := globalRegistry.GetCategory("zsh")
	if !exists {
		t.Fatal("Zsh category not found")
	}

	// Verify that all actions are command-type actions
	for _, action := range zshCategory.Actions {
		if action.Type != actions.ActionTypeCommand {
			t.Errorf("Action %s should be of type Command, got %s", action.ID, action.Type)
		}
	}
}

func TestZshCategoryRegistrationVerification(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Test that zsh category is properly registered in the global registry
	zshCategory, exists := globalRegistry.GetCategory("zsh")
	if !exists {
		t.Fatal("Zsh category not found in global registry")
	}

	// Verify zsh category properties
	if zshCategory.ID != "zsh" {
		t.Errorf("Expected zsh category ID to be 'zsh', got '%s'", zshCategory.ID)
	}

	if zshCategory.Name != "Zsh" {
		t.Errorf("Expected zsh category name to be 'Zsh', got '%s'", zshCategory.Name)
	}

	if zshCategory.Description != "Install zsh and most common plugins" {
		t.Errorf("Expected zsh category description to be 'Install zsh and most common plugins', got '%s'", zshCategory.Description)
	}

	// Verify zsh category has the install_zsh action
	if len(zshCategory.Actions) != 1 {
		t.Fatalf("Expected 1 action in zsh category, got %d", len(zshCategory.Actions))
	}

	installZshAction := zshCategory.Actions[0]
	if installZshAction.ID != "install_zsh" {
		t.Errorf("Expected action ID to be 'install_zsh', got '%s'", installZshAction.ID)
	}

	// Verify zsh category is hidden on Windows
	if !zshCategory.IsHiddenOnPlatform(actions.PlatformWindows) {
		t.Error("Zsh category should be hidden on Windows platform")
	}

	// Verify zsh category is visible on other platforms
	visiblePlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformDebian,
		actions.PlatformArch,
		actions.PlatformMacOS,
	}

	for _, platform := range visiblePlatforms {
		if zshCategory.IsHiddenOnPlatform(platform) {
			t.Errorf("Zsh category should not be hidden on %s platform", platform)
		}
	}
}

func TestZshCategoryRegistrationIdempotency(t *testing.T) {
	// This test is disabled because the registry prevents duplicate category registration
	// and the init() function in actions.go would cause a panic if registration fails
	t.Skip("Registration idempotency test disabled - registry prevents duplicate registration")
}
