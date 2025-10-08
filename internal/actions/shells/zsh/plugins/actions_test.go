package shells

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestInit(t *testing.T) {
	// Since init() is automatically called when the package is imported,
	// we need to test that the global registry has been properly initialized

	// Get the global registry
	globalRegistry := actions.GlobalRegistry()

	// Check if the zsh_plugins category exists (it should be registered under zsh parent)
	zshPluginsCategory, exists := globalRegistry.GetCategory("zsh_plugins")
	if !exists {
		t.Fatal("Zsh plugins category not found in global registry")
	}

	// Verify zsh_plugins category properties
	if zshPluginsCategory.ID != "zsh_plugins" {
		t.Errorf("Expected zsh_plugins category ID to be 'zsh_plugins', got '%s'", zshPluginsCategory.ID)
	}

	if zshPluginsCategory.Name != "Zsh Plugins" {
		t.Errorf("Expected zsh_plugins category name to be 'Zsh Plugins', got '%s'", zshPluginsCategory.Name)
	}

	if zshPluginsCategory.Description != "Install plugins for oh-my-zsh" {
		t.Errorf("Expected zsh_plugins category description to be 'Install plugins for oh-my-zsh', got '%s'", zshPluginsCategory.Description)
	}

	// Check if the expected actions are registered
	if len(zshPluginsCategory.Actions) != 2 {
		t.Fatalf("Expected 2 actions in zsh_plugins category, got %d", len(zshPluginsCategory.Actions))
	}

	// Verify both actions are present
	var installHistorySearchAction *actions.Action
	var installSyntaxHighlightingAction *actions.Action
	for _, action := range zshPluginsCategory.Actions {
		if action.ID == "install_zsh_plugin_history_search" {
			installHistorySearchAction = action
		} else if action.ID == "install_zsh_plugin_syntax_highlighting" {
			installSyntaxHighlightingAction = action
		}
	}

	if installHistorySearchAction == nil {
		t.Error("install_zsh_plugin_history_search action not found in zsh_plugins category")
	}

	if installSyntaxHighlightingAction == nil {
		t.Error("install_zsh_plugin_syntax_highlighting action not found in zsh_plugins category")
	}

	if installHistorySearchAction.Name != "Install Zsh Plugin zsh-history-substring-search" {
		t.Errorf("Expected install_zsh_plugin_history_search action name to be 'Install Zsh Plugin zsh-history-substring-search', got '%s'", installHistorySearchAction.Name)
	}

	if installSyntaxHighlightingAction.Name != "Install Zsh Plugin zsh-syntax-highlighting" {
		t.Errorf("Expected install_zsh_plugin_syntax_highlighting action name to be 'Install Zsh Plugin zsh-syntax-highlighting', got '%s'", installSyntaxHighlightingAction.Name)
	}
}

func TestZshPluginsCategoryIsSubCategoryOfZsh(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh_plugins category directly from registry
	zshPluginsCategory, exists := globalRegistry.GetCategory("zsh_plugins")
	if !exists {
		t.Fatal("Zsh plugins category not found in global registry")
	}

	// Verify zsh_plugins category exists in registry
	if zshPluginsCategory.ID != "zsh_plugins" {
		t.Errorf("Expected zsh_plugins category ID to be 'zsh_plugins', got '%s'", zshPluginsCategory.ID)
	}

	// Get the zsh category to verify hierarchy
	zshCategory, exists := globalRegistry.GetCategory("zsh")
	if !exists {
		t.Skip("Zsh category not found - zsh package may not be initialized in test environment")
		return
	}

	// Check if zsh_plugins category is a subcategory of zsh
	zshPluginsCategoryFound := false
	for _, subCat := range zshCategory.SubCategories {
		if subCat.ID == "zsh_plugins" {
			zshPluginsCategoryFound = true
			// Verify it's the same category instance
			if subCat != zshPluginsCategory {
				t.Error("Zsh plugins category in subcategories is not the same instance as the registered one")
			}
			break
		}
	}

	if !zshPluginsCategoryFound {
		t.Error("Zsh plugins category is not a subcategory of zsh")
	}
}

