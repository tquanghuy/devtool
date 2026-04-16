# devtool

A CLI tool to manage local developer services and configurations.

## Commands

### `devtool status`

Check the status of all configured developer tools (Postgres, MySQL, Docker, Telepresence).

```bash
devtool status
```

### `devtool add`

Add a supported tool to devtool's managed list.

```bash
# Add a singleton tool (may only be added once)
devtool add docker
devtool add telepresence

# Add a port-bound tool on its default port
devtool add postgres
devtool add mysql

# Add a port-bound tool on a specific port
devtool add postgres --port 5433

# Fail instead of prompting (useful in scripts)
devtool add postgres --non-interactive
```

**Supported tools**

| Tool         | Type       | Default Port |
|--------------|------------|--------------|
| docker       | Singleton  | —            |
| telepresence | Singleton  | —            |
| postgres     | Port-bound | 5432         |
| mysql        | Port-bound | 3306         |

### `devtool remove`

Remove a tool from the managed list.

```bash
# Remove a singleton tool
devtool remove docker

# Remove a port-bound tool (prompts to select instance if multiple exist)
devtool remove postgres

# Remove a specific instance by ID
devtool remove postgres-5433

# Terminate a running tool and remove it in one step
devtool remove postgres --force

# Non-interactive mode (fails if disambiguation or stop-confirmation is needed)
devtool remove postgres --non-interactive
```



## Configuration

Managed tool state is stored in `~/.devtool/managed.json` (created automatically on first use).
App configuration (database hosts/ports) is loaded from `~/.devtool.yml`.