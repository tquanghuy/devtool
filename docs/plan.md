# Project Plan: devtool

## Overview

`devtool` is a **centralized, configuration-driven development tool management system**. Its goal is to provide a single interactive entry point for managing *every* tool a developer needs, from global infrastructure (Docker, Telepresence) to project-specific services (Postgres, local APIs, Redis).

## Vision: The Single Point of Truth
- **Universal**: Manage ANY tool that can be checked, started, or stopped via shell commands.
- **Contextual**: Automatically adapt toolsets based on the current working directory.
- **Visual**: Power-user TUI for real-time monitoring and control.

## Architecture

The project follows a modular Go structure:

- **`cmd/devtool`**: The entry point.
- **`pkg/gui`**: TUI implementation using `tview`. Features a two-panel layout (`Tools` and `Connections`).
- **`pkg/config`**: Cascading configuration logic. Handles defaults, `~/.devtool.yml` (Global), and `./.devtool.yml` (Local).
- **`pkg/commands`**: Shell command execution wrapper (`OSCommand`) and core status checking logic.

## Core Workflows

### 1. Unified Status Monitoring
On startup, `devtool` merges all available tool definitions from global and local configs. It continuously polls their status using the defined `check_cmd`.

### 2. Context-Aware Management
When you open `devtool` in a project directory, it detects the local `.devtool.yml`. This allows you to manage microservices or databases specific to that project alongside your global tools like Docker.

### 3. Dynamic Tool Lifecycle
Users can "Add" a tool to their active workbench. For `singleton` tools, only one instance is tracked. For `portbound` tools, multiple instances (disambiguated by ports) can be managed simultaneously.

## Future Roadmap

For a detailed list of planned features and long-term vision, see the [Feature Roadmap](file:///Users/t.quanghuy/Dev/Code/devtool/docs/roadmap.md).
