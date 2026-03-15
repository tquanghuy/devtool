# Tasks: devtool-add-remove

**Input**: Design documents from `/specs/002-devtool-add-remove/`
**Prerequisites**: `plan.md` ✅ · `spec.md` ✅ · `data-model.md` ✅ · `contracts/cli.md` ✅ · `research.md` ✅ · `quickstart.md` ✅

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no shared dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Exact file paths are included in all task descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify project scaffolding and create the new packages required by this feature.

- [x] T001 Verify existing packages (`internal/checker`, `internal/config`, `internal/cli`) build cleanly with `go build ./...`
- [x] T002 Create `internal/manager/` package directory with placeholder `manager.go` (package declaration only)
- [x] T003 [P] Create `tests/unit/` directory with `.gitkeep` if it does not exist

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core data structures and config operations that ALL user stories depend on.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T004 Define `ToolType` enum and `ToolDefinition` struct in `internal/checker/tool_definition.go` (fields: `Name string`, `Type ToolType`, `DefaultPort int`)
- [x] T005 Define the supported tool registry (static map/slice of predefined tools) in `internal/checker/registry.go`
- [x] T006 [P] Define `ManagedInstance` struct in `internal/config/managed_instance.go` (fields: `ToolName string`, `Identifier string`, `CreatedAt time.Time`)
- [x] T007 [P] Define `ManagedConfig` struct in `internal/config/managed_config.go` (field: `Instances []ManagedInstance`) with JSON serialization tags
- [x] T008 Implement `ManagedConfig.AddInstance(instance *ManagedInstance) error` in `internal/config/managed_config.go` (validates singleton uniqueness and identifier uniqueness)
- [x] T009 Implement `ManagedConfig.RemoveInstance(identifier string) error` in `internal/config/managed_config.go`
- [x] T010 Implement config file read/write helpers (`LoadConfig() (*ManagedConfig, error)`, `SaveConfig(cfg *ManagedConfig) error`) in `internal/config/config_io.go`
- [x] T011 Implement `LookupTool(name string) (*ToolDefinition, bool)` helper in `internal/checker/registry.go`

**Checkpoint**: Foundation ready — `ManagedConfig` can be loaded, modified, and persisted. User story implementation can now begin.

---

## Phase 3: User Story 1 — Add a Singleton Tool (Priority: P1) 🎯 MVP

**Goal**: Allow a user to add a singleton tool (e.g., Docker, Telepresence) to the managed list exactly once.

**Independent Test**: Run `devtool add docker` on an empty config → tool appears in list. Run it again → error shown. Validate no network calls occur.

### Implementation for User Story 1

- [x] T012 [US1] Implement `add` command skeleton with Cobra in `internal/cli/add.go` (register `devtool add <tool-name>` with `--port` and `--non-interactive` flags; set Cobra `Short`, `Long`, and per-flag `Usage` strings for `--help` output) *(C2, C3)*
- [x] T013 [US1] Implement singleton-add logic in `internal/manager/add.go`: look up `ToolDefinition` via `LookupTool`; on unsupported tool **print the full list of supported tool names from `internal/checker/registry.go` then return error** (FR-003, *G1*); reject duplicate singleton (FR-005); create `ManagedInstance{Identifier: ToolName}`, call `config.AddInstance` and `config.SaveConfig`
- [x] T014 [US1] Wire `internal/cli/add.go` to call `internal/manager/add.go` and print success/error messages per `contracts/cli.md`
- [x] T015 [US1] Register `add` command in `internal/cli/root.go` (or equivalent command registration file)
- [x] T015a [P] [US1] Write unit tests for singleton-add in `tests/unit/add_singleton_test.go`: (a) add succeeds on empty config, (b) second add returns singleton-duplicate error, (c) unsupported tool returns error with supported-tools list *(C1)*

**Checkpoint**: `devtool add docker` works end-to-end. Singleton guard is enforced. Unit tests pass. Ready to demo/test independently.

