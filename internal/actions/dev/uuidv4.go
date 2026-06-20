package dev

import (
	"context"
	"fmt"

	"anthodev/codory/internal/actions"

	"github.com/atotto/clipboard"
	"github.com/google/uuid"
)

func NewUUIDv4Action() *actions.Action {
	return &actions.Action{
		ID:          "uuidv4",
		Name:        "Generate UUIDv4",
		Description: "Generate a new UUIDv4",
		Type:        actions.ActionTypeFunction,
		Handler:     generateUUIDv4,
	}
}

func generateUUIDv4(ctx context.Context) (string, error) {
	id := uuid.New()

	if err := clipboard.WriteAll(id.String()); err != nil {
		return "", fmt.Errorf("copy UUIDv4 to clipboard: %w", err)
	}

	return fmt.Sprintf("Generated UUID: %s, copied to clipboard!", id.String()), nil
}
