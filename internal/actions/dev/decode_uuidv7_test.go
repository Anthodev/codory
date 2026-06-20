package dev

import (
	"context"
	"testing"
	"time"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
)

func TestNewDecodeUUIDv7Action(t *testing.T) {
	action := DecodeUUIDv7Action()

	if action.ID != "decode_uuidv7" {
		t.Errorf("Expected action ID to be 'decode_uuidv7', got '%s'", action.ID)
	}

	if action.Name != "Decode UUIDv7" {
		t.Errorf("Expected action name to be 'Decode UUIDv7', got '%s'", action.Name)
	}

	if action.Description != "Decode an UUIDv7" {
		t.Errorf("Expected action description to be 'Decode an UUIDv7', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeFunction {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeFunction, action.Type)
	}

	if action.Handler == nil {
		t.Error("Expected action handler to be set")
	}

	if len(action.Arguments) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(action.Arguments))
	}

	if action.Arguments[0].Name != "uuidv7" {
		t.Errorf("Expected argument name to be 'uuidv7', got '%s'", action.Arguments[0].Name)
	}

	if !action.Arguments[0].Required {
		t.Error("Expected uuidv7 argument to be required")
	}
}

func TestDecodeUUIDv7(t *testing.T) {
	validUUIDv7 := "01944bf0-7c2d-7cc4-ba43-7b6f2f4a7c1d"
	ctx := context.WithValue(context.Background(), actions.ArgsContextKey, []string{validUUIDv7})

	result, err := decodeUUIDv7(ctx)
	if err != nil {
		t.Errorf("Expected no error for valid UUIDv7, got %v", err)
	}

	expectedPrefix := "Datetime decoded for " + validUUIDv7 + ":"
	if !utils.Contains(result, expectedPrefix) {
		t.Errorf("Expected result to start with '%s', got '%s'", expectedPrefix, result)
	}

	// The result should contain a valid RFC3339 timestamp
	// Extract the timestamp part
	timestampStart := len(expectedPrefix)
	if len(result) <= timestampStart {
		t.Error("Result too short to contain timestamp")
		return
	}

	timestampStr := result[timestampStart:]
	timestampStr = trimSpace(timestampStr)

	// Try to parse the timestamp
	_, err = parseTimeRFC3339(timestampStr)
	if err != nil {
		t.Errorf("Expected valid RFC3339 timestamp, got '%s' (error: %v)", timestampStr, err)
	}
}

func TestDecodeUUIDv7_WithoutHyphens(t *testing.T) {
	validUUIDv7 := "01944bf07c2d7cc4ba437b6f2f4a7c1d"
	ctx := context.WithValue(context.Background(), actions.ArgsContextKey, []string{validUUIDv7})

	result, err := decodeUUIDv7(ctx)
	if err != nil {
		t.Errorf("Expected no error for valid UUIDv7 without hyphens, got %v", err)
	}

	expectedPrefix := "Datetime decoded for " + validUUIDv7 + ":"
	if !utils.Contains(result, expectedPrefix) {
		t.Errorf("Expected result to start with '%s', got '%s'", expectedPrefix, result)
	}
}

func TestDecodeUUIDv7_InvalidUUID(t *testing.T) {
	testCases := []struct {
		name     string
		uuid     string
		expected string
	}{
		{
			name:     "empty UUID",
			uuid:     "",
			expected: "invalid UUIDv7",
		},
		{
			name:     "too short UUID",
			uuid:     "01944bf0-7c2d-7cc4-ba43-7b6f2f4a7c",
			expected: "invalid UUIDv7",
		},
		{
			name:     "too long UUID",
			uuid:     "01944bf0-7c2d-7cc4-ba43-7b6f2f4a7c1d-extra",
			expected: "invalid UUIDv7",
		},
		{
			name:     "invalid characters in timestamp part",
			uuid:     "01944gfo-7c2d-7cc4-ba43-7b6f2f4a7c1d", // 'g' is not valid hex in timestamp part
			expected: "invalid UUIDv7",
		},
		{
			name:     "completely invalid format",
			uuid:     "not-a-uuid-at-all",
			expected: "invalid UUIDv7",
		},
		{
			name:     "UUID is not version 7",
			uuid:     "01944bf0-7c2d-4cc4-ba43-7b6f2f4a7c1d",
			expected: "invalid UUIDv7",
		},
		{
			name:     "UUID has non-RFC variant",
			uuid:     "01944bf0-7c2d-7cc4-7a43-7b6f2f4a7c1d",
			expected: "invalid UUIDv7",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), actions.ArgsContextKey, []string{tc.uuid})
			_, err := decodeUUIDv7(ctx)
			if err == nil {
				t.Errorf("Expected error for %s, got none", tc.name)
			}
			if err != nil && !utils.Contains(err.Error(), tc.expected) {
				t.Errorf("Expected error containing '%s', got '%v'", tc.expected, err)
			}
		})
	}
}

