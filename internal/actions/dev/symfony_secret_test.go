package dev

import (
	"context"
	"testing"

	"anthodev/codory/internal/actions"
	"anthodev/codory/pkg/utils"
)

func TestNewSymfonySecretAction(t *testing.T) {
	action := NewSymfonySecretAction()

	// Test action properties
	if action.ID != "symfony_secret" {
		t.Errorf("Expected action ID to be 'symfony_secret', got '%s'", action.ID)
	}

	if action.Name != "Generate Symfony secret" {
		t.Errorf("Expected action name to be 'Generate Symfony secret', got '%s'", action.Name)
	}

	if action.Description != "Generate a new Symfony secret" {
		t.Errorf("Expected action description to be 'Generate a new Symfony secret', got '%s'", action.Description)
	}

	if action.Type != actions.ActionTypeFunction {
		t.Errorf("Expected action type to be '%s', got '%s'", actions.ActionTypeFunction, action.Type)
	}

	if action.Handler == nil {
		t.Error("Expected action handler to be set")
	}
}

func TestGenerateSymfonySecret(t *testing.T) {
	ctx := context.Background()
	result, err := generateSymfonySecret(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Test that the result contains "Generated Symfony secret:" and "copied to clipboard!"
	if !utils.Contains(result, "Generated Symfony secret:") {
		t.Errorf("Expected result to contain 'Generated Symfony secret:', got '%s'", result)
	}

	if !utils.Contains(result, "copied to clipboard!") {
		t.Errorf("Expected result to contain 'copied to clipboard!', got '%s'", result)
	}

	// Extract the secret from the result
	secretStart := findSubstringIndex(result, "Generated Symfony secret: ")
	if secretStart == -1 {
		t.Error("Could not find secret in result")
		return
	}
	secretStart += len("Generated Symfony secret: ")

	secretEnd := findSubstringIndex(result[secretStart:], ", copied to clipboard!")
	if secretEnd == -1 {
		t.Error("Could not find end of secret in result")
		return
	}

	secret := result[secretStart : secretStart+secretEnd]

	// Test that the secret is 64 characters long
	if len(secret) != 64 {
		t.Errorf("Expected secret length to be 64, got %d", len(secret))
	}

	// Test that the secret contains only valid characters (letters and digits)
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range secret {
		found := false
		for _, validChar := range validChars {
			if char == validChar {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Secret contains invalid character '%c'", char)
		}
	}
}

func TestGenerateSymfonySecret_UniqueGeneration(t *testing.T) {
	// Test that multiple calls generate different secrets
	ctx := context.Background()

	result1, err1 := generateSymfonySecret(ctx)
	if err1 != nil {
		t.Errorf("Expected no error on first call, got %v", err1)
	}

	result2, err2 := generateSymfonySecret(ctx)
	if err2 != nil {
		t.Errorf("Expected no error on second call, got %v", err2)
	}

	// Extract just the secret parts for comparison
	secret1Start := findSubstringIndex(result1, "Generated Symfony secret: ") + len("Generated Symfony secret: ")
	secret1End := findSubstringIndex(result1[secret1Start:], ", copied to clipboard!")
	secret1 := result1[secret1Start : secret1Start+secret1End]

	secret2Start := findSubstringIndex(result2, "Generated Symfony secret: ") + len("Generated Symfony secret: ")
	secret2End := findSubstringIndex(result2[secret2Start:], ", copied to clipboard!")
	secret2 := result2[secret2Start : secret2Start+secret2End]

	if secret1 == secret2 {
		t.Error("Expected different secrets on multiple calls")
	}
}

func TestGenerateSymfonySecret_CryptographicRandomness(t *testing.T) {
	// Test that we get different secrets on multiple calls (indicating good randomness)
	ctx := context.Background()
	secrets := make(map[string]bool)
	numTests := 10

	for i := 0; i < numTests; i++ {
		result, err := generateSymfonySecret(ctx)
		if err != nil {
			t.Errorf("Expected no error on call %d, got %v", i, err)
			continue
		}

		// Extract the secret
		secretStart := findSubstringIndex(result, "Generated Symfony secret: ") + len("Generated Symfony secret: ")
		secretEnd := findSubstringIndex(result[secretStart:], ", copied to clipboard!")
		secret := result[secretStart : secretStart+secretEnd]

		if secrets[secret] {
			t.Errorf("Duplicate secret generated on call %d: %s", i, secret)
		}
		secrets[secret] = true
	}

	// We should have generated unique secrets
	if len(secrets) != numTests {
		t.Errorf("Expected %d unique secrets, got %d", numTests, len(secrets))
	}
}

func TestGenerateSymfonySecret_ContextHandling(t *testing.T) {
	// Test that the function works with different context types
	ctx := context.Background()
	result, err := generateSymfonySecret(ctx)

	if err != nil {
		t.Errorf("Expected no error with background context, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result with background context")
	}

	// Test with a context that has values (should still work)
	ctxWithValue := context.WithValue(context.Background(), "test", "value")
	result2, err2 := generateSymfonySecret(ctxWithValue)

	if err2 != nil {
		t.Errorf("Expected no error with context with value, got %v", err2)
	}

	if result2 == "" {
		t.Error("Expected non-empty result with context with value")
	}
}

func TestSymfonySecretActionRegistration(t *testing.T) {
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

	// Register the SymfonySecret action
	symfonySecretAction := NewSymfonySecretAction()
	err = registry.RegisterAction("dev", symfonySecretAction)
	if err != nil {
		t.Fatalf("Failed to register SymfonySecret action: %v", err)
	}

	// Verify the action was registered
	devCat, exists := registry.GetCategory("dev")
	if !exists {
		t.Fatal("Dev category not found")
	}

	if len(devCat.Actions) != 1 {
		t.Errorf("Expected 1 action in dev category, got %d", len(devCat.Actions))
	}

	if devCat.Actions[0].ID != "symfony_secret" {
		t.Errorf("Expected action ID to be 'symfony_secret', got '%s'", devCat.Actions[0].ID)
	}
}

func TestGenerateSymfonySecret_SecretCharacterDistribution(t *testing.T) {
	// Test that the generated secrets have a good distribution of characters
	ctx := context.Background()
	result, err := generateSymfonySecret(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Extract the secret
	secretStart := findSubstringIndex(result, "Generated Symfony secret: ") + len("Generated Symfony secret: ")
	secretEnd := findSubstringIndex(result[secretStart:], ", copied to clipboard!")
	secret := result[secretStart : secretStart+secretEnd]

	// Count character types
	lowercaseCount := 0
	uppercaseCount := 0
	digitCount := 0

	for _, char := range secret {
		switch {
		case char >= 'a' && char <= 'z':
			lowercaseCount++
		case char >= 'A' && char <= 'Z':
			uppercaseCount++
		case char >= '0' && char <= '9':
			digitCount++
		}
	}

	// Basic sanity checks - we should have all character types represented
	if lowercaseCount == 0 {
		t.Error("Expected at least one lowercase letter in secret")
	}
	if uppercaseCount == 0 {
		t.Error("Expected at least one uppercase letter in secret")
	}
	if digitCount == 0 {
		t.Error("Expected at least one digit in secret")
	}

	// The secret should have a reasonable distribution (not all one type)
	totalChars := len(secret)
	if lowercaseCount == totalChars {
		t.Error("Secret should not be all lowercase letters")
	}
	if uppercaseCount == totalChars {
		t.Error("Secret should not be all uppercase letters")
	}
	if digitCount == totalChars {
		t.Error("Secret should not be all digits")
	}
}
