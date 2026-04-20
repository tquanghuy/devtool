# Feature Roadmap: devtool

This document outlines the planned features and long-term vision for `devtool`. These items are sorted by functional area.

## 1. Observability & Monitoring
- [ ] **Real-time Log Streaming**: Dedicated panel to tail `stdout/stderr` for managed tools.
- [ ] **Health Check History**: A timeline showing recent status changes and exact error outputs from failed checks.
- [ ] **Process Detail View**: Modal showing PID, start time, parent process, and child threads for a selected tool.
- [ ] **Port Conflict Resolver**: A visualization of all listening ports on the system to help identify and fix collisions.
- [ ] **Metric Sparklines**: inline ~60s sparklines of CPU/RAM/latency per tool in the Tools table. *(new)*
- [ ] **Health-check Latency Budget**: mark a tool "degraded" (amber) when `check_cmd` exceeds a threshold. *(new)*
- [ ] **Structured Event Log**: in-memory circular event log for start/stop/check results, via a `:events` panel. *(new)*
- [x] **Async Status Polling**: move per-tool status checks off the render path into goroutines. *(new)*

## 2. Orchestration & Stacks
- [ ] **Tool Dependencies**: Support for `depends_on` in tool definitions (e.g., API requires Postgres).
- [ ] **Tool Groups (Stacks)**: Define a collection of tools as a "Stack" to be started/stopped together (e.g., `stack: data-pipeline`).
- [ ] **Profile Management**: Save and switch between sets of "Active" tools (e.g., "Frontend Dev" vs "SRE Debugging").
- [ ] **Interactive Onboarding**: A wizard to help users configure their first global and local tools without editing YAML.
- [ ] **Parallel Stack Start with DAG**: sort `depends_on` branches and start concurrently with a live tree view. *(new)*
- [ ] **Health Gating**: treat a tool "started" only after `check_cmd` passes or `wait_for` expression succeeds. *(new)*
- [ ] **Dry-run Mode**: `devtool stack up --dry-run` prints the planned order; TUI previews it in a modal. *(new)*
- [ ] **Rollback on Failure**: on any node failure, stop already-started tools in reverse order. *(new)*

## 3. Advanced UI/UX
- [ ] **Command Palette**: `Ctrl+P` / `Cmd+K` interface for fuzzy searching tools, connections, and actions.
- [ ] **Custom Keybindings**: User-configurable shortcuts defined in `~/.devtool.yml`.
- [ ] **Themes & Aesthetics**: Support for custom color schemes and better visual indicators (e.g., icons, progress bars).
- [ ] **Search & Filter**: Real-time filtering for Tool and Connection tables.
- [ ] **Detached Log/Details Pane**: tear off Details into a full-screen or side-by-side split (toggled with `|` / `_`). *(new)*
- [ ] **Status-bar Breadcrumbs**: show `project · profile · active-stack · git-branch` for constant context. *(new)*
- [ ] **Keyboard-Driven Multi-select**: `Space` to mark multiple tools, `b` for batch start/stop/restart. *(new)*
- [ ] **Action Hints in Rows**: inline `[enter] [r] [l]` hints in the selected row so shortcuts are discoverable. *(new)*
- [ ] **Responsive Layout**: collapse the three-panel layout into tabs when terminal width drops below ~120 cols. *(new)*
- [ ] **Notification Toasts**: auto-dismissing toasts for action success/failure so modals stop stealing focus. *(new)*
- [ ] **Empty States**: centered helper text plus `[a] add your first tool` when no tools are managed. *(new)*
- [ ] **Consistent Modal Geometry**: extract a `modal.Center(width, height)` helper so every modal looks identical. *(new)*
- [ ] **Context-sensitive Status Bar**: shortcut bands that change as focus moves between Tools / Connections. *(new)*
- [ ] **Theme Tokens**: replace hard-coded colors with semantic tokens (`accent`, `success`, `warn`, `danger`) from config. *(new)*
- [ ] **High-contrast & Colorblind-safe Preset**: ship a built-in `theme: accessible` using Okabe-Ito hues. *(new)*
- [ ] **Optional Mouse Support**: `mouse: true` config toggle for pairing sessions, without compromising keyboard flow. *(new)*

