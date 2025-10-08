package testutil

import (
	"errors"
	"testing"

	"anthodev/codory/internal/actions"
)

func TestMockExecutor_Execute(t *testing.T) {
	tests := []struct {
		name        string
		executeFunc func(*actions.Action) (string, error)
		action      *actions.Action
		wantResult  string
		wantErr     bool
		errMsg      string
	}{
		{
			name: "successful execution",
			executeFunc: func(action *actions.Action) (string, error) {
				return "success", nil
			},
			action:     &actions.Action{ID: "test-action"},
			wantResult: "success",
			wantErr:    false,
		},
		{
			name: "execution with error",
			executeFunc: func(action *actions.Action) (string, error) {
				return "", errors.New("execution failed")
			},
			action:     &actions.Action{ID: "test-action"},
			wantResult: "",
			wantErr:    true,
			errMsg:     "execution failed",
		},
		{
			name:        "no execute function defined",
			executeFunc: nil,
			action:      &actions.Action{ID: "test-action"},
			wantResult:  "",
			wantErr:     true,
			errMsg:      "no execute function defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := &MockExecutor{
				ExecuteFunc: tt.executeFunc,
			}

			result, err := mockExecutor.Execute(tt.action)

			if result != tt.wantResult {
				t.Errorf("Execute() result = %v, want %v", result, tt.wantResult)
			}

			if tt.wantErr && err == nil {
				t.Error("Execute() expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Execute() unexpected error = %v", err)
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Execute() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestNewMockExecutor(t *testing.T) {
	executeFunc := func(action *actions.Action) (string, error) {
		return "test result", nil
	}

	mockExecutor := NewMockExecutor(executeFunc)

	if mockExecutor.ExecuteFunc == nil {
		t.Error("NewMockExecutor() ExecuteFunc is nil")
	}

	action := &actions.Action{ID: "test"}
	result, err := mockExecutor.Execute(action)
	if err != nil {
		t.Errorf("NewMockExecutor() unexpected error = %v", err)
	}
	if result != "test result" {
		t.Errorf("NewMockExecutor() result = %v, want %v", result, "test result")
	}
}

func TestNewMockExecutorWithValidation(t *testing.T) {
	tests := []struct {
		name         string
		expectedID   string
		actionID     string
		returnResult string
		returnError  error
		wantResult   string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid action ID",
			expectedID:   "test-action",
			actionID:     "test-action",
			returnResult: "success",
			returnError:  nil,
			wantResult:   "success",
			wantErr:      false,
		},
		{
			name:         "invalid action ID",
			expectedID:   "test-action",
			actionID:     "wrong-action",
			returnResult: "",
			returnError:  nil,
			wantResult:   "",
			wantErr:      true,
			errMsg:       "unexpected action ID: wrong-action, expected: test-action",
		},
		{
			name:         "with error return",
			expectedID:   "test-action",
			actionID:     "test-action",
			returnResult: "",
			returnError:  errors.New("custom error"),
			wantResult:   "",
			wantErr:      true,
			errMsg:       "custom error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := NewMockExecutorWithValidation(tt.expectedID, tt.returnResult, tt.returnError)
			action := &actions.Action{ID: tt.actionID}

			result, err := mockExecutor.Execute(action)

			if result != tt.wantResult {
				t.Errorf("Execute() result = %v, want %v", result, tt.wantResult)
			}

			if tt.wantErr && err == nil {
				t.Error("Execute() expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Execute() unexpected error = %v", err)
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Execute() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
