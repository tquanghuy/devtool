# Contracts: devtool-add-remove

## CLI Commands

### 1. `devtool add`
Adds a supported tool to the managed list.

**Syntax**: 
`devtool add <tool-name> [flags]`

**Arguments**:
- `<tool-name>` (Required): The name of the predefined tool to add (e.g., `docker`, `postgres`).

**Flags**:
- `--port <number>` (Optional): Specifically request a port for port-bound tools. Ignored or causes error for singletons.
- `--non-interactive` (Optional): Fails instead of prompting if port mapping resolution is required for port-bound tools.

**Outputs**:
- Success (0): "Successfully added `<tool-name>` to managed list."
- Error (1): "Tool `<tool-name>` is not supported."
- Error (1): "Tool `<tool-name>` is a singleton and is already managed."
- Error (1): "Port conflict for `<tool-name>`. Please specify a different port or run interactively."

### 2. `devtool remove`
Removes a tool from the managed list.

**Syntax**:
`devtool remove <tool-name>|<instance-id> [flags]`

**Arguments**:
- `<tool-name>|<instance-id>` (Required): The name of the tool or the specific instance identifier to remove.

**Flags**:
- `--force` / `-f` (Optional): Bypasses the prompt to gracefully stop the tool and forcefully terminates it before removal.
- `--non-interactive` (Optional): Fails if the tool is currently running (requires manual stop first) or if multiple instances exist without a specific `<instance-id>` provided.

**Outputs**:
- Success (0): "Successfully removed `<instance-id>` from managed list."
- Error (1): "Tool `<tool-name>` is not currently managed."
- Error (1): "Multiple instances of `<tool-name>` found. Please specify exact instance ID or run interactively."
- Error (1): "Tool is currently running. Stop it first or run with --force."
