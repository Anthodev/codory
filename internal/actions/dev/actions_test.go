package dev

import (
	"context"
	"testing"

	"anthodev/codory/internal/actions"
)

func TestInit(t *testing.T) {
	// Since init() is automatically called when the package is imported,
	// we need to test that the global registry has been properly initialized

	// Get the global registry
	globalRegistry := actions.GlobalRegistry()

	// Check if the dev category exists
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

	// Check if all expected actions are registered (UUIDv4, UUIDv7, DecodeUUIDv7, SymfonySecret)
	if len(devCategory.Actions) != 4 {
		t.Fatalf("Expected 4 actions in dev category, got %d", len(devCategory.Actions))
	}

	// Verify all expected actions are present
	expectedActions := map[string]string{
		"uuidv4":         "Generate UUIDv4",
		"uuidv7":         "Generate UUIDv7",
		"decode_uuidv7":  "Decode UUIDv7",
		"symfony_secret": "Generate Symfony secret",
	}

	for actionID, expectedName := range expectedActions {
		found := false
		for _, action := range devCategory.Actions {
			if action.ID == actionID {
				found = true
				if action.Name != expectedName {
					t.Errorf("Expected %s action name to be '%s', got '%s'", actionID, expectedName, action.Name)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected action %s not found in dev category", actionID)
		}
	}
}

func TestDevCategoryIsChildOfRoot(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the root category
	rootCategory := globalRegistry.GetRoot()
	if rootCategory == nil {
		t.Fatal("Root category not found")
	}

	// Check if dev category is a subcategory of root
	devCategoryFound := false
	for _, subCat := range rootCategory.SubCategories {
		if subCat.ID == "dev" {
			devCategoryFound = true
			break
		}
	}

	if !devCategoryFound {
		t.Error("Dev category is not a subcategory of root")
	}
}

func TestDevCategoryActionsAreProperlyRegistered(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the dev category
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	// Verify that all actions in the dev category are properly initialized
	for _, action := range devCategory.Actions {
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

func TestDevCategoryHasNoSubCategories(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the dev category
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	// Verify that dev category has no subcategories (it's a leaf category)
	if !devCategory.IsLeaf() {
		t.Error("Dev category should be a leaf category with no subcategories")
	}

	if len(devCategory.SubCategories) != 0 {
		t.Errorf("Dev category should have no subcategories, got %d", len(devCategory.SubCategories))
	}
}

func TestDevCategoryHasActions(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the dev category
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	// Verify that dev category has actions
	if !devCategory.HasActions() {
		t.Error("Dev category should have actions")
	}

	if len(devCategory.Actions) == 0 {
		t.Error("Dev category should have at least one action")
	}
}

func TestDevCategoryPlatformVisibility(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the dev category
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	// Test visibility on different platforms
	platforms := []actions.Platform{
		actions.PlatformAny,
		actions.PlatformLinux,
		actions.PlatformWindows,
		actions.PlatformMacOS,
		actions.PlatformDebian,
		actions.PlatformArch,
	}

	for _, platform := range platforms {
		if !devCategory.IsVisibleOnPlatform(platform) {
			t.Errorf("Dev category should be visible on platform %s", platform)
		}
		if devCategory.IsHiddenOnPlatform(platform) {
			t.Errorf("Dev category should not be hidden on platform %s", platform)
		}
	}
}

func TestDevActionsHaveCorrectTypes(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the dev category
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	// Verify that all actions are function-type actions
	for _, action := range devCategory.Actions {
		if action.Type != actions.ActionTypeFunction {
			t.Errorf("Action %s should be of type Function, got %s", action.ID, action.Type)
		}
	}
}

func TestDecodeUUIDv7ActionHasArguments(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the dev category
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	// Find the decode_uuidv7 action
	var decodeAction *actions.Action
	for _, action := range devCategory.Actions {
		if action.ID == "decode_uuidv7" {
			decodeAction = action
			break
		}
	}

	if decodeAction == nil {
		t.Fatal("DecodeUUIDv7 action not found")
	}

	// Verify it has the required argument
	if len(decodeAction.Arguments) != 1 {
		t.Fatalf("DecodeUUIDv7 action should have 1 argument, got %d", len(decodeAction.Arguments))
	}

	arg := decodeAction.Arguments[0]
	if arg.Name != "uuidv7" {
		t.Errorf("Expected argument name to be 'uuidv7', got '%s'", arg.Name)
	}
	if arg.Description != "The UUIDv7 to decode (36 characters maximum)" {
		t.Errorf("Expected argument description to match, got '%s'", arg.Description)
	}
	if !arg.Required {
		t.Error("UUIDv7 argument should be required")
	}
}

func TestDevActionsHandlersAreCallable(t *testing.T) {
	globalRegistry := actions.GlobalRegistry()

	// Get the dev category
	devCategory, exists := globalRegistry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	// Test that all function handlers are callable
	for _, action := range devCategory.Actions {
		if action.Type == actions.ActionTypeFunction && action.Handler != nil {
			// Create a test context
			ctx := context.Background()

			// For decode_uuidv7, we need to provide arguments
			if action.ID == "decode_uuidv7" {
				ctx = context.WithValue(ctx, actions.ArgsContextKey, []string{"018f4c1c-9b7a-7f91-8e8f-933f3b3b3b3b"})
			}

			// Call the handler - we expect it to not panic
			result, err := action.Handler(ctx)

			// For UUID actions, we expect success
			if action.ID == "uuidv4" || action.ID == "uuidv7" || action.ID == "symfony_secret" {
				if err != nil {
					t.Errorf("Action %s handler returned error: %v", action.ID, err)
				}
				if result == "" {
					t.Errorf("Action %s handler returned empty result", action.ID)
				}
			}

			// For decode action with valid UUID, we expect success
			if action.ID == "decode_uuidv7" {
				if err != nil {
					t.Errorf("DecodeUUIDv7 handler returned error with valid UUID: %v", err)
				}
				if result == "" {
					t.Errorf("DecodeUUIDv7 handler returned empty result")
				}
			}
		}
	}
}
