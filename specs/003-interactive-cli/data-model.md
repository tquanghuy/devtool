# Data Model: Interactive CLI (003-interactive-cli)

## DevtoolProfile

Represents a user-configured developer tool entry. Stored in a YAML list in `~/.devtools.yml`.

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | `string` | Yes | Human-readable display name shown in the interactive menu |
| `executable` | `string` | Yes | Absolute or `$PATH`-resolvable path to the tool's binary |
| `args` | `string` | No | Space-separated default arguments passed to the executable on launch |

### Validation Rules

- `name` must be non-empty and unique across all profiles.
- `executable` must be non-empty.
- `args` may be empty string.

### Example `~/.devtools.yml`

```yaml
devtools:
  - name: "Postgres REPL"
    executable: "/usr/local/bin/psql"
    args: "-h localhost -d mydb"
  - name: "Docker TUI"
    executable: "/usr/local/bin/lazydocker"
    args: ""
```

## DevtoolsConfig (root document)

Top-level YAML document that wraps the list of profiles.

| Field | Type | Description |
|-------|------|-------------|
| `devtools` | `[]DevtoolProfile` | Ordered list of configured devtool profiles |

### State Transitions

```
No file exists  →  empty DevtoolsConfig (no error on load)
Empty list      →  user adds tool → list grows
Non-empty list  →  user removes tool (with confirmation) → list shrinks
```
