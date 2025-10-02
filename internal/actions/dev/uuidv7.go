package dev

import (
	"anthodev/codory/internal/actions"
	"context"
	"fmt"
	"log"

	"github.com/atotto/clipboard"
	"github.com/google/uuid"
)

func NewUUIDv7Action() *actions.Action {
	return &actions.Action{
		ID:          "uuidv7",
		Name:        "Generate UUIDv7",
		Description: "Generate a new UUIDv7",
		Type:        actions.ActionTypeFunction,
		Handler:     generateUUIDv7,
	}
}

func generateUUIDv7(ctx context.Context) (string, error) {
	id, err := uuid.NewV7()

	if err != nil {
		log.Fatal(err)
	}

	if err := clipboard.WriteAll(id.String()); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("Generated UUIDv7: %s, copied to clipboard!", id.String()), nil
}
