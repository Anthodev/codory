package dev

import (
	"context"
	"testing"

	"anthodev/codory/internal/actions"
)

func TestNewUUIDv4Action(t *testing.T) {
	action := NewUUIDv4Action()

	// Test action properties
	if action.ID != "uuidv4" {
		t.Errorf("Expected action ID to be 'uuidv4', got '%s'", action.ID)
	}

	if action.Name != "Generate UUIDv4" {
		t.Errorf("Expected action name to be 'Generate UUIDv4', got '%s'", action.Name)
	}

	if action.Description != "Generate a new UUIDv4" {
		t.Errorf("Expected action description to be 'Generate a new UUIDv4', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeFunction {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeFunction, action.Type)
	}

	if action.Handler == nil {
		t.Error("Expected action handler to be set")
	}
}

func TestGenerateUUIDv4(t *testing.T) {
	ctx := context.Background()
	result, err := generateUUIDv4(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Test that the result contains "Generated UUID:" and "copied to clipboard!"
	if !contains(result, "Generated UUID:") {
		t.Errorf("Expected result to contain 'Generated UUID:', got '%s'", result)
	}

	if !contains(result, "copied to clipboard!") {
		t.Errorf("Expected result to contain 'copied to clipboard!', got '%s'", result)
	}

	// Test that the UUID part is a valid UUID format (basic check)
	// Extract UUID from result string
	uuidStart := findSubstringIndex(result, "Generated UUID: ")
	if uuidStart == -1 {
		t.Error("Could not find UUID in result")
		return
	}
	uuidStart += len("Generated UUID: ")

	uuidEnd := findSubstringIndex(result[uuidStart:], ", copied to clipboard!")
	if uuidEnd == -1 {
		t.Error("Could not find end of UUID in result")
		return
	}

	uuidStr := result[uuidStart : uuidStart+uuidEnd]

	// Basic UUID format check (8-4-4-4-12 pattern)
	if len(uuidStr) != 36 {
		t.Errorf("Expected UUID length to be 36, got %d", len(uuidStr))
	}

	// Check hyphens at positions 8, 13, 18, 23
	expectedHyphens := []int{8, 13, 18, 23}
	for _, pos := range expectedHyphens {
		if uuidStr[pos] != '-' {
			t.Errorf("Expected hyphen at position %d, got '%c'", pos, uuidStr[pos])
		}
	}
}

func TestGenerateUUIDv4_ClipboardError(t *testing.T) {
	// This test verifies that the function still generates a UUID even if clipboard fails
	// We can't easily mock the clipboard in the current implementation, but we can
	// at least verify that the function doesn't panic and returns a valid result

	// The current implementation uses log.Fatal on clipboard error, which would stop the test
	// In a real application, this should be handled more gracefully
	// For now, we'll just ensure the basic functionality works

	ctx := context.Background()
	result, err := generateUUIDv4(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result even if clipboard fails")
	}
}

func TestUUIDv4ActionRegistration(t *testing.T) {
	registry := actions.NewRegistry()

	// Create the Dev category
	devCategory := &actions.Category{
		ID:          "dev",
		Name:        "Development",
		Description: "Development tools and utilities",
		Actions:     make([]*actions.Action, 0),
	}

	err := registry.RegisterCategory("root", devCategory)
	if err != nil {
		t.Fatalf("Failed to register dev category: %v", err)
	}

	// Register the UUIDv4 action
	uuidv4Action := NewUUIDv4Action()
	err = registry.RegisterAction("dev", uuidv4Action)
	if err != nil {
		t.Fatalf("Failed to register UUIDv4 action: %v", err)
	}

	// Verify the action was registered
	devCat, exists := registry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	if len(devCat.Actions) != 1 {
		t.Errorf("Expected 1 action in dev category, got %d", len(devCat.Actions))
	}

	if devCat.Actions[0].ID != "uuidv4" {
		t.Errorf("Expected action ID to be 'uuidv4', got '%s'", devCat.Actions[0].ID)
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

func TestGenerateUUIDv4_UniqueGeneration(t *testing.T) {
	// Test that multiple calls generate different UUIDs
	ctx := context.Background()

	result1, err1 := generateUUIDv4(ctx)
	if err1 != nil {
		t.Errorf("Expected no error on first call, got %v", err1)
	}

	result2, err2 := generateUUIDv4(ctx)
	if err2 != nil {
		t.Errorf("Expected no error on second call, got %v", err2)
	}

	// Extract just the UUID parts for comparison
	uuid1Start := findSubstringIndex(result1, "Generated UUID: ") + len("Generated UUID: ")
	uuid1End := findSubstringIndex(result1[uuid1Start:], ", copied to clipboard!")
	uuid1 := result1[uuid1Start : uuid1Start+uuid1End]

	uuid2Start := findSubstringIndex(result2, "Generated UUID: ") + len("Generated UUID: ")
	uuid2End := findSubstringIndex(result2[uuid2Start:], ", copied to clipboard!")
	uuid2 := result2[uuid2Start : uuid2Start+uuid2End]

	if uuid1 == uuid2 {
		t.Error("Expected different UUIDs on multiple calls")
	}
}

func findSubstringIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