---

## Phase 4: User Story 2 — Add a Port-bound Tool (Priority: P1)

**Goal**: Allow a user to add a port-bound tool (e.g., PostgreSQL, MySQL) multiple times with distinct port/identifier assignments.

**Independent Test**: Run `devtool add postgres` for the first time → succeeds on default port. Run `devtool add postgres` a second time without `--port` → **fails** with error asking user to specify `--port`. Run `devtool add postgres --port 5433` → succeeds.

### Implementation for User Story 2

- [x] T016 [US2] Extend `internal/manager/add.go` with port-bound path: attempt `ToolDefinition.DefaultPort` as identifier `<ToolName>-<Port>`; if that identifier already exists in config → **fail immediately with error "Port <N> already in use. Please specify a different port with --port."** (Option B, no auto-assign); if `--port` is given explicitly and already exists → same fail behavior *(I1 fix)*
- [x] T017 [US2] Wire `--non-interactive` flag read in `internal/cli/add.go`: pass flag value to manager; when set and port conflict occurs, skip any prompt and return error directly per `contracts/cli.md` *(G3)*
- [x] T018 [US2] Handle the case where all preferred ports are occupied in `internal/manager/add.go`: fail with instructive error "No available port found. Please specify a port explicitly with --port." (Option B)
- [x] T019 [P] [US2] Write unit tests for port-bound add in `tests/unit/add_portbound_test.go`: (a) first add succeeds on default port, (b) second add without `--port` fails with conflict error, (c) add with explicit free `--port` succeeds, (d) `--non-interactive` flag skips prompt and fails on conflict *(C1)*

**Checkpoint**: `devtool add postgres` and `devtool add postgres --port 5433` both work. Port conflicts fail with clear errors. `--non-interactive` behaves correctly.

---

## Phase 5: User Story 3 — Remove a Tool (Priority: P1)

**Goal**: Allow a user to remove a managed tool instance from the list, with disambiguation for multiple port-bound instances and a safety prompt for running tools.

**Independent Test**: Add a tool, then remove it — it no longer appears in the list. Remove a non-managed tool → error. Remove one of two postgres instances → only the targeted instance is removed.

### Implementation for User Story 3

- [x] T020 [US3] Implement `remove` command skeleton with Cobra in `internal/cli/remove.go` (register `devtool remove <tool-name>|<instance-id>` with `--force` and `--non-interactive` flags; set Cobra `Short`, `Long`, and per-flag `Usage` strings for `--help` output) *(C2, C3)*
- [x] T021 [US3] Implement remove logic in `internal/manager/remove.go`:
  - Query config for all instances matching tool name
  - 0 found → print "not managed" error (FR-002)
  - 1 found → proceed to running-check step
  - >1 found → if interactive, display numbered selection list of `Identifier`s; if `--non-interactive`, **immediately fail** with error per `contracts/cli.md` ("Multiple instances found. Please specify exact instance ID or run interactively.") *(G2)*
- [x] T021a [US3] Wire `--non-interactive` flag read in `internal/cli/remove.go`: pass flag value to manager so multi-instance disambiguation and running-tool checks respect it *(G2)*
- [x] T022 [US3] Implement running-instance check in `internal/manager/remove.go`: if tool process is detected as running, prompt user to confirm stop (interactive) or fail (non-interactive / no `--force`); with `--force`, terminate process and proceed
- [x] T023 [US3] Call `config.RemoveInstance` and `config.SaveConfig` on confirmed removal; print success/error messages per `contracts/cli.md`
- [x] T024 [US3] Wire `internal/cli/remove.go` to call `internal/manager/remove.go`
- [x] T025 [US3] Register `remove` command in `internal/cli/root.go` (or equivalent command registration file)
- [x] T025a [P] [US3] Write unit tests for remove in `tests/unit/remove_test.go`: (a) remove managed singleton succeeds, (b) remove non-managed tool returns error, (c) multi-instance interactive selection removes only selected instance, (d) `--non-interactive` fails when multiple instances exist *(C1)*

