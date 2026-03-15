# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Add functionality to `devtool` to manage local development tools by adding them to or removing them from the user's managed list. The system supports singleton tools (max 1 instance, like Docker) and port-bound tools (multiple instances allowed, like PostgreSQL). The `add` command will act as a configuration manager without strictly verifying host executables, while the `remove` command will interactively handle running instances and disambiguate multiple port-bound instances.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.21+
**Primary Dependencies**: `spf13/cobra` (CLI framework), internal config/checker packages
**Storage**: Local configuration file — **JSON format**, default path `~/.devtool/managed.json` (via internal config package) *(U1 fix)*
**Testing**: Go `testing` package, `testify`
**Target Platform**: Developer Workstations (macOS/Linux/Windows)
**Project Type**: CLI Tool
**Performance Goals**: Command execution < 3 seconds
**Constraints**: Zero external network requests required for add/remove logic
**Scale/Scope**: Local developer environment, up to ~50 managed tools

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **Clean Code**: Design keeps logic separated (CLI parsing vs. configuration management vs. tool logic).
- [x] **Simple CLI Interface**: Commands will be standard `devtool add <tool>` and `devtool remove <tool>`.
- [x] **Easy to Use CLI Tool**: Interactive prompts for ambiguous removals and helpful error messages for unsupported tools.
- [x] **Centralized Developer Tools**: Integrates seamlessly into the existing `devtool` paradigm.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
```text
# Internal library packages
internal/
├── checker/       # (Existing) Tool definitions and status checking
├── config/        # (Existing) Managed list configuration read/write
├── cli/           # (Existing) Cobra commands, adding add.go and remove.go
├── manager/       # (New or Existing) Core logic for adding/removing instances

# Tests
tests/
└── unit/          # Unit tests for new add/remove logic
```

**Structure Decision**: The project is an existing Go CLI application. We will follow the established standard Go project layout, placing new commands in `internal/cli` and core business logic in `internal/manager` or `internal/config`. No new top-level directories are needed.

> **Fill ONLY if Constitution Check has violations that must be justified**

No constitution violations found. Code complexity is kept strictly necessary for the CLI tool specification.
