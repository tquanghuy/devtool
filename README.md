# devtool

`devtool` is a **centralized, configuration-driven development tool management system**. It acts as a single point of truth for every tool a developer needs, from global infrastructure like Docker and Telepresence to project-specific microservices and databases.

## The Goal

The goal of `devtool` is to simplify the local development environment by providing a unified interface to:
- **Orchestrate Tools**: Start, stop, and monitor the health of any developer tool.
- **Context-Specific Configuration**: Automatically load different toolsets based on the project you are working in (Global vs. Local `.devtool.yml`).
- **Universal Extensibility**: Define *any* tool (Redis, Kafka, local APIs) using simple shell commands.
- **Visual Management**: An interactive TUI (Terminal User Interface) influenced by `lazydocker` and `k9s` for at-a-glance status and control.

## Usage

Simply run `devtool` to launch the interactive TUI:

```bash
devtool
```

### Key Features
- **Cascading Config**: Merges embedded defaults, global `~/.devtool.yml`, and project-local `./.devtool.yml`.
- **Tool Kinds**: Supports `singleton` (daemon-like) and `portbound` (database-like) resources.
- **Rich TUI**: Manage connections, tools, and logs in a unified, keyboard-centric interface.

## Configuration

### Global Configuration (`~/.devtool.yml`)
Define tools and settings that you want available across all projects.

### Local Configuration (`./.devtool.yml`)
Define tools specific to your current repository. Local definitions can override global ones for tailored project workflows.

```yaml
tools:
  my-local-api:
    name: "local-api"
    kind: "singleton"
    start_cmd: "go run main.go"
    stop_cmd: "pkill -f local-api"
    check_cmd: "curl -s localhost:8080/health"
```

## Documentation

More detailed information can be found in the `docs/` directory:
- [Quickstart Guide](docs/quickstart.md): Get up and running in minutes.
- [Project Structure](docs/project-structure.md): Deep dive into the codebase and architecture.
- [Data Model](docs/data-model.md): Specification of configuration and state files.
- [Roadmap](docs/roadmap.md): Planned features and future directions.
- [CLI Contract](docs/contracts/cli.md): Commands and flags specification.

## Storage
- **Managed Tool Entries**: `~/.devtool/managed.json` tracks which tools you have "added" to your active workbench.
- **Dynamic State**: Status is checked in real-time using your defined `check_cmd`.