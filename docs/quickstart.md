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

3. Run the CLI locally:
   ```bash
   go run ./cmd/devtool
   ```
   *This will launch the TUI.*

4. Test basic CLI commands (without compiling):
   ```bash
   # Status check
   go run ./cmd/devtool status

   # Add Postgres
   go run ./cmd/devtool add postgres --port 5433

   # Remove Postgres
   go run ./cmd/devtool remove postgres-5433
   ```

## Development Iterate Loop

When developing new features (e.g., adding a new database engine or tool type):

1. **Update the Manager**: Register the tool type in `internal/manager/` and define whether it's a Singleton or Port-Bound tool.
2. **Update the Checker**: Implement necessary ping functionality in `internal/checker/` to verify the new tool's live status.
3. **Update the TUI**: If the tool introduces unique workflows (similar to the Telepresence specific forms), build the necessary `.go` forms under `internal/tui/`. Ensure changes align with the current two-panel UI layout.
4. **Test Configurations**: Manually verify `~/.devtool.yml` and `~/.devtool/managed.json` reflect the expected state changes. Clear these out if you need a clean slate during testing:
   ```bash
   rm ~/.devtool.yml ~/.devtool/managed.json
   ```