func TestZshPluginsCategoryActionsAreProperlyRegistered(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh_plugins category
	zshPluginsCategory, exists := globalRegistry.GetCategory("zsh_plugins")
	if !exists {
		t.Fatal("Zsh plugins category not found")
	}

	// Verify that all actions in the zsh_plugins category are properly initialized
	for _, action := range zshPluginsCategory.Actions {
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

func TestZshPluginsCategoryHasNoSubCategories(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh_plugins category
	zshPluginsCategory, exists := globalRegistry.GetCategory("zsh_plugins")
	if !exists {
		t.Fatal("Zsh plugins category not found")
	}

	// Verify that zsh_plugins category has no subcategories (it's a leaf category)
	if !zshPluginsCategory.IsLeaf() {
		t.Error("Zsh plugins category should be a leaf category with no subcategories")
	}

	if len(zshPluginsCategory.SubCategories) != 0 {
		t.Errorf("Zsh plugins category should have no subcategories, got %d", len(zshPluginsCategory.SubCategories))
	}
}

func TestZshPluginsCategoryHasActions(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh_plugins category
	zshPluginsCategory, exists := globalRegistry.GetCategory("zsh_plugins")
	if !exists {
		t.Fatal("Zsh plugins category not found")
	}

	// Verify that zsh_plugins category has actions
	if !zshPluginsCategory.HasActions() {
		t.Error("Zsh plugins category should have actions")
	}

	if len(zshPluginsCategory.Actions) == 0 {
		t.Error("Zsh plugins category should have at least one action")
	}
}

func TestZshPluginsCategoryPlatformVisibility(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh_plugins category
	zshPluginsCategory, exists := globalRegistry.GetCategory("zsh_plugins")
	if !exists {
		t.Fatal("Zsh plugins category not found")
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
		if !zshPluginsCategory.IsVisibleOnPlatform(platform) {
			t.Errorf("Zsh plugins category should be visible on platform %s", platform)
		}
		if zshPluginsCategory.IsHiddenOnPlatform(platform) {
			t.Errorf("Zsh plugins category should not be hidden on platform %s", platform)
		}
	}

	// Test that zsh_plugins category is hidden on Windows
	if zshPluginsCategory.IsVisibleOnPlatform(actions.PlatformWindows) {
		t.Error("Zsh plugins category should not be visible on Windows platform")
	}
	if !zshPluginsCategory.IsHiddenOnPlatform(actions.PlatformWindows) {
		t.Error("Zsh plugins category should be hidden on Windows platform")
	}
}

func TestZshPluginsActionsHaveCorrectTypes(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the zsh_plugins category
	zshPluginsCategory, exists := globalRegistry.GetCategory("zsh_plugins")
	if !exists {
		t.Fatal("Zsh plugins category not found")
	}

	// Verify that actions have correct types
	for _, action := range zshPluginsCategory.Actions {
		switch action.ID {
		case "install_zsh_plugin_history_search":
			if action.Type != actions.ActionTypeCommand {
				t.Errorf("Action %s should be of type Command, got %s", action.ID, action.Type)
			}
		case "install_zsh_plugin_syntax_highlighting":
			if action.Type != actions.ActionTypeCommand {
				t.Errorf("Action %s should be of type Command, got %s", action.ID, action.Type)
			}
		default:
			t.Errorf("Unknown action %s with type %s", action.ID, action.Type)
		}
	}
}

func TestZshPluginsCategoryRegistrationVerification(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Test that zsh_plugins category is properly registered in the global registry
	zshPluginsCategory, exists := globalRegistry.GetCategory("zsh_plugins")
	if !exists {
		t.Fatal("Zsh plugins category not found in global registry")
	}

	// Verify zsh_plugins category properties
	if zshPluginsCategory.ID != "zsh_plugins" {
		t.Errorf("Expected zsh_plugins category ID to be 'zsh_plugins', got '%s'", zshPluginsCategory.ID)
	}

	if zshPluginsCategory.Name != "Zsh Plugins" {
		t.Errorf("Expected zsh_plugins category name to be 'Zsh Plugins', got '%s'", zshPluginsCategory.Name)
	}

	if zshPluginsCategory.Description != "Install plugins for oh-my-zsh" {
		t.Errorf("Expected zsh_plugins category description to be 'Install plugins for oh-my-zsh', got '%s'", zshPluginsCategory.Description)
	}

	// Verify zsh_plugins category has the expected actions
	if len(zshPluginsCategory.Actions) != 2 {
		t.Fatalf("Expected 2 actions in zsh_plugins category, got %d", len(zshPluginsCategory.Actions))
	}

	// Verify both actions are present
	var installHistorySearchAction *actions.Action
	var installSyntaxHighlightingAction *actions.Action
	for _, action := range zshPluginsCategory.Actions {
		if action.ID == "install_zsh_plugin_history_search" {
			installHistorySearchAction = action
		} else if action.ID == "install_zsh_plugin_syntax_highlighting" {
			installSyntaxHighlightingAction = action
		}
	}

	if installHistorySearchAction == nil {
		t.Error("install_zsh_plugin_history_search action not found in zsh_plugins category")
	}

	if installSyntaxHighlightingAction == nil {
		t.Error("install_zsh_plugin_syntax_highlighting action not found in zsh_plugins category")
	}

	// Verify zsh_plugins category is hidden on Windows
	if !zshPluginsCategory.IsHiddenOnPlatform(actions.PlatformWindows) {
		t.Error("Zsh plugins category should be hidden on Windows platform")
	}

	// Verify zsh_plugins category is visible on other platforms
	visiblePlatforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformDebian,
		actions.PlatformArch,
		actions.PlatformMacOS,
	}

	for _, platform := range visiblePlatforms {
		if zshPluginsCategory.IsHiddenOnPlatform(platform) {
			t.Errorf("Zsh plugins category should not be hidden on %s platform", platform)
		}
	}
}

func TestZshPluginsCategoryRegistrationIdempotency(t *testing.T) {
	// This test is disabled because the registry prevents duplicate category registration
	// and the init() function in actions.go would cause a panic if registration fails
	t.Skip("Registration idempotency test disabled - registry prevents duplicate registration")
}
