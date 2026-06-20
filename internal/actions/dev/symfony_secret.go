package dev

import (
	"anthodev/codory/internal/actions"
	"context"
	crypto_rand "crypto/rand"
	"fmt"

	"github.com/atotto/clipboard"
)

func NewSymfonySecretAction() *actions.Action {
	return &actions.Action{
		ID:          "symfony_secret",
		Name:        "Generate Symfony secret",
		Description: "Generate a new Symfony secret",
		Type:        actions.ActionTypeFunction,
		Handler:     generateSymfonySecret,
	}
}

func generateSymfonySecret(ctx context.Context) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, 64)

	_, err := crypto_rand.Read(bytes)

	if err != nil {
		return "", fmt.Errorf("generate Symfony secret: %w", err)
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	result := string(bytes)

	if err := clipboard.WriteAll(result); err != nil {
		return "", fmt.Errorf("copy Symfony secret to clipboard: %w", err)
	}

	return fmt.Sprintf("Generated Symfony secret: %s, copied to clipboard!", result), nil
}
