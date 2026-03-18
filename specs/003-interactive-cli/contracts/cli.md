# CLI Contract: Interactive Mode (003-interactive-cli)

## Entry Point

When `devtool` is invoked with **no arguments**, the interactive TUI mode launches.
Existing subcommands (`devtool add`, `devtool remove`, `devtool status`) remain available and are unaffected.

## Interactive Main Menu

The root TUI presents a **navigable list** of configured devtools.

```
┌─ devtool ─────────────────────────────────────────┐
│                                                    │
│  Your devtools  (↑/↓ to navigate, Enter to open) │
│                                                    │
│  ▸ Postgres REPL                                  │
│    Docker TUI                                      │
│                                                    │
│  [a] Add tool   [d] Delete tool   [q] Quit        │
└────────────────────────────────────────────────────┘
```

### Empty State

```
┌─ devtool ─────────────────────────────────────────┐
│                                                    │
│  No devtools configured yet.                       │
│  Press [a] to add your first tool.                 │
│                                                    │
│  [a] Add tool   [q] Quit                          │
└────────────────────────────────────────────────────┘
```

## Keyboard Bindings

| Key    | Action |
|--------|--------|
| `↑`/`↓` | Navigate the list |
| `Enter` | Open the selected tool's interactive interface |
| `a`    | Open the add-tool form |
| `d`    | Open the remove confirmation for the selected tool |
| `q` / `Ctrl+C` | Quit |

## Add Tool Form

Sequential text-input prompts collected one at a time:

1. **Name** (required, must be unique): `Enter a display name for this tool:`
2. **Executable** (required): `Enter the path to the executable:`
3. **Default arguments** (optional): `Enter default arguments (leave blank to skip):`

On completion, the profile is saved to `~/.devtools.yml` and the menu is refreshed.
On cancel (`Esc`), return to the main menu without saving.

## Remove Confirmation

```
  Remove "Postgres REPL"? [y/N]:
```

- `y` / `Y`: Removes entry from `~/.devtools.yml`, returns to main menu.
- Any other key: Cancels, returns to main menu.

## Tool Execution

On `Enter` in the main menu:
1. The Bubble Tea program suspends its rendering.
2. The selected tool is launched via `os/exec` with the configured executable and args, with stdin/stdout/stderr connected to the terminal.
3. When the tool process exits, the Bubble Tea program resumes.

## Configuration File: `~/.devtools.yml`

```yaml
devtools:
  - name: "Postgres REPL"
    executable: "/usr/local/bin/psql"
    args: "-h localhost -d mydb"
```

- File is loaded once at startup.
- Manual edits while the CLI is running are not dynamically detected; a restart is required.
- If the file is missing, an empty list is used (no error).
- If the file is malformed YAML, an error is displayed and the program exits.
