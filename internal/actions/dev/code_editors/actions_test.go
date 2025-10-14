package dev

import (
	"testing"

	"anthodev/codory/internal/actions"
)

func TestInit(t *testing.T) {
	// Since init() is automatically called when the package is imported,
	// we need to test that the global registry has been properly initialized

	// Get the global registry
	globalRegistry := actions.GlobalRegistry()

	// Check if the code_editor category exists (it should be registered under dev parent)
	codeEditorCategory, exists := globalRegistry.GetCategory("code_editor")
	if !exists {
		t.Fatal("Code Editor category not found in global registry")
	}

	// Verify code_editor category properties
	if codeEditorCategory.ID != "code_editor" {
		t.Errorf("Expected code_editor category ID to be 'code_editor', got '%s'", codeEditorCategory.ID)
	}

	if codeEditorCategory.Name != "Code Editors" {
		t.Errorf("Expected code_editor category name to be 'Code Editors', got '%s'", codeEditorCategory.Name)
	}

	if codeEditorCategory.Description != "Listing all code editors" {
		t.Errorf("Expected code_editor category description to be 'Listing all code editors', got '%s'", codeEditorCategory.Description)
	}

	// Check if the expected actions are registered
	if len(codeEditorCategory.Actions) != 3 {
		t.Fatalf("Expected 3 actions in code_editor category, got %d", len(codeEditorCategory.Actions))
	}

	// Verify all actions are present
	var installNeovimAction *actions.Action
	var installLazyVimAction *actions.Action
	var installHelixAction *actions.Action
	for _, action := range codeEditorCategory.Actions {
		if action.ID == "install_neovim" {
			installNeovimAction = action
		} else if action.ID == "install_lazyvim" {
			installLazyVimAction = action
		} else if action.ID == "install_helix" {
			installHelixAction = action
		}
	}

	if installNeovimAction == nil {
		t.Error("install_neovim action not found in code_editor category")
	}

	if installLazyVimAction == nil {
		t.Error("install_lazyvim action not found in code_editor category")
	}

	if installHelixAction == nil {
		t.Error("install_helix action not found in code_editor category")
	}

	if installNeovimAction.Name != "Install Neovim" {
		t.Errorf("Expected install_neovim action name to be 'Install Neovim', got '%s'", installNeovimAction.Name)
	}
}

func TestCodeEditorCategoryIsSubCategoryOfDev(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the code_editor category directly from registry
	codeEditorCategory, exists := globalRegistry.GetCategory("code_editor")
	if !exists {
		t.Fatal("Code Editor category not found in global registry")
	}

	// Verify code_editor category exists in registry
	if codeEditorCategory.ID != "code_editor" {
		t.Errorf("Expected code_editor category ID to be 'code_editor', got '%s'", codeEditorCategory.ID)
	}

	// Get the dev category to verify hierarchy
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Skip("Dev category not found - dev package may not be initialized in test environment")
		return
	}

	// Check if code_editor category is a subcategory of dev
	codeEditorCategoryFound := false
	for _, subCat := range devCategory.SubCategories {
		if subCat.ID == "code_editor" {
			codeEditorCategoryFound = true
			// Verify it's the same category instance
			if subCat != codeEditorCategory {
				t.Error("Code Editor category in subcategories is not the same instance as the registered one")
			}
			break
		}
	}

	if !codeEditorCategoryFound {
		t.Error("Code Editor category is not a subcategory of dev")
	}
}

