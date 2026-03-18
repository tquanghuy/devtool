# Task Breakdown: Interactive CLI

## Implementation Strategy

We will follow an MVP-first approach, starting with the configuration storage layer, then building the TUI components story-by-story, and finally wiring it up to the CLI entry point.

## Phase 1: Setup
- [x] T001 Initialize project structure for `internal/devtools` and `internal/tui`
- [x] T002 Add `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/bubbles` to `go.mod`

## Phase 2: Foundational (Config Storage)
- [x] T003 Define `DevtoolProfile` and `DevtoolsConfig` structs in `internal/devtools/config.go`
- [x] T004 Implement `Load()` config logic in `internal/devtools/config_io.go`
- [x] T005 Implement `Save()` config logic with atomic write in `internal/devtools/config_io.go`
- [x] T006 Implement `Add()` profile logic with validation in `internal/devtools/config_io.go`
- [x] T007 Implement `Remove()` profile logic in `internal/devtools/config_io.go`
- [x] T008 [P] Implement unit tests for config I/O in `tests/unit/devtools_config_test.go`

## Phase 3: User Story 1 - First-Time Empty State Experience [US1]
- [x] T009 [US1] Implement root TUI model with empty state message in `internal/tui/model.go`
- [x] T010 [US1] Create basic list view using `bubbles/list` in `internal/tui/model.go`
- [x] T011 [US1] Implement "Add First Tool" prompt in empty state view

## Phase 4: User Story 2 - Adding a New Devtool Interactively [US2]
- [x] T012 [P] [US2] Implement multi-step add form in `internal/tui/add_form.go`
- [x] T013 [US2] Wire `a` keybinding to launch add form in `internal/tui/model.go`
- [x] T014 [US2] Implement form submission logic calling `devtools.Add` and `Save`
- [x] T015 [US2] Add validation feedback for required fields in add form

## Phase 5: User Story 3 - Loading and Displaying Configured Tools [US3]
- [x] T016 [US3] Update root model to load config on startup in `internal/tui/model.go`
- [x] T017 [US3] Map `devtools.DevtoolProfile` to `list.Item` for rendering
- [x] T018 [US3] Implement arrow-key navigation for the tools list

## Phase 6: User Story 4 - Removing a Configured Devtool [US4]
- [x] T019 [P] [US4] Implement delete confirmation dialog in `internal/tui/confirm.go`
- [x] T020 [US4] Wire `d` keybinding to launch confirmation in `internal/tui/model.go`
- [x] T021 [US4] Implement deletion logic calling `devtools.Remove` and `Save`

## Phase 7: Tool Execution & CLI Integration
- [x] T022 Implement tool execution logic using `tea.ExecProcess` in `internal/tui/model.go`
- [x] T023 Modify `internal/cli/root.go` to launch TUI when no arguments are provided
- [x] T024 [P] Implement state-machine tests for TUI models in `tests/unit/tui_model_test.go`

## Phase 8: Polish & Cross-cutting
- [x] T025 Verify error handling for invalid/corrupted `devtools.yml`
- [x] T026 Final manual verification of all user scenarios in the terminal
- [x] T027 Ensure existing CLI subcommands still function correctly

## Dependencies
- Phase 2 (Foundational) blocks all User Story phases.
- US1 is the entry point for US2.
- US2 must be completed to verify US3 (with actual data).
- US4 can be implemented in parallel with US3 after US1/US2 are baseline.

## Parallel Execution Examples
- T008 (Config Tests) can run in parallel with T009 (TUI Root Model).
- T012 (Add Form) and T019 (Confirm Dialog) can be built in parallel as they are separate models.
- T024 (TUI Tests) can run as soon as models are defined.
