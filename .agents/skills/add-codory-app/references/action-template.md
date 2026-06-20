# Codory Action Template

Create one constructor in the selected category package.

Path:

```text
internal/actions/<category>/install_<app_snake>.go
```

Before using this template, inspect:

- `internal/actions/action.go` for current `Platform`, `PackageSource`, `ActionType`, and `PlatformCommand` definitions.
- `internal/actions/*/install_*.go` for current command, fallback, `CheckCommand`, and `Interactive` patterns.
- `internal/actions/*/actions.go` for category IDs and registration style.

Template. Platform entries below are illustrative only; replace them with platforms and package sources supported by current project constants and upstream docs:

```go
package <category_package>

import "anthodev/codory/internal/actions"

func Install<AppPascal>() *actions.Action {
	return &actions.Action{
		ID:          "install_<app_snake>",
		Name:        "Install <AppName>",
		Description: "Install <short description>",
		Type:        actions.ActionTypeCommand,
		PlatformCommands: map[actions.Platform]actions.PlatformCommand{
			actions.PlatformArch: {
				Command:       "sudo pacman -S <package>",
				PackageSource: actions.PackageSourceOfficial,
				CheckCommand:  "<binary>",
				Interactive:   true,
			},
			actions.PlatformMacOS: {
				Command:       "brew install <formula>",
				PackageSource: actions.PackageSourceBrew,
				CheckCommand:  "<binary>",
			},
		},
		// Add visibility fields only when current project patterns need them.
	}
}
```

Register in existing category package using the category ID and style from that package:

```go
func init() {
	registry := actions.GlobalRegistry()

	// existing category setup...

	registry.RegisterAction("<category_id>", Install<AppPascal>())
}
```

If creating a new action package, add blank import in `main.go`:

```go
_ "anthodev/codory/internal/actions/<category>"
```

Notes:

- Delete unsupported platforms from the map.
- Use fallback platforms only when upstream docs or existing Codory patterns support them.
- Use package source constants exactly as defined in `internal/actions/action.go`.
- Use `Interactive: true` for `sudo`, prompts, scripts, or installer UIs.
