# Quickstart: devtool-add-remove

This quickstart guides you through the new devtool manage features.

## Prerequisites

- Install the latest version of the devtool CLI.
- Ensure your host environment allows applications to bind to local ports if you're managing port-bound tools.

## Adding a Tool

To start managing a tool (e.g., PostgreSQL) locally:

```bash
$ devtool add postgres
Successfully added postgres to managed list.
```

If it is a port-bound tool, you can specify a specific port:

```bash
$ devtool add postgres --port 5433
Successfully added postgres to managed list.
```

If you try to add a singleton tool (e.g., Docker) that is already managed:

```bash
$ devtool add docker
Error: Tool docker is a singleton and is already managed.
```

## Removing a Tool

To remove a tool that you are no longer using:

```bash
$ devtool remove docker
Successfully removed docker from managed list.
```

If you have multiple instances of a port-bound tool, the CLI will prompt you to select which one to remove:

```bash
$ devtool remove postgres
Multiple instances found:
1) postgres (port 5432)
2) postgres (port 5433)
Select instance to remove [1-2]: 2
Successfully removed postgres (port 5433) from managed list.
```

If the tool is currently running, `devtool` will ask if you want to stop it before removal. You can bypass this check and forcefully remove it using the `--force` flag:

```bash
$ devtool remove postgres --force
Successfully removed postgres from managed list.
```
