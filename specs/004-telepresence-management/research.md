# Research Notes: Telepresence Management

## Technical Context

- **Language**: Go
- **Primary Dependencies**: `os/exec` for invoking the Telepresence CLI, `charmbracelet/bubbletea` for the TUI integration.
- **Target Platform**: OS with Telepresence installed (macOS/Linux).

## Analysis of Technical Requirements

### 1. Telepresence CLI integration
- **Check Binary Existence**: Can be done using `exec.LookPath("telepresence")`.
- **Connection Status**:
  - Run `telepresence status`.
  - Parse the output. If it contains "Root Daemon: Running" and "User Daemon: Running", it's connected. If it says "Not running" or "telepresence: command not found", it's not connected.
  - Another approach is to check if `telepresence quit` throws an error or just check `telepresence status`.
- **Connect**: Run `telepresence connect`.
- **Disconnect**: Run `telepresence quit`.
- **Restart**: Can be done by running `telepresence quit` followed by `telepresence connect`.
- **Installation Prompt**: If `telepresence` binary is not found, prompt the user. We will show a markdown or list in the TUI:
  - macOS: `brew install datawire/blackbird/telepresence`
  - npm: `npm install -g telepresence` (Wait, telepresence is typically installed via brew or curl on linux).

### Decision: CLI Invocation Strategy
**Decision**: Use Go's `os/exec` package to wrapper Telepresence commands.
**Rationale**: Telepresence is primarily a CLI-driven tool. Wrapping the CLI commands directly provides a reliable way to manage it without needing complex API integrations or gRPC clients that might change across Telepresence versions.

### Decision: Asynchronous Execution for TUI
**Decision**: Execute connection operations asynchronously (via `tea.Cmd`) to avoid blocking the Bubbletea event loop.
**Rationale**: `telepresence connect` can take seconds to complete depending on cluster connectivity. Blocking the UI thread would lead to a frozen terminal interface.

### Alternatives considered:
- **Direct API integration**: Discarded. Telepresence is a complex tool with its own daemons; interacting with its internal APIs is unstable and undocumented compared to the CLI.

## Security & Privacy Impact
- Requires cluster access configured (kubeconfig). We rely entirely on the underlying environment's kubeconfig and do not manage credentials directly.

## Edge Cases Addressed
1. **Missing Binary**: Dealt with by catching `exec.Error` on `exec.LookPath` and displaying an installation guide in the TUI instead of crashing.
2. **Connection Timeouts**: Handled by context timeouts in `os/exec` if necessary, or simply bubbling up the CLI's own timeout error to the TUI.
