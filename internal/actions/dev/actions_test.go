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

	// Check if the UUIDv4 action is registered
	if len(devCategory.Actions) != 1 {
		t.Fatalf("Expected 1 action in dev category, got %d", len(devCategory.Actions))
	}

	uuidv4Action := devCategory.Actions[0]
	if uuidv4Action.ID != "uuidv4" {
		t.Errorf("Expected action ID to be 'uuidv4', got '%s'", uuidv4Action.ID)
	}

	if uuidv4Action.Name != "Generate UUIDv4" {
		t.Errorf("Expected action name to be 'Generate UUIDv4', got '%s'", uuidv4Action.Name)
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
