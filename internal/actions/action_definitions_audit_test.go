package actions_test

import (
	"strings"
	"testing"

	"anthodev/codory/internal/actions"
	_ "anthodev/codory/internal/actions/dev"
	_ "anthodev/codory/internal/actions/dev/code_editors"
	_ "anthodev/codory/internal/actions/dev/docker"
	_ "anthodev/codory/internal/actions/package_managers"
	_ "anthodev/codory/internal/actions/shells"
	_ "anthodev/codory/internal/actions/shells/zsh"
	_ "anthodev/codory/internal/actions/shells/zsh/plugins"
	_ "anthodev/codory/internal/actions/terminal"
	_ "anthodev/codory/internal/actions/tools"
)

func TestRegisteredCommandActionsPassAudit(t *testing.T) {
	auditCategory(t, actions.GlobalRegistry().GetRoot())
}

func TestCommandActionAuditCatchesInvalidDefinitions(t *testing.T) {
	tests := []struct {
		name   string
		action *actions.Action
		want   string
	}{
		{
			name: "empty command",
			action: auditAction(actions.PlatformArch, actions.PlatformCommand{
				Command:       " ",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "tool",
			}),
			want: "empty command",
		},
		{
			name: "empty check",
			action: auditAction(actions.PlatformArch, actions.PlatformCommand{
				Command:       "pacman -S tool",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  " ",
			}),
			want: "empty check command",
		},
		{
			name: "package manager check",
			action: auditAction(actions.PlatformMacOS, actions.PlatformCommand{
				Command:       "brew install tool",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "brew",
			}),
			want: "package-manager-only check command",
		},
		{
			name: "sudo not interactive",
			action: auditAction(actions.PlatformDebian, actions.PlatformCommand{
				Command:       "sudo apt install tool",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "tool",
			}),
			want: "sudo command must be interactive",
		},
		{
			name: "source mismatch",
			action: auditAction(actions.PlatformMacOS, actions.PlatformCommand{
				Command:       "apt install tool",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "tool",
			}),
			want: "brew source command must start with brew",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := auditCommandAction(tt.action); !strings.Contains(got, tt.want) {
				t.Fatalf("auditCommandAction() = %q, want substring %q", got, tt.want)
			}
		})
	}
}

func auditCategory(t *testing.T, category *actions.Category) {
	t.Helper()
	for _, action := range category.Actions {
		if action.Type != actions.ActionTypeCommand {
			continue
		}
		if got := auditCommandAction(action); got != "" {
			t.Error(got)
		}
	}
	for _, child := range category.SubCategories {
		auditCategory(t, child)
	}
}

func auditAction(platform actions.Platform, cmd actions.PlatformCommand) *actions.Action {
	return &actions.Action{
		ID:          "audit_tool",
		Name:        "Audit tool",
		Description: "Audit tool",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			platform: cmd,
		},
	}
}

func auditCommandAction(action *actions.Action) string {
	if strings.TrimSpace(action.ID) == "" {
		return "action has empty ID"
	}
	if strings.TrimSpace(action.Name) == "" {
		return action.ID + ": empty name"
	}
	if strings.TrimSpace(action.Description) == "" {
		return action.ID + ": empty description"
	}
	if len(action.PlatformCommands) == 0 {
		return action.ID + ": empty platform commands"
	}

	for platform, cmd := range action.PlatformCommands {
		where := action.ID + ":" + string(platform) + ": "
		command := strings.TrimSpace(cmd.Command)
		check := strings.TrimSpace(cmd.CheckCommand)
		if command == "" {
			return where + "empty command"
		}
		if check == "" {
			return where + "empty check command"
		}
		if isPackageManagerOnlyCheck(check) {
			return where + "package-manager-only check command " + check
		}
		if cmd.PackageSource == "" {
			return where + "empty package source"
		}
		if strings.Contains(command, "sudo") && !cmd.Interactive {
			return where + "sudo command must be interactive"
		}
		if got := auditPackageSourceCommand(cmd.PackageSource, command); got != "" {
			return where + got
		}
	}
	return ""
}

func isPackageManagerOnlyCheck(check string) bool {
	switch check {
	case "brew", "winget", "winget.exe", "yay", "apt", "apt-get", "pacman":
		return true
	default:
		return false
	}
}

func auditPackageSourceCommand(source actions.PackageSource, command string) string {
	switch source {
	case actions.PackageSourceBrew:
		if !strings.HasPrefix(command, "brew") {
			return "brew source command must start with brew"
		}
	case actions.PackageSourceWinget:
		if !strings.HasPrefix(command, "winget") && !strings.HasPrefix(command, "winget.exe") {
			return "winget source command must start with winget or winget.exe"
		}
	case actions.PackageSourceAUR:
		if !strings.HasPrefix(command, "yay") {
			return "AUR source command must start with yay"
		}
	}
	return ""
}
