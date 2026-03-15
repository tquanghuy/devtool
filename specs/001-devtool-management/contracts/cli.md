# Contracts

This CLI tool does not expose an external HTTP API. Its primary interface contract is the terminal CLI and the configuration file.

## `~/.devtool.yml` Configuration Schema

Users must define their database connections in a YAML file in their home directory (or via equivalent environment variables supported by Viper).

```yaml
tools:
  docker:
    enabled: true
  telepresence:
    enabled: true
  postgres:
    enabled: true
    connection:
      host: "localhost"
      port: 5432
      user: "postgres"
      password: "password123"
      database: "devlocal"
  mysql:
    enabled: true
    connection:
      host: "127.0.0.1"
      port: 3306
      user: "root"
      password: "password"
      database: "devlocal"
```

## CLI Output Format (Standard Out)

Success Output (Tabular Format):
```
TOOL            STATUS      LATENCY    MESSAGE
Docker          Running     12ms       Daemon active
Telepresence    Stopped     5ms        Daemon not found
Postgres        Connected   45ms       Ping successful
MySQL           Failed      3000ms     Connection refused
```
