# Implementation Plan: Interactive CLI

**Branch**: `003-interactive-cli` | **Date**: 2026-03-15 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/003-interactive-cli/spec.md`

## Summary

Transform the `devtool` CLI from a pure subcommand-based tool into one that launches a Bubble Tea interactive TUI when run with no arguments. Users can add, remove, and launch developer tools from an arrow-key navigable menu. Tool profiles are persisted in `~/.devtools.yml` (YAML). Existing CLI subcommands remain intact for scripted usage.

## Technical Context

**Language/Version**: Go 1.25  
**Primary Dependencies**: `github.com/spf13/cobra` (existing), `gopkg.in/yaml.v3` (existing), `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/bubbles` (new)  
**Storage**: `~/.devtools.yml` (YAML, via `yaml.v3`)  
**Testing**: `github.com/stretchr/testify` (existing), `go test ./...`  
**Target Platform**: macOS / Linux terminals  
**Project Type**: CLI tool  
**Performance Goals**: Startup < 200ms, menu navigation instant  
**Constraints**: Minimize new dependencies; only Charm TUI ecosystem added  
**Scale/Scope**: Single user, local tool list (10s of devtools expected)

## Constitution Check

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Clean Code | ✅ Pass | New `internal/tui` and `internal/devtools` packages keep concerns separate |
| II. Simple CLI Interface | ✅ Pass | No-args entry = interactive TUI; all subcommands remain available |
| III. Easy to Use CLI Tool | ✅ Pass | Empty-state guidance, arrow-key navigation, form-based add |
| IV. Centralized Developer Tools | ✅ Pass | Single devtool binary as gateway to configured tools |
| Testing | ✅ Pass | Unit tests for all config I/O logic; TUI model state tests |
| Minimal Dependencies | ✅ Pass | Only adding charmbracelet/bubbletea + bubbles |

## Project Structure

### Documentation (this feature)

```text
specs/003-interactive-cli/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/cli.md     # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code Changes

```text
internal/
├── devtools/              [NEW] Devtool profile config I/O
│   ├── config.go          [NEW] DevtoolProfile + DevtoolsConfig structs + YAML tags
│   └── config_io.go       [NEW] Load(), Save(), Add(), Remove() helpers
├── tui/                   [NEW] Bubble Tea TUI models
│   ├── model.go           [NEW] Root model (list view, keybindings, state machine)
│   ├── add_form.go        [NEW] Multi-step add tool form
│   └── confirm.go         [NEW] Delete confirmation dialog
└── cli/
    └── root.go            [MODIFY] Add RunE: launch TUI when no subcommand given

cmd/devtool/main.go        [no change]

tests/unit/
├── devtools_config_test.go  [NEW] Unit tests for Load/Save/Add/Remove
└── tui_model_test.go        [NEW] Pure state-machine tests for TUI model

go.mod / go.sum            [MODIFY] Add bubbletea + bubbles dependencies
```

## Implementation Phases

### Phase 1: Devtool Profile Storage (`internal/devtools`)

1. Define `DevtoolProfile` and `DevtoolsConfig` structs with yaml tags.
2. Implement `Load()` — reads `~/.devtools.yml`, returns empty config if file not found, error if malformed.
3. Implement `Save(cfg)` — atomic write via temp file + rename (same pattern as `config_io.go`).
4. Implement `Add(cfg, profile)` — validates non-empty name+executable, uniqueness check, appends.
5. Implement `Remove(cfg, name)` — removes by name, error if not found.

### Phase 2: TUI Models (`internal/tui`)

1. **`model.go`**: Root Bubble Tea model with a `bubbles/list` for the devtool list. Handles keybindings: `↑`/`↓` (navigation), `Enter` (launch), `a` (add form), `d` (remove confirm), `q`/`Ctrl+C` (quit). On startup: load config into list items. Empty-state message when list is empty.
2. **`add_form.go`**: Sequential `bubbles/textinput` form with 3 steps (name → executable → args). On submit: calls `devtools.Add` + `devtools.Save`, returns new item to root model. On cancel (`Esc`): returns to main menu.
3. **`confirm.go`**: Simple yes/no confirmation model for deletion. `y` → calls `devtools.Remove` + `devtools.Save`, removes from list. Any other key → cancel, return to main menu.
4. **Tool execution**: `bubbletea.Cmd` that suspends the program (`tea.ExecProcess`), runs the tool via `os/exec`, and resumes on exit.

### Phase 3: Wire Up CLI Entry Point

Modify `internal/cli/root.go`:
- Add `PersistentPreRunE` or check `cobra.NoArgs` on the root command.
- When `devtool` is called with no args, run `tui.Start()`.
- Add `rootCmd.Args = cobra.NoArgs` guard with redirect to TUI instead of error.

### Phase 4: Tests

- `tests/unit/devtools_config_test.go`: Table-driven tests for `Load`, `Save`, `Add`, `Remove` using `t.TempDir()` + `t.Setenv("HOME", …)`.
- `tests/unit/tui_model_test.go`: Instantiate Bubble Tea models in test, send `tea.KeyMsg` messages, assert model state (selected item, form values, list length) without actual rendering.

## Verification Plan

### Automated Tests

```sh
# Run all unit tests (existing + new)
go test ./tests/unit/... -v
```

Expected: all existing tests pass; new devtools_config and tui_model tests pass.

```sh
# Verify build compiles cleanly
go build ./cmd/devtool/...
```

### Manual Verification

1. **Empty state**: Run `./devtool` with no `~/.devtools.yml` — should show empty-state message with prompt to add a tool.
2. **Add tool**: Press `a`, fill in name/executable/args, confirm — tool should appear in list and `~/.devtools.yml` should contain the new entry (`cat ~/.devtools.yml`).
3. **Remove tool**: Navigate to the tool, press `d`, confirm with `y` — tool should disappear from list and from `~/.devtools.yml`.
4. **Cancel remove**: Press `d`, then press any key other than `y` — tool should remain.
5. **Launch tool**: Navigate to a valid tool (e.g., `/bin/sh`), press `Enter` — shell should launch; on exit, devtool menu should resume.
6. **Existing subcommands unaffected**: `devtool add docker`, `devtool remove docker`, `devtool status` should all still work.
