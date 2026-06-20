---
name: add-codory-app
description: Add a new Codory registry app/tool action with platform install commands and definition tests. Use when adding apps, tools, package actions, or registry entries to Codory.
---

# Add Codory App

Add one Codory action for installing a developer/system app.

## Gather input

Ask one focused question if missing:

- app/tool name, or
- official project/package link.

## Inspect Codory first

Before choosing platforms, installers, categories, or imports:

1. Read `internal/actions/action.go` for current `Platform`, `PackageSource`, `ActionType`, and `PlatformCommand` fields.
2. Read existing action files under `internal/actions/*/install_*.go` for command patterns, fallback usage, `CheckCommand`, `Interactive`, and naming.
3. Read `internal/actions/*/actions.go` for current category IDs and registration style.
4. Check `main.go` imports before adding a new action package.

Treat project constants and nearby action patterns as workflow truth. Do not rely on stale platform or installer lists in this skill.

## Research upstream

1. Web-search the app name or open provided link.
2. Prefer official sources: homepage, README, install docs, package registry pages.
3. If multiple plausible projects match, stop and ask user to choose.
4. Match upstream-supported install methods to current Codory `Platform` and `PackageSource` constants.
5. Prefer project-specific platform commands over broad fallbacks when upstream commands differ.
6. Use generic/fallback platforms only when upstream documents them or existing Codory patterns support them.
7. Find `CheckCommand` from the installed binary name.
8. If platform support, package ID, installer source, or binary name is unclear, ask user before coding.

## Pick category

Choose from existing categories registered in `internal/actions/*/actions.go`. Use nearby actions as examples for category intent.

If category is ambiguous, ask user.

## References

Use these reference templates before coding:

- [Action template](references/action-template.md)
- [Action test template](references/action-test-template.md)

## Implement

1. Read the reference templates above.
2. Add `internal/actions/<category>/install_<app_snake>.go`.
3. Register action in `internal/actions/<category>/actions.go` with existing category ID and style.
4. If adding a new action package, blank-import it from `main.go`.
5. Add `internal/actions/<category>/install_<app_snake>_test.go`.
6. Keep tests as definition checks only. Never execute installer commands.
7. Run:

```bash
gofmt -w internal/actions/<category>/install_<app_snake>.go internal/actions/<category>/install_<app_snake>_test.go
CODORY_TEST=1 GO_TEST=1 go test ./...
```

## Command rules

- Use only `Platform` and `PackageSource` constants found in `internal/actions/action.go`.
- Copy command shape from existing action files when it matches current project patterns.
- Set `Interactive: true` for `sudo`, prompts, curl-pipe installers, GUI installers, or command flows needing user input.
- Prefer platform-specific commands over broad fallback when commands differ.
- Delete unsupported platform placeholders.
- Add unsupported-platform checks for intentionally unsupported targets.

## Report

Return:

- selected project URL
- category
- platforms added/skipped
- files changed
- test result
