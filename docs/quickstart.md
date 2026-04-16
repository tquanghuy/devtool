# Quickstart Guide

This guide will help you set up and run `devtool` locally for development and testing.

## Prerequisites

- **Go**: Version 1.20 or later.
- Ensure ports `5432` (Postgres default) and `3306` (MySQL default) are unblocked if you plan on testing those tools locally.
- **Docker** (optional, but required if you want to test docker-managed instances).

## Running Locally

1. Clone the repository and navigate to the project root:
   ```bash
   cd devtool
   ```

2. Download and verify modules:
   ```bash
   go mod tidy
   ```

3. Build and Run the application:
   ```bash
   make build
   ./devtool
   ```
   *This will launch the interactive TUI.*

## Development Workflow

`devtool` is designed to be **configuration-driven**. To add or modify tools:

### 1. Unified Configuration
Manage your development environment via `.devtool.yml`. You can define new tool types here:

```yaml
tools:
  my-custom-service:
    name: "my-service"
    kind: "singleton"
    start_cmd: "npm run dev"
    stop_cmd: "pkill -f 'npm run dev'"
    check_cmd: "curl -s localhost:3000"
```

### 2. Cascading Configs
- **Global**: `~/.devtool.yml` (Applies everywhere)
- **Local**: `./.devtool.yml` (Applies to the current project only)

### 3. Modifying the Core
If you need to change the application's internal behavior:
- **`pkg/gui`**: Update the Tview-based TUI panels and modals.
- **`pkg/config`**: Modify how configuration and managed states are handled.
- **`pkg/commands`**: Update the command execution or status checking logic.

## Resetting State
If you need a clean slate during testing, you can remove the managed tool state:
```bash
rm ~/.devtool/managed.json
```
To reset your global configuration:
```bash
rm ~/.devtool.yml
```