**Checkpoint**: Full lifecycle works — add then remove cleanly. Multi-instance disambiguation verified. Running-tool guard works. Unit tests pass.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Error message polish, edge-case hardening, documentation, and constitution compliance validation.

- [x] T026 [P] Validate all user-facing error messages match the exact strings in `contracts/cli.md` and `quickstart.md`
- [x] T027 Verify edge cases from `spec.md`: unsupported tool name → informative error listing supported tools; double-add singleton → single clear error
- [x] T028 [P] Run `go vet ./...` and fix any issues across `internal/manager/`, `internal/cli/`, `internal/config/`, `internal/checker/`
- [x] T029 [P] Update `README.md` with `devtool add` and `devtool remove` usage examples drawn from `quickstart.md`
- [x] T030 Run full `quickstart.md` walkthrough manually and verify all example outputs match exactly: (a) `devtool add postgres` → "Successfully added postgres to managed list.", (b) `devtool add docker` twice → second shows singleton error, (c) `devtool remove postgres` with 2 instances → selection prompt appears, (d) `devtool remove postgres --force` → success without prompt *(A2 fix)*
- [x] T031 [P] Validate `--help` output for both commands: run `devtool add --help` and `devtool remove --help`, confirm `Short`, `Long`, and all flag descriptions are present and accurate per `contracts/cli.md` *(C2)*

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — **blocks all user stories**
- **Phase 3–5 (User Stories)**: All depend on Phase 2 completion
  - US1 (singleton add) should be completed before US2 (port-bound add), as US2 extends the same `add.go` manager
  - US3 (remove) depends on Phase 2 only; can proceed independently of US2
- **Phase 6 (Polish)**: Depends on all user stories being complete

### User Story Dependencies

| Story | Depends On         | Can run in parallel with |
|-------|--------------------|--------------------------|
| US1   | Phase 2            | US3                      |
| US2   | Phase 2 + US1 done | US3                      |
| US3   | Phase 2            | US1                      |

### Within Each User Story

- Data structures (manager logic) before CLI wiring
- CLI wiring before smoke testing

### Parallel Opportunities

- T006 and T007 (struct definitions) can run in parallel
- T028 (go vet) and T029 (README) can run in parallel
- US1 and US3 can be developed in parallel by two developers

---

## Parallel Example: Foundational Phase

```
Parallel (different files, no conflict):
  T006  internal/config/managed_instance.go
  T007  internal/config/managed_config.go

Sequential (depends on T006/T007):
  T008  ManagedConfig.AddInstance  →  T009  ManagedConfig.RemoveInstance
  T010  config_io.go
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (**CRITICAL — blocks all stories**)
3. Complete Phase 3: User Story 1 (singleton add)
4. **STOP and VALIDATE**: `devtool add docker` works; second add is rejected
5. Demo/merge if ready

### Incremental Delivery

1. Setup + Foundational → foundation ready
2. + US1 → `devtool add <singleton>` works (MVP!)
3. + US2 → `devtool add <port-bound>` works (multiple instances)
4. + US3 → `devtool remove` works (full lifecycle)
5. + Phase 6 → polished, documented, production-ready

---

## Notes

- `[P]` tasks operate on different files and have no in-flight dependencies
- `[Story]` label maps each task to its user story for traceability
- Commit after each phase checkpoint to enable easy rollback
- Unit tests are required by the project constitution (§Dev Guidelines) — `tests/unit/` tasks (T015a, T019, T025a) are mandatory, not optional
- Avoid cross-story file conflicts: `add.go` is extended by US2 but must be complete for US1 first
- **Constitution compliance checklist before PR**: ✅ unit tests present · ✅ `--help` output verified · ✅ README updated · ✅ `go vet` clean