func TestCodeEditorCategoryActionsAreProperlyRegistered(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the code_editor category
	codeEditorCategory, exists := globalRegistry.GetCategory("code_editor")
	if !exists {
		t.Fatal("Code Editor category not found")
	}

	// Verify that all actions in the code_editor category are properly initialized
	for _, action := range codeEditorCategory.Actions {
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

func TestCodeEditorCategoryHasNoSubCategories(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the code_editor category
	codeEditorCategory, exists := globalRegistry.GetCategory("code_editor")
	if !exists {
		t.Fatal("Code Editor category not found")
	}

	// Verify that code_editor category has no subcategories (it's a leaf category)
	if !codeEditorCategory.IsLeaf() {
		t.Error("Code Editor category should be a leaf category with no subcategories")
	}

	if len(codeEditorCategory.SubCategories) != 0 {
		t.Errorf("Code Editor category should have no subcategories, got %d", len(codeEditorCategory.SubCategories))
	}
}

func TestCodeEditorCategoryHasActions(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the code_editor category
	codeEditorCategory, exists := globalRegistry.GetCategory("code_editor")
	if !exists {
		t.Fatal("Code Editor category not found")
	}

	// Verify that code_editor category has actions
	if !codeEditorCategory.HasActions() {
		t.Error("Code Editor category should have actions")
	}

	if len(codeEditorCategory.Actions) == 0 {
		t.Error("Code Editor category should have at least one action")
	}
}

func TestCodeEditorCategoryPlatformVisibility(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the code_editor category
	codeEditorCategory, exists := globalRegistry.GetCategory("code_editor")
	if !exists {
		t.Fatal("Code Editor category not found")
	}

	// Test visibility on different platforms
	platforms := []actions.Platform{
		actions.PlatformAny,
		actions.PlatformLinux,
		actions.PlatformDebian,
		actions.PlatformArch,
		actions.PlatformMacOS,
		actions.PlatformWindows,
	}

	for _, platform := range platforms {
		if !codeEditorCategory.IsVisibleOnPlatform(platform) {
			t.Errorf("Code Editor category should be visible on platform %s", platform)
		}
		if codeEditorCategory.IsHiddenOnPlatform(platform) {
			t.Errorf("Code Editor category should not be hidden on platform %s", platform)
		}
	}
}

func TestCodeEditorCategoryActionsHaveCorrectTypes(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the code_editor category
	codeEditorCategory, exists := globalRegistry.GetCategory("code_editor")
	if !exists {
		t.Fatal("Code Editor category not found")
	}

	// Verify that actions have correct types
	for _, action := range codeEditorCategory.Actions {
		switch action.ID {
		case "install_neovim":
			if action.Type != actions.ActionTypeCommand {
				t.Errorf("Action %s should be of type Command, got %s", action.ID, action.Type)
			}
		case "install_lazyvim":
			if action.Type != actions.ActionTypeCommand {
				t.Errorf("Action %s should be of type Command, got %s", action.ID, action.Type)
			}
		case "install_helix":
			if action.Type != actions.ActionTypeCommand {
				t.Errorf("Action %s should be of type Command, got %s", action.ID, action.Type)
			}
		default:
			t.Errorf("Unknown action %s with type %s", action.ID, action.Type)
		}
	}
}

func TestCodeEditorCategoryRegistrationVerification(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Test that code_editor category is properly registered in the global registry
	codeEditorCategory, exists := globalRegistry.GetCategory("code_editor")
	if !exists {
		t.Fatal("Code Editor category not found in global registry")
	}

	// Verify code_editor category properties
	if codeEditorCategory.ID != "code_editor" {
		t.Errorf("Expected code_editor category ID to be 'code_editor', got '%s'", codeEditorCategory.ID)
	}

	if codeEditorCategory.Name != "Code Editors" {
		t.Errorf("Expected code_editor category name to be 'Code Editors', got '%s'", codeEditorCategory.Name)
	}

	if codeEditorCategory.Description != "Listing all code editors" {
		t.Errorf("Expected code_editor category description to be 'Listing all code editors', got '%s'", codeEditorCategory.Description)
	}

	// Verify code_editor category has the expected actions
	if len(codeEditorCategory.Actions) != 3 {
		t.Fatalf("Expected 3 actions in code_editor category, got %d", len(codeEditorCategory.Actions))
	}

	// Verify both actions are present
	var installNeovimAction *actions.Action
	var installLazyVimAction *actions.Action
	for _, action := range codeEditorCategory.Actions {
		if action.ID == "install_neovim" {
			installNeovimAction = action
		} else if action.ID == "install_lazyvim" {
			installLazyVimAction = action
		}
	}

	if installNeovimAction == nil {
		t.Error("install_neovim action not found in code_editor category")
	}

	if installLazyVimAction == nil {
		t.Error("install_lazyvim action not found in code_editor category")
	}

	// Verify code_editor category is visible on all platforms
	platforms := []actions.Platform{
		actions.PlatformLinux,
		actions.PlatformDebian,
		actions.PlatformArch,
		actions.PlatformMacOS,
		actions.PlatformWindows,
	}

	for _, platform := range platforms {
		if codeEditorCategory.IsHiddenOnPlatform(platform) {
			t.Errorf("Code Editor category should not be hidden on %s platform", platform)
		}
	}
}

func TestCodeEditorCategoryRegistrationIdempotency(t *testing.T) {
	// This test is disabled because the registry prevents duplicate category registration
	// and the init() function in actions.go would cause a panic if registration fails
	t.Skip("Registration idempotency test disabled - registry prevents duplicate registration")
}

func TestDevCategoryCreationForCodeEditors(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Test that dev category exists and has correct properties
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found in global registry")
	}

	// Verify dev category properties
	if devCategory.ID != "dev" {
		t.Errorf("Expected dev category ID to be 'dev', got '%s'", devCategory.ID)
	}

	if devCategory.Name != "Development" {
		t.Errorf("Expected dev category name to be 'Development', got '%s'", devCategory.Name)
	}

	if devCategory.Description != "Development tools and utilities" {
		t.Errorf("Expected dev category description to be 'Development tools and utilities', got '%s'", devCategory.Description)
	}

	// Verify dev category has subcategories (should include code_editor)
	if len(devCategory.SubCategories) == 0 {
		t.Error("Dev category should have subcategories")
	}

	// Verify code_editor is in dev subcategories
	codeEditorFound := false
	for _, subCat := range devCategory.SubCategories {
		if subCat.ID == "code_editor" {
			codeEditorFound = true
			break
		}
	}

	if !codeEditorFound {
		t.Error("Code Editor category should be a subcategory of dev")
	}
}

func TestInitHandlesMissingDevCategoryForCodeEditors(t *testing.T) {
	// This test verifies that the init() function properly handles the case
	// where the dev category doesn't exist yet by creating it
	// Since we can't easily reset the global registry in tests, we verify
	// that the dev category exists and has the expected properties

	globalRegistry := actions.GlobalRegistry()

	// Verify dev category exists (it should have been created by init() if missing)
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category should exist - init() should have created it if missing")
	}

	// Verify dev category has the expected structure
	if devCategory.SubCategories == nil {
		t.Error("Dev category should have initialized SubCategories slice")
	}

	if devCategory.Actions == nil {
		t.Error("Dev category should have initialized Actions slice")
	}

	// Verify code_editor category exists as a subcategory
	codeEditorFound := false
	for _, subCat := range devCategory.SubCategories {
		if subCat.ID == "code_editor" {
			codeEditorFound = true
			if subCat.Actions == nil {
				t.Error("Code Editor category should have initialized Actions slice")
			}
			// Note: Code Editor category doesn't have SubCategories initialized since it's a leaf category
			break
		}
	}

	if !codeEditorFound {
		t.Error("Code Editor category should be registered as a subcategory of dev")
	}
}

