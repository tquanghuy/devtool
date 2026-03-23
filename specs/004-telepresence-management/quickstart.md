# Quickstart: Telepresence Management

## Prerequisites

- Go 1.21+ installed
- Telepresence CLI tool installed (or the capability to install it when prompted)

## Setup

1. Build the devtool binary in the root directory:
   ```bash
   go build -o devtool ./cmd/devtool
   ```

2. Run the tool:
   ```bash
   ./devtool
   ```

## Usage

1. **Add Telepresence**:
   - In the devtool TUI, navigate to the `Add` section.
   - Select `Telepresence` from the list of available tools.
   - The tool will be added to your managed tools list.

2. **Manage Connection**:
   - In your managed tools list, select `Telepresence`.
   - The UI will display the current connection status (Connected / Not Connected).
   - Press Enter to view actions: Connect, Disconnect, Restart, Remove.
   - Select the desired action to interact with your cluster.

3. **Installation Prompt**:
   - If Telepresence is not installed, selecting it or trying to connect will show an installation prompt with commands for brew/npm.
   - Install the CLI outside of the application, then return to use it.
