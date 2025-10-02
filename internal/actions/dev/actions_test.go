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

	// Check if UUIDv4, UUIDv7, and DecodeUUIDv7 actions are registered
	if len(devCategory.Actions) != 3 {
		t.Fatalf("Expected 3 actions in dev category, got %d", len(devCategory.Actions))
	}

	// Verify UUIDv4, UUIDv7, and DecodeUUIDv7 actions
	var uuidv4Action *actions.Action
	var uuidv7Action *actions.Action
	var decodeUUIDv7Action *actions.Action

	for _, action := range devCategory.Actions {
		if action.ID == "uuidv4" {
			uuidv4Action = action
		} else if action.ID == "uuidv7" {
			uuidv7Action = action
		} else if action.ID == "decode_uuidv7" {
			decodeUUIDv7Action = action
		}
	}

	if uuidv4Action == nil {
		t.Fatal("UUIDv4 action not found in dev category")
	}

	if uuidv4Action.Name != "Generate UUIDv4" {
		t.Errorf("Expected UUIDv4 action name to be 'Generate UUIDv4', got '%s'", uuidv4Action.Name)
	}

	// Verify UUIDv7 action
	if uuidv7Action == nil {
		t.Fatal("UUIDv7 action not found in dev category")
	}

	if uuidv7Action.Name != "Generate UUIDv7" {
		t.Errorf("Expected UUIDv7 action name to be 'Generate UUIDv7', got '%s'", uuidv7Action.Name)
	}

	// Verify DecodeUUIDv7 action
	if decodeUUIDv7Action == nil {
		t.Fatal("DecodeUUIDv7 action not found in dev category")
	}

	if decodeUUIDv7Action.Name != "Decode UUIDv7" {
		t.Errorf("Expected DecodeUUIDv7 action name to be 'Decode UUIDv7', got '%s'", decodeUUIDv7Action.Name)
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
