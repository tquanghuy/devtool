# Quickstart: devtool-management

This tool allows you to quickly check the status of local developer dependencies (Docker, Telepresence, Postgres, and MySQL).

## Prerequisites

- Go 1.25+ installed.

## Configuration

Before checking connection-based tools (like databases), you need to configure their connection details. Create a `.devtool.yml` file in your home directory (`~/.devtool.yml`):

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
```

**Secure Credentials**:
Passwords are NOT stored in the configuration file. You must set them securely in your OS Native Keychain using the devtool CLI:
```bash
devtool config set-credential postgres
```

*Note: Process-based tools like Docker and Telepresence do not strictly require configuration to check if they are running.*

## Building and Running

1. **Build the CLI**:
   ```bash
   go build -o devtool ./cmd/devtool
   ```

2. **Check Developer Tools Status**:
   ```bash
   ./devtool status
   ```

   **Example Output**:
   ```text
   TOOL                   STATUS    LATENCY    MESSAGE
   PostgreSQL Database    UP        45ms       -
   MySQL Database         DOWN      0ms        dial tcp [::1]:3306: connect: connection refused
   Docker Daemon          UP        12ms       -
   Telepresence           DOWN      5ms        daemon not running
   ```