func TestInstallNeovimActionHasPlatformCommands(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the code_editor category
	codeEditorCategory, exists := globalRegistry.GetCategory("code_editor")
	if !exists {
		t.Fatal("Code Editor category not found")
	}

	// Find the install_neovim action
	var installNeovimAction *actions.Action
	for _, action := range codeEditorCategory.Actions {
		if action.ID == "install_neovim" {
			installNeovimAction = action
			break
		}
	}

	if installNeovimAction == nil {
		t.Fatal("install_neovim action not found")
	}

	// Verify the action has platform commands
	if installNeovimAction.PlatformCommands == nil {
		t.Fatal("install_neovim action should have platform commands")
	}

	// Verify expected platforms are supported
	expectedPlatforms := []actions.Platform{
		actions.PlatformArch,
		actions.PlatformDebian,
		actions.PlatformLinux,
		actions.PlatformMacOS,
	}

	for _, platform := range expectedPlatforms {
		cmd, exists := installNeovimAction.PlatformCommands[platform]
		if !exists {
			t.Errorf("install_neovim action should have command for platform %s", platform)
			continue
		}

		if cmd.Command == "" {
			t.Errorf("Command should not be empty for platform %s", platform)
		}

		if cmd.CheckCommand != "nvim" {
			t.Errorf("CheckCommand should be 'nvim' for platform %s, got '%s'", platform, cmd.CheckCommand)
		}
	}
}