## 4. Ecosystem & Integration
- [ ] **K8s Integration**: Context and namespace switcher for users running tools like Telepresence.
- [ ] **Git Context**: Show current branch and working directory status in the status bar.
- [ ] **Environment Variable Sync**: Manage and reload `.env` files for specific tool contexts.
- [ ] **Plugin System**: Allow third-party tool definitions to be shared and imported easily.
- [ ] **Docker Compose Import**: `devtool import compose ./docker-compose.yml` generates `.devtool.yml` entries. *(new)*
- [ ] **Tilt/Skaffold Bridge**: point devtool at an existing Tiltfile config to get read-only status surfaces. *(new)*
- [ ] **Shell Hook Generator**: `devtool init zsh` emits a prompt hook summarizing current project health. *(new)*
- [ ] **Editor Quick Actions**: VS Code / JetBrains companion exposing "Start stack" commands via `managed.json`. *(new)*

## 5. Persistence & Sync
- [ ] **Remote Config Sync**: Sync your `~/.devtool.yml` across multiple machines using a git repo or cloud storage.
- [ ] **State Persistence**: Preserve the state of the TUI layout and selection across restarts.
- [ ] **Workspace Sessions**: named snapshots of active tools + layout; restore with `devtool session load <name>`. *(new)*
- [ ] **Config Validation & Lint**: `devtool config lint` — schema-validated YAML with line-numbered error messages. *(new)*
- [ ] **Secrets Indirection**: `${env:VAR}` and `${op:...}` placeholders in commands so `.devtool.yml` is committed safely. *(new)*

## 6. Advanced Developer Experience (DX)
- [ ] **Integrated Database REPL**: Lightweight SQL/Redis/Mongo interactive shells directly within the TUI.
- [ ] **Traffic Interception (Proxy)**: Inspect HTTP/gRPC traffic between services bound to local ports.
- [ ] **AI-Powered Diagnostics**: Automated analysis of error logs and startup failures with suggested fixes.
- [ ] **Performance Profiling**: One-click `pprof` or flamegraph generation for supported tool runtimes.
- [x] **Just-in-time Port Allocator**: when adding a tool, suggest the next free port and pre-fill the identifier. *(new)*
- [ ] **Log Follow + Filter**: `f` on any tool opens a tailing log view with regex highlight and in-buffer search. *(new)*
- [ ] **Recipe Sharing**: `devtool add <url>` pulls a signed tool definition and prompts for approval before trust. *(new)*
- [ ] **AI Log Triage**: on a failed `start_cmd`, offer a `?` action that summarizes stderr via an LLM endpoint. *(new)*

## 7. Team & Collaboration
- [ ] **Environment Snapshots**: Capture logs, state, and volume data into a sharable bundle for bug reproduction.
- [ ] **Config Sharing**: Generate "Join Codes" or ephemeral configs to help teammates spin up identical environments.
- [ ] **Local Service Discovery**: Broadcast presence of local services to other `devtool` instances on the network.
- [ ] **Read-only Share URL**: time-boxed signed link exposing a read-only web view of the current stack's status. *(new)*
- [ ] **Reproducible Bug Bundles**: `devtool bundle` zips managed.json, config, recent events, and redacted logs. *(new)*

## 8. Security & Trust
- [ ] **Command Sandbox Profile**: allow-list patterns for commands so unreviewed YAML can't run destructive actions. *(new)*
- [ ] **Signed Tool Definitions**: optional SSH-signature or minisign verification when importing from a registry. *(new)*
- [ ] **Least-privilege Mode**: surface each tool's required capabilities (network, files) in Details before first run. *(new)*
- [ ] **Audit Trail**: append-only log capturing every lifecycle action with user, cwd, and exit status. *(new)*

## 9. Testing, Mocks & CI Parity
- [ ] **Mock Backend**: `mock:` tool backend that simulates status/ports without invoking real binaries. *(new)*
- [ ] **CI Smoke Recipe**: `devtool ci --stack <name>` non-interactive runner that starts, gates, runs, and tears down. *(new)*
- [ ] **Golden Config Tests**: snapshot-test the cascaded merge of defaults + global + local configs. *(new)*

## Proposed Prioritization (Checklist)

This checklist tracks the implementation progress of new items, sorted by prioritization (Impact × Effort).

