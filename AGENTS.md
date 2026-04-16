# Antigravity Context: devtool

`devtool` is a **centralized development tool management system**. Its goal is to provide a single interactive entry point for managing *every* tool a developer needs, whether they are global infrastructure components (Docker, Telepresence) or project-specific services (Postgres instances, local APIs, Redis).

## Project Goal
To serve as the "Single Point of Truth" for local developer tool lifecycles, using a configuration-driven approach that supports global and per-project tailoring.

## Project Stack & Tools
- **Language**: Go
- **UI Paradigm**: 
  - TUI-First: A keyboard-centric interactive terminal interface built with `tview`, influenced by `lazydocker` and `k9s`.
  - Minimal CLI: The CLI primarily serves as an entry point for the TUI.
- **Dependencies**: Uses `tview` and `tcell` for the TUI, `logrus` for logging, and `yaml.v3` for configuration.
- **Config Storage**: 
  - **Cascading Config**: Merges defaults, `~/.devtool.yml` (Global), and `./.devtool.yml` (Local).
  - **Managed State**: `~/.devtool/managed.json` tracks user-selected tool instances.

## AI Agent Guidelines
- **Configuration over Hardcoding**: Never hardcode new tools. Add them to `pkg/config/defaults.go` if they are core, or document how users can add them to their YAML configs.
- **TUI Patterns**: Follow the two-panel layout (`Tools` and `Connections`). Ensure new features are accessible via keyboard shortcuts or the action modals.
- **Context Awareness**: Always respect the cascading config. Tools defined in a local directory should seamlessly override or augment global tools.
- **Command Safety**: When implementing start/stop/check logic, ensure commands are executed via the unified `OSCommand` wrapper in `pkg/commands`.
- **Documentation**: Provide concise comments. Maintain this `AGENTS.md` and the `docs/` specs as the project evolves.

By following these guidelines, you ensure `devtool` remains the extensible, centralized hub it was designed to be.

