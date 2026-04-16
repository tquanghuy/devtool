# Antigravity Context: devtool

This directory contains `devtool`, a Go-based CLI and TUI application designed to manage local developer services and configurations such as Docker, Telepresence, Postgres, and MySQL.

## Project Stack & Tools
- **Language**: Go
- **UI Paradigm**: 
  - Standard CLI commands (`add`, `remove`, `status`, etc.).
  - TUI (Text User Interface) for interactive tabbed exploration, influenced by tools like lazydocker or k9s.
- **Dependencies**: Expect standard Go ecosystem CLI/TUI frameworks (e.g., Cobra, Bubble Tea/Charm).
- **Config Storage**: 
  - Managed tool state: `~/.devtool/managed.json`
  - App configuration: `~/.devtool.yml`
  - Credentials: Uses native keychain integration.

## AI Agent Guidelines
- **Go Best Practices**: Adhere to idiomatic Go. Use small, composable functions. Handle errors explicitly.
- **TUI/CLI Patterns**: When updating the interface or creating new pages, follow existing layout paradigms (e.g. the lazydocker-styled tabbed interface). Do not break non-interactive CLI usage (`--non-interactive` flag) when adding features.
- **Context Awareness**: The `devtool` can manage singleton tools (Docker, Telepresence) and port-bound tools (Postgres, MySQL). Ensure you consider disambiguation via ports when working with port-bound resources.
- **Documentation**: Provide concise comments and docstrings. Maintain documentation integrity.
- **Feature Workflows**: Reference the `specs/` directory for technical plans and specifications.

By following these guidelines, you ensure consistency, reliability, and an exceptional user experience in the `devtool` project.