- [x] **#01: Async Status Polling** (Section 1) • Impact: L • Effort: S
- [x] **#02: Just-in-time Port Allocator** (Section 6) • Impact: L • Effort: S (Deps: Identifier bug fix)
- [ ] **#03: Log Follow + Filter** (Section 6) • Impact: L • Effort: M (Deps: Async polling)
- [ ] **#04: Docker Compose Import** (Section 4) • Impact: L • Effort: M
- [ ] **#05: Metric Sparklines** (Section 1) • Impact: M • Effort: S (Deps: Async polling)
- [ ] **#06: Structured Event Log** (Section 1) • Impact: M • Effort: S
- [ ] **#07: Notification Toasts** (Section 3) • Impact: M • Effort: S
- [ ] **#08: Empty States** (Section 3) • Impact: S • Effort: S
- [ ] **#09: Consistent Modal Geometry** (Section 3) • Impact: S • Effort: S
- [ ] **#10: Theme Tokens** (Section 3) • Impact: M • Effort: S
- [ ] **#11: Responsive Layout** (Section 3) • Impact: M • Effort: M
- [ ] **#12: Status-bar Breadcrumbs** (Section 3) • Impact: S • Effort: S
- [ ] **#13: Context-sensitive Status Bar** (Section 3) • Impact: S • Effort: S
- [ ] **#14: High-contrast Preset** (Section 3) • Impact: M • Effort: S (Deps: Theme tokens)
- [ ] **#15: Optional Mouse Support** (Section 3) • Impact: S • Effort: S
- [ ] **#16: Action Hints in Rows** (Section 3) • Impact: M • Effort: S
- [ ] **#17: Keyboard-Driven Multi-select** (Section 3) • Impact: M • Effort: M
- [ ] **#18: Detached Log/Details Pane** (Section 3) • Impact: M • Effort: M (Deps: Log Follow)
- [ ] **#19: Health-check Latency Budget** (Section 1) • Impact: M • Effort: S (Deps: Async polling)
- [ ] **#20: Parallel Stack Start with DAG** (Section 2) • Impact: XL • Effort: L (Deps: `depends_on` support)
- [ ] **#21: Health Gating** (Section 2) • Impact: L • Effort: M (Deps: `depends_on` support)
- [ ] **#22: Dry-run Mode** (Section 2) • Impact: M • Effort: S (Deps: `depends_on` support)
- [ ] **#23: Rollback on Failure** (Section 2) • Impact: L • Effort: M (Deps: Health Gating)
- [ ] **#24: Workspace Sessions** (Section 5) • Impact: M • Effort: M
- [ ] **#25: Config Validation & Lint** (Section 5) • Impact: M • Effort: M
- [ ] **#26: Secrets Indirection** (Section 5) • Impact: M • Effort: M
- [ ] **#27: Shell Hook Generator** (Section 4) • Impact: M • Effort: S (Deps: `devtool status` CLI)
- [ ] **#28: Editor Quick Actions** (Section 4) • Impact: M • Effort: L (Deps: `devtool status` CLI)
- [ ] **#29: Tilt/Skaffold Bridge** (Section 4) • Impact: S • Effort: L
- [ ] **#30: Recipe Sharing** (Section 6) • Impact: M • Effort: M (Deps: Signed Tool Definitions)
- [ ] **#31: AI Log Triage** (Section 6) • Impact: M • Effort: M (Deps: Log Follow)
- [ ] **#32: Read-only Share URL** (Section 7) • Impact: M • Effort: L (Deps: Event Log)
- [ ] **#33: Reproducible Bug Bundles** (Section 7) • Impact: M • Effort: S (Deps: Event Log)
- [ ] **#34: Command Sandbox Profile** (Section 8) • Impact: L • Effort: M
- [ ] **#35: Signed Tool Definitions** (Section 8) • Impact: L • Effort: M
- [ ] **#36: Least-privilege Mode** (Section 8) • Impact: M • Effort: S
- [ ] **#37: Audit Trail** (Section 8) • Impact: M • Effort: S (Deps: Event Log)
- [ ] **#38: Mock Backend** (Section 9) • Impact: L • Effort: M
- [ ] **#39: CI Smoke Recipe** (Section 9) • Impact: L • Effort: M (Deps: Mock Backend)
- [ ] **#40: Golden Config Tests** (Section 9) • Impact: M • Effort: S

### Pick these up first (top 5, ordered)

1. **Async Status Polling** — removes a latent UI-blocking bug and unlocks sparklines, degraded-health, and log follow.
2. **Just-in-time Port Allocator** — closes the singleton/port-bound identifier collision devs hit on day one.
3. **Log Follow + Filter** — delivers the most-requested TUI capability and justifies the pitch.
4. **Docker Compose Import** — fastest path to real-world adoption; removes the hand-crafted-YAML barrier.
5. **Structured Event Log + Notification Toasts** — shared substrate for audit trail, bug bundles, and UI feedback.
