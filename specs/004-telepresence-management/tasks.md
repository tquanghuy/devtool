# Tasks: Telepresence Management

**Input**: Design documents from `/specs/004-telepresence-management/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure for Telepresence integration.

- [ ] T001 Initialize Telepresence package/file structure in `internal/tools/telepresence/telepresence.go` (if creating new package) or `internal/tui/telepresence.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented.

- [ ] T002 Implement helper to check binary existence using `exec.LookPath("telepresence")` in `internal/tui/telepresence.go`
- [ ] T003 [P] Implement parser for `telepresence status` command output in `internal/tui/telepresence.go`

**Checkpoint**: Foundation ready - user story implementation can now begin.

---

## Phase 3: User Story 1 - Add and View Telepresence Tool (Priority: P1) 🎯 MVP

**Goal**: As a developer, I want to add Telepresence to my list of managed development tools and view its current connection status.

**Independent Test**: Can be fully tested by adding the tool to the devtool and verifying that it appears in the list with a "connected" or "not connected" status.

### Implementation for User Story 1

- [ ] T004 [P] [US1] Create `DevtoolItem` initialization for Telepresence in `internal/tui/delegate.go`
- [ ] T005 [P] [US1] Implement status-checking `tea.Cmd` for Telepresence in `internal/tui/telepresence.go`
- [ ] T006 [US1] Extend UI view to display Telepresence status alongside the item in `internal/tui/delegate.go`
- [ ] T007 [US1] Add logic to periodically refresh Telepresence status or fetch it on view load

**Checkpoint**: At this point, User Story 1 should be fully functional; Telepresence is visible and its status is correct.

---

## Phase 4: User Story 2 - Manage Telepresence Connection (Priority: P1)

**Goal**: As a developer, I want to be able to connect, disconnect, or restart my Telepresence session directly from the devtool interface.

**Independent Test**: Can be tested by selecting Telepresence in the devtool list and executing the connect, disconnect, and restart actions, then verifying the state changes in the UI and via terminal.

### Implementation for User Story 2

- [ ] T008 [P] [US2] Implement asynchronous `tea.Cmd` executor for `telepresence connect` in `internal/tui/telepresence.go`
- [ ] T009 [P] [US2] Implement asynchronous `tea.Cmd` executor for `telepresence quit` in `internal/tui/telepresence.go`
- [ ] T010 [US2] Handle UI state transitions (e.g., showing a connecting spinner) during async operations in `internal/tui/delegate.go`
- [ ] T011 [US2] Add Connect, Disconnect, Restart, and Remove actions to the Telepresence context menu in `internal/tui/delegate.go`
- [ ] T012 [US2] Implement fallback UI prompt for missing Telepresence binary listing brew/npm installation options

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories.

- [ ] T013 Verify asynchronous operations never block the main TUI event loop
- [ ] T014 Add error logging for failed `exec.Cmd` calls
- [ ] T015 Verify the feature manually using steps from `quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on all user stories being complete

### Parallel Opportunities

- T002 and T003 can be worked on in parallel.
- T004 and T005 can be worked on in parallel.
- T008 and T009 can be worked on in parallel.

## Implementation Strategy

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Telepresence shows up in UI
3. Add User Story 2 → Test actions independently → Complete MVP
