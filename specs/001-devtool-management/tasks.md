# Tasks: devtool-management

**Input**: Design documents from `/specs/001-devtool-management/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create project structure per implementation plan (`cmd/`, `internal/`)
- [x] T002 Initialize Go modules (`go mod init devtool`)
- [x] T003 [P] Add primary dependencies (`spf13/cobra`, `zalando/go-keyring`, SQL drivers)
- [x] T004 Build main application entry point scaffolding in `cmd/devtool/main.go`
- [x] T005 Setup Cobra CLI root command in `internal/cli/root.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T006 Implement configuration loading (`.devtool.yml` via raw parsing or `viper`) in `internal/config/config.go`
- [x] T007 Implement secure credential storage wrappers (`Config.SetPassword`, `Config.GetPassword`) using `zalando/go-keyring` in `internal/config/secure.go`
- [x] T008 Add `devtool config set-credential` command to CLI in `internal/cli/config.go`
- [x] T009 Define core data models (`ToolType`, `ToolConfig`, `DatabaseConfig`, `StatusResult`) in `internal/checker/models.go`
- [x] T010 Implement the main concurrency checker orchestrator in `internal/checker/checker.go`

**Checkpoint**: Foundation ready - CLI can route commands and securely load configuration.

---

## Phase 3: User Story 1 - Environment Status Check (Priority: P1) 🎯 MVP

**Goal**: Check the status of all essential local developer tools from a single CLI command with summary output.

**Independent Test**: Can be tested by running the management command while manually configuring dummy services or mock checkers to verify the tabular output matches the aggregated state.

### Implementation for User Story 1

- [x] T011 [P] [US1] Add `devtool status` command to CLI in `internal/cli/status.go`
- [x] T012 [P] [US1] Implement basic Postgres connection ping using `lib/pq` or `pgx` in `internal/checker/postgres.go` (waits for timeout SC-003)
- [x] T013 [P] [US1] Implement basic MySQL connection ping using `go-sql-driver/mysql` in `internal/checker/mysql.go` (waits for timeout SC-003)
- [x] T014 [US1] Wire the DB checkers into the main `checker.go` orchestrator to run concurrently (FR-006)
- [x] T015 [US1] Implement table output formatting (tabwriter or similar) for the final `StatusReport` in `internal/cli/status.go`

**Checkpoint**: At this point, User Story 1 should be fully functional (can check DBs concurrently and print the table).

---

## Phase 4: User Story 2 - CLI-based Daemon Detection (Priority: P2)

**Goal**: Accurately detect if Docker and Telepresence are running by using their native CLI commands.

**Independent Test**: Can be fully tested by verifying command execution logic against known running/stopped daemon states locally.

### Implementation for User Story 2

- [x] T016 [P] [US2] Implement Docker CLI status parsing (`docker info`) via `os/exec` in `internal/checker/docker.go`
- [x] T017 [P] [US2] Implement Telepresence CLI status parsing (`telepresence status`) via `os/exec` in `internal/checker/telepresence.go`
- [x] T018 [US2] Wire the Daemon checkers into the main `checker.go` orchestrator
- [x] T019 [P] [US2] Add unit tests for parsing Docker CLI output in `tests/unit/docker_test.go`
- [x] T020 [P] [US2] Add unit tests for parsing Telepresence CLI output in `tests/unit/telepresence_test.go`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently and the tool is feature-complete for the 4 requested services.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T021 [P] Ensure all checkers respect the 3-second timeout gracefully (SC-003) and don't panic on missing CLI tools (FR-008). 
- [x] T022 Code cleanup, standardizing error formatting and logging messages.
- [x] T023 Run quickstart.md validation locally to verify instructions match implementation behavior.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel if mocked out, or sequentially in priority order (P1 → P2).
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### Parallel Opportunities

- Pings/Checkers logic (T012, T013, T016, T017) can be written independently of each other.
- Unit tests can be written in parallel with the implementation.
