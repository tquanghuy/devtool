# Phase 1: Data Model & Contracts

## Entities

### `ToolType` (Enum/Type)
- **Values**: `Docker`, `Telepresence`, `Postgres`, `MySQL`
- **Purpose**: Strongly type the tools the system can check.

### `ToolConfig`
- **Fields**:
  - `Name` (ToolType)
  - `Enabled` (bool)
  - `Connection` (DatabaseConfig, optional)
- **Purpose**: Represents a single tool's configuration as loaded from `.devtool.yml`.

### `DatabaseConfig`
- **Fields**:
  - `Host` (string)
  - `Port` (int)
  - `User` (string)
  - `Database` (string)
- **Purpose**: Stores routing/connection details for network-based tools. 
- **Validation Rules**: If a DB tool is enabled, Port and User must be present. *(Note: Passwords are not stored here. They are fetched from the OS Native Keychain at runtime).*

### `StatusResult`
- **Fields**:
  - `Tool` (ToolType)
  - `IsActive` (bool)
  - `Error` (error, optional)
  - `Latency` (time.Duration)
- **Purpose**: The result of a single check execution.

## State Transitions
1. `Pending` -> `Running` (when the checker starts).
2. `Running` -> `Active` (check passes).
3. `Running` -> `Inactive` (check fails or times out).

## Contracts
No external web APIs or schemas exposed. The primary contract is the CLI output format (tabular text out to stdout).
