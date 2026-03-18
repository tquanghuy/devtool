# Quick Start: Interactive CLI (003-interactive-cli)

## Prerequisites

- Go 1.25+
- The `devtool` binary built and in your `$PATH`, or run via `go run ./cmd/devtool`

## Launch the Interactive UI

```sh
devtool
```

If `~/.devtools.yml` does not exist yet, the tool starts with an empty list.

## Add Your First Tool

1. Press `a` in the main menu.
2. Enter a display name (e.g., `My Postgres`).
3. Enter the executable path (e.g., `/usr/local/bin/psql`).
4. Enter default arguments (e.g., `-h localhost -d mydb`), or leave blank.

The tool is saved to `~/.devtools.yml` and immediately appears in the menu.

## Run a Tool

1. Use `↑`/`↓` to navigate to the tool.
2. Press `Enter` — the tool's own interactive interface launches in the same terminal.
3. When you exit the tool, the devtool menu resumes.

## Remove a Tool

1. Navigate to the tool with `↑`/`↓`.
2. Press `d`.
3. Confirm with `y`. The tool is removed from `~/.devtools.yml`.

## Quit

Press `q` or `Ctrl+C`.

## Configuration File Location

`~/.devtools.yml` — human-editable YAML. Restart `devtool` after manual edits.
