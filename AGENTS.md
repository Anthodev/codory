# Codory Project Context

Purpose: concise orientation for coding agents. Focus on project goal, structure, and safe development workflow.

## Goal

Codory is a cross-platform Go TUI CLI for installing and managing developer/system tools. It uses Bubble Tea for navigation, platform-aware action registration, and command/function actions for tasks such as installing tools, generating UUIDs/secrets, and checking package managers.

## Tech Stack

- Language: Go (`go 1.25.1` in `go.mod`)
- Module: `anthodev/codory`
- TUI: Bubble Tea (`github.com/charmbracelet/bubbletea`)
- UI components/styles: Bubbles + Lip Gloss
- Utilities: clipboard support and UUID generation
- License: MIT

## Runtime Flow

1. `main.go` blank-imports action packages so their `init()` functions register actions.
2. `main.go` starts Bubble Tea with `ui.NewModel()`.
3. `internal/ui` reads registered categories/actions from `actions.GlobalRegistry()`.
4. UI displays only actions visible for detected platform.
5. Selected action is executed by `actions.Executor` as either:
   - Go function handler
   - platform-specific shell command

## Project Structure

| Path | Role |
|---|---|
| `main.go` | CLI entry point; starts Bubble Tea program; imports action packages |
| `internal/ui/` | Bubble Tea model, navigation states, rendering, styles |
| `internal/actions/` | Core action/category types, global registry, executor |
| `internal/actions/dev/` | Developer utilities and dev-tool actions |
| `internal/actions/package_managers/` | Yay, Homebrew, Winget actions |
| `internal/actions/shells/` | Shell, Zsh, Oh My Zsh, Zsh plugin actions |
| `internal/actions/terminal/` | Terminal emulator/multiplexer actions |
| `internal/actions/tools/` | Small CLI/system tool actions |
| `internal/platform/` | OS/distro/package-manager detection and installers |
| `pkg/utils/` | Shared utility functions |
| `test/testutil/` | Test helpers and mocks |
| `.github/workflows/test.yml` | CI test/coverage workflow |

## Architecture Notes

- Architecture style: modular monolithic CLI with registry-driven action plugins.
- Categories/actions are registered during package initialization.
- Platform support is modeled with `actions.Platform` and per-action `PlatformCommands`.
- Debian/Arch can fall back to generic Linux commands when no distro-specific command exists.
- No HTTP server, API routes, database, repositories, or background jobs.

## Adding or Changing Actions

- Add action constructor in relevant `internal/actions/...` package.
- Register it from that package `init()` function.
- If adding a new action package, blank-import it in `main.go`.
- Command actions should define `PlatformCommands`, `PackageSource`, and `CheckCommand` when possible.
- Mark commands requiring user interaction/sudo prompts with `Interactive: true`.
- Do not run real install commands in tests; mock or validate action definitions.

## Commands

Run from repo root.

| Task | Command |
|---|---|
| Run app | `go run .` |
| Build | `go build -o codory .` |
| Format | `gofmt -w <files>` |
| Test | `CODORY_TEST=1 GO_TEST=1 go test ./...` |
| CI-style test | `CODORY_TEST=1 GO_TEST=1 xvfb-run -a go test -v ./...` |
| Coverage | `CODORY_TEST=1 GO_TEST=1 go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out` |

## Testing

- Test runner: Go `testing` via `go test ./...`.
- CI uses Ubuntu, Go 1.25, Xvfb, clipboard utilities, and Codecov upload.
- Unit tests exist for most action/platform/helper packages.
- `internal/ui` currently has no tests.
- No E2E test suite/config found.
- Safe test mode uses `CODORY_TEST=1` and `GO_TEST=1`.

## Tooling

- Go modules: `go.mod`, `go.sum`.
- Formatting: `gofmt`; `.editorconfig` present.
- CI: GitHub Actions in `.github/workflows/test.yml`.
- No Makefile, golangci-lint config, staticcheck config, Dockerfile, or release config found.

## Known Issues To Watch

- README appears incomplete/truncated.
- README mentions package-manager auto-install/prompting, but UI flow does not fully wire that behavior.
- Some command strings/check commands look suspicious (`apt get`, `brew installl`, VS Code check command, macOS bat check command).
- Interactive commands run through `sh -c`, which may be problematic on Windows.
- Some function handlers use `log.Fatal`, which exits instead of returning UI errors.
- `.rules` describes a stricter Clean Architecture layout than the actual registry-based CLI layout.

## Safe Change Checklist

1. Modify only relevant files.
2. Run `gofmt` on edited Go files.
3. Run `CODORY_TEST=1 GO_TEST=1 go test ./...`.
4. Avoid real system install commands in tests.
5. Update tests when changing action behavior or command definitions.
