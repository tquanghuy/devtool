# Data Model Specification

The state and configuration for `devtool` are maintained across two local files, ensuring a separation between user-defined connection settings and tool lifecycle tracking.

## 1. AppConfig (`~/.devtool.yml`)

Stores static application configuration and user-defined connections.

```yaml
postgres:
  host: "localhost"
  port: 5432
  user: "postgres"
  database: "postgres"
mysql:
  host: "localhost"
  port: 3306
  user: "root"
  database: "mysql"
connections:
  # Additional dynamic connections
```

### Structures

- **`AppConfig`**
  - **`Postgres`** (`DatabaseConfig`): Default configuration for Postgres.
  - **`MySQL`** (`DatabaseConfig`): Default configuration for MySQL.
  - **`Connections`** (`map[string]DatabaseConfig`): Dynamically added custom connections.
  
- **`DatabaseConfig`**
  - `Host` (string)
  - `Port` (int)
  - `User` (string)
  - `Database` (string)

## 2. ManagedConfig (`~/.devtool/managed.json`)

Tracks the lifecycle of tools that have been "added" to `devtool` and are actively managed.

```json
{
  "instances": [
    {
      "tool_name": "docker",
      "identifier": "docker",
      "created_at": "2026-03-15T15:00:00Z"
    },
    {
      "tool_name": "postgres",
      "identifier": "postgres-5433",
      "created_at": "2026-03-15T15:01:00Z"
    }
  ]
}
```

### Structures

- **`ManagedConfig`**
  - **`Instances`** (`[]ManagedInstance`): The root array of managed instances.

- **`ManagedInstance`**
  - **`ToolName`** (string): The canonical name of the tool (e.g., `docker`, `postgres`).
  - **`Identifier`** (string): The unique key for the instance. For singletons, this equals `ToolName`. For port-bound tools, it follows the format `<ToolName>-<port>`.
  - **`CreatedAt`** (time.Time): Timestamp of when the tool was added to management.

### Validation Rules

- **Singleton Check**: A tool designated as a singleton (like `docker` or `telepresence`) cannot be added if an instance with the same `ToolName` already exists.
- **Identifier Uniqueness**: The `Identifier` must be globally unique within the `Instances` array. For example, `postgres-5432` can only appear once.