func TestDecodeUUIDv7_MissingArgument(t *testing.T) {
	ctx := context.WithValue(context.Background(), actions.ArgsContextKey, []string{})

	_, err := decodeUUIDv7(ctx)
	if err == nil {
		t.Fatal("Expected error when no arguments provided, got nil")
	}
	if !utils.Contains(err.Error(), "invalid UUIDv7") {
		t.Errorf("Expected invalid UUIDv7 error, got %v", err)
	}
}

func TestUUID7stringToAtom(t *testing.T) {
	testCases := []struct {
		name        string
		uuid        string
		shouldError bool
		contains    string // substring that should be in the result
	}{
		{
			name:        "valid UUIDv7 with hyphens",
			uuid:        "01944bf0-7c2d-7cc4-ba43-7b6f2f4a7c1d",
			shouldError: false,
			contains:    "T", // RFC3339 timestamp should contain 'T'
		},
		{
			name:        "valid UUIDv7 without hyphens",
			uuid:        "01944bf07c2d7cc4ba437b6f2f4a7c1d",
			shouldError: false,
			contains:    "T", // RFC3339 timestamp should contain 'T'
		},
		{
			name:        "too short UUID",
			uuid:        "01944bf0-7c2d-7cc4-ba43-7b6f2f4a7c",
			shouldError: true,
		},
		{
			name:        "too long UUID",
			uuid:        "01944bf0-7c2d-7cc4-ba43-7b6f2f4a7c1d-extra",
			shouldError: true,
		},
		{
			name:        "invalid hex characters in timestamp part",
			uuid:        "01944gfo-7c2d-7cc4-ba43-7b6f2f4a7c1d", // 'g' in timestamp part
			shouldError: true,
		},
		{
			name:        "empty string",
			uuid:        "",
			shouldError: true,
		},
		{
			name:        "wrong UUID version",
			uuid:        "01944bf0-7c2d-4cc4-ba43-7b6f2f4a7c1d",
			shouldError: true,
		},
		{
			name:        "wrong UUID variant",
			uuid:        "01944bf0-7c2d-7cc4-7a43-7b6f2f4a7c1d",
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := uuid7stringToAtom(tc.uuid)

			if tc.shouldError {
				if err == nil {
					t.Errorf("Expected error for %s, got none", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for %s, got %v", tc.name, err)
				}
				if tc.contains != "" && !utils.Contains(result, tc.contains) {
					t.Errorf("Expected result to contain '%s', got '%s'", tc.contains, result)
				}
			}
		})
	}
}

func TestDecodeUUIDv7ActionRegistration(t *testing.T) {
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

	// Register the DecodeUUIDv7 action
	decodeAction := DecodeUUIDv7Action()
	err = registry.RegisterAction("dev", decodeAction)
	if err != nil {
		t.Fatalf("Failed to register DecodeUUIDv7 action: %v", err)
	}

	// Verify the action was registered
	devCat, exists := registry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	// Count actions with ID decode_uuidv7
	decodeActionCount := 0
	for _, action := range devCat.Actions {
		if action.ID == "decode_uuidv7" {
			decodeActionCount++
		}
	}

	if decodeActionCount != 1 {
		t.Errorf("Expected 1 DecodeUUIDv7 action in dev category, got %d", decodeActionCount)
	}
}

func TestDecodeUUIDv7_Consistency(t *testing.T) {
	// Test that decoding the same UUID multiple times gives consistent results
	validUUIDv7 := "01944bf0-7c2d-7cc4-ba43-7b6f2f4a7c1d"
	ctx := context.WithValue(context.Background(), actions.ArgsContextKey, []string{validUUIDv7})

	result1, err1 := decodeUUIDv7(ctx)
	if err1 != nil {
		t.Errorf("Expected no error on first call, got %v", err1)
	}

	result2, err2 := decodeUUIDv7(ctx)
	if err2 != nil {
		t.Errorf("Expected no error on second call, got %v", err2)
	}

	if result1 != result2 {
		t.Errorf("Expected consistent results, got '%s' and '%s'", result1, result2)
	}
}

// Helper function for trimming whitespace
func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	for start < end && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}

func parseTimeRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
