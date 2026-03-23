---
description: End-to-end feature development using the spec-kit skills
---
This workflow guides you through developing a new feature from initial requirements to complete implementation using the spec-kit methodology.

1. **Specification Phase**: Use the `speckit-specify` skill to create a detailed feature specification (`spec.md`) from natural language descriptions or user requests.
2. **Clarification Phase** (Optional): If the requirements in the specification are underspecified or ambiguous, use the `speckit-clarify` skill to resolve them through a structured questioning workflow. Use the `speckit-checklist` skill to generate custom quality checklists for validating the completeness.
3. **Planning Phase**: Once the specification is clear, use the `speckit-plan` skill to generate a technical implementation plan (`plan.md`) describing the architecture and tech stack.
4. **Task Breakdown Phase**: Use the `speckit-tasks` skill to break down the technical plan into an actionable and dependency-aware task list (`tasks.md`). Optionally, use `speckit-taskstoissues` to convert these into project management issues.
5. **Consistency Check Phase**: Use the `speckit-analyze` skill to perform cross-artifact consistency analysis across `spec.md`, `plan.md`, and `tasks.md` to identify gaps and discrepancies before implementation starts.
6. **Implementation Phase**: Use the `speckit-implement` skill to execute all tasks in the breakdown to build the feature following a Test-Driven Development (TDD) approach where applicable.
