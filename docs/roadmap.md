# Feature Roadmap: devtool

This document outlines the planned features and long-term vision for `devtool`. These items are sorted by functional area.

## 1. Observability & Monitoring
- [ ] **Real-time Log Streaming**: Dedicated panel to tail `stdout/stderr` for managed tools.
- [ ] **Health Check History**: A timeline showing recent status changes and exact error outputs from failed checks.
- [ ] **Process Detail View**: Modal showing PID, start time, parent process, and child threads for a selected tool.
- [ ] **Port Conflict Resolver**: A visualization of all listening ports on the system to help identify and fix collisions.

## 2. Orchestration & Stacks
- [ ] **Tool Dependencies**: Support for `depends_on` in tool definitions (e.g., API requires Postgres).
- [ ] **Tool Groups (Stacks)**: Define a collection of tools as a "Stack" to be started/stopped together (e.g., `stack: data-pipeline`).
- [ ] **Profile Management**: Save and switch between sets of "Active" tools (e.g., "Frontend Dev" vs "SRE Debugging").
- [ ] **Interactive Onboarding**: A wizard to help users configure their first global and local tools without editing YAML.

## 3. Advanced UI/UX
- [ ] **Command Palette**: `Ctrl+P` / `Cmd+K` interface for fuzzy searching tools, connections, and actions.
- [ ] **Custom Keybindings**: User-configurable shortcuts defined in `~/.devtool.yml`.
- [ ] **Themes & Aesthetics**: Support for custom color schemes and better visual indicators (e.g., icons, progress bars).
- [ ] **Search & Filter**: Real-time filtering for Tool and Connection tables.

## 4. Ecosystem & Integration
- [ ] **K8s Integration**: Context and namespace switcher for users running tools like Telepresence.
- [ ] **Git Context**: Show current branch and working directory status in the status bar.
- [ ] **Environment Variable Sync**: Manage and reload `.env` files for specific tool contexts.
- [ ] **Plugin System**: Allow third-party tool definitions to be shared and imported easily.

## 5. Persistence & Sync
- [ ] **Remote Config Sync**: Sync your `~/.devtool.yml` across multiple machines using a git repo or cloud storage.
- [ ] **State Persistence**: Preserve the state of the TUI layout and selection across restarts.

## 6. Advanced Developer Experience (DX)
- [ ] **Integrated Database REPL**: Lightweight SQL/Redis/Mongo interactive shells directly within the TUI.
- [ ] **Traffic Interception (Proxy)**: Inspect HTTP/gRPC traffic between services bound to local ports.
- [ ] **AI-Powered Diagnostics**: Automated analysis of error logs and startup failures with suggested fixes.
- [ ] **Performance Profiling**: One-click `pprof` or flamegraph generation for supported tool runtimes.

## 7. Team & Collaboration
- [ ] **Environment Snapshots**: Capture logs, state, and volume data into a sharable bundle for bug reproduction.
- [ ] **Config Sharing**: Generate "Join Codes" or ephemeral configs to help teammates spin up identical environments.
- [ ] **Local Service Discovery**: Broadcast presence of local services to other `devtool` instances on the network.
