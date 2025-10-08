// Package utils provides common utility functions used throughout the application
package utils

import "strings"

// Contains checks if a string contains a substring using a recursive approach.
// This is a basic implementation that checks for substring presence.
func Contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}

	// Check if substring matches at current position
	if s[:len(substr)] == substr {
		return true
	}

	// Check if substring matches at end position
	if len(s) > len(substr) && s[len(s)-len(substr):] == substr {
		return true
	}

	// Recursive check on remaining string
	return Contains(s[1:], substr)
}

// ContainsSubstring checks if a string contains a substring using an iterative approach.
// This is more efficient than the recursive Contains function for larger strings.
func ContainsSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ContainsString is a wrapper around strings.Contains for consistency with the custom functions.
// It uses the standard library implementation which is optimized for performance.
func ContainsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
