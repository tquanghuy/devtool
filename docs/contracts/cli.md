# CLI Contracts

The `devtool` CLI uses `cobra` to parse and execute commands. 

## Base Command
`devtool` (No Arguments)
- **Action**: Launches the interactive Terminal User Interface (TUI).
- **Behavior**: If no subcommands are provided, `RunE` delegates control to `tui.Start()`, which initializes the Bubble Tea application representing the default visual workspace.

## `devtool status`
- **Action**: Outputs the running status of all managed tools to the terminal standard output.
- **Behavior**: Iterates over `ManagedConfig`, pings local ports or checks tool APIs, and prints active vs inactive states.

## `devtool add <tool-name>`
- **Action**: Registers a new tool in the `ManagedConfig` list.
- **Arguments**:
  - `<tool-name>` (Required): The type of tool to add. Valid values include `docker`, `telepresence`, `postgres`, `mysql`.
- **Flags**:
  - `--port <int>` (Optional): Specifically request a port for port-bound tools (e.g., databases). If used on a singleton, it is ignored or causes an error. Default: `0` (resolves to tool default).
  - `--non-interactive` (Optional): Forces the command to fail instead of prompting the user for input if port mapping resolution or disambiguation is required.
- **Behavior**: Singletons are added using their name as the identifier. Port-bound tools generate an identifier like `postgres-5432`.

## `devtool remove <tool-name>|<instance-id>`
- **Action**: Removes a tool from `ManagedConfig`.
- **Arguments**:
  - `<tool-name>|<instance-id>` (Required): The identifier of the instance or the generic tool name.
- **Flags**:
  - `--force`, `-f` (Optional): Bypasses the confirmation prompt to gracefully stop the tool. Forcefully terminates the process before removal.
  - `--non-interactive` (Optional): Fails if the tool is currently running (requires manual stop first) or if multiple instances of the same tool name exist without a specific `<instance-id>` provided.
- **Behavior**: 
  - If a tool name like `postgres` is provided but multiple instances exist (e.g., `postgres-5432`, `postgres-5433`), the user is prompted to select which one to remove via an interactive prompt (unless `--non-interactive` is passed).
  - If the tool is currently running, the user is prompted to force stop it before removal (unless `--force` or `--non-interactive` is utilized).
