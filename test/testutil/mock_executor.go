// Package testutil provides common testing utilities and mocks used throughout the application tests
package testutil

import (
	"fmt"

	"anthodev/codory/internal/actions"
)

// MockExecutor is a test double that simulates action execution without running actual commands.
// This allows for safe testing of actions without side effects or dependencies on external systems.
type MockExecutor struct {
	ExecuteFunc func(*actions.Action) (string, error)
}

// Execute implements the executor interface by calling the configured ExecuteFunc.
// If no ExecuteFunc is defined, it returns an error.
func (m *MockExecutor) Execute(action *actions.Action) (string, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(action)
	}
	return "", fmt.Errorf("no execute function defined")
}

// NewMockExecutor creates a new MockExecutor with the provided execute function.
func NewMockExecutor(executeFunc func(*actions.Action) (string, error)) *MockExecutor {
	return &MockExecutor{
		ExecuteFunc: executeFunc,
	}
}

// NewMockExecutorWithValidation creates a new MockExecutor that validates the action ID
// before executing the mock function.
func NewMockExecutorWithValidation(expectedID string, returnResult string, returnError error) *MockExecutor {
	return &MockExecutor{
		ExecuteFunc: func(action *actions.Action) (string, error) {
			if action.ID != expectedID {
				return "", fmt.Errorf("unexpected action ID: %s, expected: %s", action.ID, expectedID)
			}
			return returnResult, returnError
		},
	}
}
