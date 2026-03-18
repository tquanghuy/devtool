# Research: Interactive CLI (003-interactive-cli)

## 1. Interactive TUI Library

### Decision: Bubble Tea (github.com/charmbracelet/bubbletea + bubbles)

**Rationale**:
- The spec requires interactive arrow-key navigation with the ability to popup a per-tool interactive interface. This is precisely the use-case Bubble Tea models well.
- `bubbles/list` provides an out-of-the-box arrowkey-navigable list with consistent UX.
- `bubbles/textinput` provides the single-line inputs needed for the add-tool form (Name, Executable path, Default arguments).
- Actively maintained, large ecosystem, and the recommended successor to the now-archived `survey` library.
- Constitution requires minimizing external dependencies: adding `bubbletea + bubbles + lipgloss` is a single logical unit from the Charm ecosystem. `yaml.v3` is already present.

**Alternatives Considered**:
- `promptui`: Viable for simple prompts but does not natively support a persistent full-screen TUI. The spec requires the root command itself to be an interactive UI that persists across multiple operations (add, remove, select), which goes beyond what promptui handles.
- `survey` (AlecAivazis): **Archived** as of 2024. Not suitable for new projects.

---

## 2. Profile Storage Format

### Decision: YAML (`devtools.yml` at `~/.devtools.yml`)

**Rationale**:
- `gopkg.in/yaml.v3` is already a direct dependency in `go.mod`.
- YAML is human-readable and hand-editable (aligned with the DX-first constitution).
- The user explicitly specified the filename `devtools.yml` in the feature description.
- Existing `config.go` uses `~/.devtool.yml` for the DB config. The new file is `~/.devtools.yml` (note the 's') — separate file, no collision.

**File location**: `~/.devtools.yml` (resolved via `os.UserHomeDir()`).

**Format**:
```yaml
devtools:
  - name: "My Postgres"
    executable: "/usr/local/bin/psql"
    args: "-h localhost -d mydb"
```

---

## 3. Tool Execution (popup interactive interface)

### Decision: Use `os/exec.Cmd` with stdin/stdout/stderr attached to the terminal

**Rationale**:
- The user's answer "Popup the new interactive interface for chosen tool" means the selected tool's own TUI/interactive interface should take over the terminal.
- The correct pattern is: suspend the Bubble Tea program, fork/exec the tool with os/exec with `Cmd.Stdin/Stdout/Stderr = os.Stdin/Stdout/Stderr`, wait for it to exit, then resume the Bubble Tea program.
- This is the standard pattern for TUI launchers (e.g., lazygit, k9s plugin mode).

---

## 4. Architecture Impact

The current Cobra root command (`devtool`) dispatches to subcommands (`add`, `remove`, `status`). This feature introduces a new **interactive root mode**:

- When `devtool` is run with **no arguments**, it enters the interactive TUI.
- Existing subcommands (`add`, `remove`, `status`) remain intact for non-interactive/scripted use.
- A new `internal/tui` package will house all Bubble Tea model code.
- A new `internal/devtools` package (or extend `internal/config`) will handle the `devtools.yml` CRUD.

---

## 5. Test Strategy

- **Unit tests**: Test the `devtools.yml` CRUD layer (load, save, add entry, remove entry) in isolation with `t.TempDir()` and `t.Setenv("HOME", …)` — consistent with the existing test pattern in `tests/unit/`.
- **TUI tests**: Bubble Tea programs can be tested using `bubbletea/v2`'s test helpers or by extracting pure model logic and testing the state transitions without rendering.
