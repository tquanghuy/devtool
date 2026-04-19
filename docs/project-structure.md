# Project Structure

This document outlines the directory structure and architectural patterns used in `devtool`.

## Directory Overview

```text
.
├── cmd/                # Application entry points
│   └── devtool/        # Main CLI application
│       └── main.go     # Entry point: initializes and runs the app
├── docs/               # Documentation (Roadmap, Design Specs, etc.)
│   ├── contracts/      # Interface specifications (e.g., CLI)
│   └── ...             # Feature-specific docs
├── pkg/                # Core application logic (Library code)
│   ├── app/            # Application lifecycle management
│   ├── commands/       # OS command execution wrappers & tool logic
│   ├── config/         # Cascading configuration system (YAML/JSON)
│   ├── gui/            # TUI implementation using tview
│   └── utils/          # Shared utility functions (currently empty)
├── AGENTS.md           # Instructions and context for AI coding assistants
├── Makefile            # Build and development task automation
├── go.mod/sum          # Go module dependency management
└── README.md           # Project introduction and usage
```

## Key Files

- `cmd/devtool/main.go`: The main entry point.
- `pkg/config/defaults.go`: Hardcoded default configurations.
- `pkg/config/config.go`: Cascading YAML configuration logic.
- `pkg/config/managed.go`: JSON state management for active tools.
- `pkg/gui/gui.go`: Central TUI state and initialization.
- `pkg/commands/os_command.go`: Shell command execution wrapper.
- `AGENTS.md`: Guidelines for AI agents (crucial for maintaining project patterns).

## Architectural Layers

### 1. Entry Point (`cmd/`)
The `cmd/devtool/main.go` file is intentionally minimal. It parses flags and hands off control to `pkg/app`. This pattern keeps the application logic separate from the command-line interface.

### 2. Application Layer (`pkg/app`)
The `app` package orchestrates the startup and shutdown of the application. It:
- Initializes the configuration.
- Sets up the GUI.
- Manages the main application loop.

### 3. Configuration Layer (`pkg/config`)
`devtool` uses a cascading configuration system:
- **Defaults**: Hardcoded in `pkg/config/defaults.go`.
- **Global**: User-defined in `~/.devtool.yml`.
- **Local**: Project-specific in `./.devtool.yml`.
- **Managed State**: Persisted state in `~/.devtool/managed.json` (e.g., which tools are currently active).

### 4. User Interface Layer (`pkg/gui`)
Built with `tview` and `tcell`, this layer follows a panel-based architecture:
- **`gui.go`**: Main GUI struct and initialization.
- **`layout.go`**: Defines the two-panel layout and overall TUI structure.
- **`keybindings.go`**: Centralizes keyboard shortcuts.
- **`*_panel.go`**: Individual components (Tools, Connections, Details, Resources).

### 5. Execution Layer (`pkg/commands`)
This layer provides a unified way to interact with the host operating system. The `OSCommand` wrapper ensures:
- Commands are logged.
- Error handling is consistent.
- Execution is safe across different environments.

## Development Patterns

- **TUI-First**: All management features should be accessible via the terminal interface.
- **Configuration over Hardcoding**: New tools should be added via configuration files rather than being baked into the source code whenever possible.
- **State Persistence**: User selections and tool statuses are tracked in `managed.json` to remain consistent across sessions.
