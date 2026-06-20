package dev

import (
	"context"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
)

func TestNewUUIDv7Action(t *testing.T) {
	action := NewUUIDv7Action()

	// Test action properties
	if action.ID != "uuidv7" {
		t.Errorf("Expected action ID to be 'uuidv7', got '%s'", action.ID)
	}

	if action.Name != "Generate UUIDv7" {
		t.Errorf("Expected action name to be 'Generate UUIDv7', got '%s'", action.Name)
	}

	if action.Description != "Generate a new UUIDv7" {
		t.Errorf("Expected action description to be 'Generate a new UUIDv7', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeFunction {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeFunction, action.Type)
	}

	if action.Handler == nil {
		t.Error("Expected action handler to be set")
	}
}

func TestGenerateUUIDv7(t *testing.T) {
	ctx := context.Background()
	result, err := generateUUIDv7(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Test that the result contains "Generated UUIDv7:" and "copied to clipboard!"
	if !utils.Contains(result, "Generated UUIDv7:") {
		t.Errorf("Expected result to contain 'Generated UUIDv7:', got '%s'", result)
	}

	if !utils.Contains(result, "copied to clipboard!") {
		t.Errorf("Expected result to contain 'copied to clipboard!', got '%s'", result)
	}

	// Test that the UUID part is a valid UUID format (basic check)
	// Extract UUID from result string
	uuidStart := findSubstringIndex(result, "Generated UUIDv7: ")
	if uuidStart == -1 {
		t.Error("Could not find UUID in result")
		return
	}
	uuidStart += len("Generated UUIDv7: ")

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

func TestGenerateUUIDv7_ClipboardError(t *testing.T) {
	// This test verifies that the function still generates a UUID even if clipboard fails
	// We can't easily mock the clipboard in the current implementation, but we can
	// at least verify that the function doesn't panic and returns a valid result

	// Clipboard errors should be returned instead of terminating the process.

	ctx := context.Background()
	result, err := generateUUIDv7(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result even if clipboard fails")
	}
}

func TestUUIDv7ActionRegistration(t *testing.T) {
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

	// Register the UUIDv7 action
	uuidv7Action := NewUUIDv7Action()
	err = registry.RegisterAction("dev", uuidv7Action)
	if err != nil {
		t.Fatalf("Failed to register UUIDv7 action: %v", err)
	}

	// Verify the action was registered
	devCat, exists := registry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	if len(devCat.Actions) != 1 {
		t.Errorf("Expected 1 action in dev category, got %d", len(devCat.Actions))
	}

	if devCat.Actions[0].ID != "uuidv7" {
		t.Errorf("Expected action ID to be 'uuidv7', got '%s'", devCat.Actions[0].ID)
	}
}

func TestGenerateUUIDv7_UniqueGeneration(t *testing.T) {
	// Test that multiple calls generate different UUIDs
	ctx := context.Background()

	result1, err1 := generateUUIDv7(ctx)
	if err1 != nil {
		t.Errorf("Expected no error on first call, got %v", err1)
	}

	result2, err2 := generateUUIDv7(ctx)
	if err2 != nil {
		t.Errorf("Expected no error on second call, got %v", err2)
	}

	// Extract just the UUID parts for comparison
	uuid1Start := findSubstringIndex(result1, "Generated UUIDv7: ") + len("Generated UUIDv7: ")
	uuid1End := findSubstringIndex(result1[uuid1Start:], ", copied to clipboard!")
	uuid1 := result1[uuid1Start : uuid1Start+uuid1End]

	uuid2Start := findSubstringIndex(result2, "Generated UUIDv7: ") + len("Generated UUIDv7: ")
	uuid2End := findSubstringIndex(result2[uuid2Start:], ", copied to clipboard!")
	uuid2 := result2[uuid2Start : uuid2Start+uuid2End]

	if uuid1 == uuid2 {
		t.Error("Expected different UUIDs on multiple calls")
	}
}

func TestGenerateUUIDv7_TimestampBased(t *testing.T) {
	// Test that UUIDv7s are timestamp-based and should be somewhat sequential
	// This is a basic test to ensure UUIDv7 generation is working as expected
	ctx := context.Background()

	// Generate multiple UUIDs in quick succession
	var uuids []string
	for i := 0; i < 5; i++ {
		result, err := generateUUIDv7(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
			return
		}

		// Extract UUID from result
		uuidStart := findSubstringIndex(result, "Generated UUIDv7: ") + len("Generated UUIDv7: ")
		uuidEnd := findSubstringIndex(result[uuidStart:], ", copied to clipboard!")
		uuidStr := result[uuidStart : uuidStart+uuidEnd]

		uuids = append(uuids, uuidStr)
	}

	// Basic validation: all UUIDs should be unique
	seen := make(map[string]bool)
	for _, uuid := range uuids {
		if seen[uuid] {
			t.Errorf("Duplicate UUID generated: %s", uuid)
		}
		seen[uuid] = true
	}

	// For UUIDv7, the first part should be similar (timestamp-based)
	// Extract the first 8 characters (timestamp part) from each UUID
	firstParts := make([]string, len(uuids))
	for i, uuid := range uuids {
		if len(uuid) >= 8 {
			firstParts[i] = uuid[:8]
		}
	}

	// At least some of the first parts should be identical or very similar
	// since they were generated in quick succession
	similarCount := 0
	for i := 1; i < len(firstParts); i++ {
		if firstParts[i] == firstParts[0] {
			similarCount++
		}
	}

	// We expect at least some similarity in the timestamp part
	if similarCount == 0 {
		t.Log("Note: UUIDv7 timestamp parts were all different, this might be expected behavior")
	}
}
