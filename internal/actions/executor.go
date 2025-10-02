package actions

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"anthodev/codory/internal/platform"
)

// Executor executes actions
type Executor struct {
	platformInfo platform.Info
}


// NewExecutor creates a new executor
func NewExecutor() *Executor {
	return &Executor{
		platformInfo: platform.DetectInfo(),
	}
}

// Execute executes an action
func (e *Executor) Execute(ctx context.Context, action *Action) (string, error) {
	switch action.Type {
	case ActionTypeFunction:
		if action.Handler == nil {
			return "", fmt.Errorf("no handler defined for action %s", action.ID)
		}
		return action.Handler(ctx)

	case ActionTypeCommand:
		return e.executeCommand(ctx, action)

	default:
		return "", fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// executeCommand executes a system command
func (e *Executor) executeCommand(ctx context.Context, action *Action) (string, error) {
	platformCmd, found := action.GetPlatformCommand(Platform(e.platformInfo.OS))
	if !found {
		if !found {
			// Try with PlatformAny
			platformCmd, found = action.GetPlatformCommand(PlatformAny)
			if !found {
				return "", fmt.Errorf("no command defined for platform %s", e.platformInfo.OS)
			}
		}
	}

	// Execute the command
	cmdStr := platformCmd.Command
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf("command failed: %w\n%s", err, output)
	}

	return string(output), nil
}

// GetPlatformInfo returns the platform information
func (e *Executor) GetPlatformInfo() platform.Info {
	return e.platformInfo
}
