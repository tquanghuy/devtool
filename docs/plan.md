# Project Plan: devtool

## Overview

`devtool` is a Go-based CLI and TUI application designed to seamlessly manage local developer services and configurations. It orchestrates tools like Docker, Telepresence, Postgres, and MySQL, simplifying the process of checking statuses, managing multiple database instances, and controlling proxy tunnels.

## Architecture

The project is structured following standard Go practices, roughly divided between standard CLI invocation and an interactive Terminal User Interface (TUI).

- **`cmd/devtool/main.go`**: The entry point.
- **`internal/cli`**: Contains Cobra-based command definitions (`add`, `remove`, `status`). It handles both standard flag-based execution and routes to the TUI if no subcommands are provided.
- **`internal/tui`**: The interactive interface built with Bubble Tea. Features a two-panel layout (formerly tabbed) for visual, interactive management of Connections and Tools.
- **`internal/manager`**: The core business logic for adding, removing, and changing the state of tools. It bridges the CLI/TUI layer and the configuration state.
- **`internal/checker`**: Responsible for determining if a service is actively running and accessible.
- **`internal/config`**: Manages the reading, writing, and validation of persistent states (`~/.devtool.yml` for database hosts/ports and `~/.devtool/managed.json` for managed tool lifecycles).

## Core Workflows

### 1. Initialization and Status
When invoked via `devtool status`, the application iterates over the `ManagedConfig` list, pings each service using the `checker`, and outputs their live status. When invoked as `devtool` (no arguments), the TUI loads, presenting real-time status in an interactive format.

### 2. Managing Connections (Port-Bound Tools)
Users can add multiple instances of port-bound tools (like Postgres or MySQL) by specifying differe ports. The manager handles disambiguation, prompts for port mapping, and guarantees uniqueness across the `Identifier` field (e.g., `postgres-5432`).

### 3. Managing Singleton Tools
Tools like Docker and Telepresence are Singletons. The system enforces that only one instance of these tools can exist in the managed configuration at a time. The TUI provides dedicated forms or flows, like connecting/disconnecting Telepresence directly from the interface.

## AI Agent Integration

The project has established guidelines in `AGENTS.md` for AI agent contributions, emphasizing:
- Preserving the TUI/CLI paradigms.
- Disambiguating port-bound resources.
- Following idiomatic Go and maintaining documentation integrity.
