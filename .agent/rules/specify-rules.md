# devtool Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-03-23

## Active Technologies
- Go (latest) (001-devtool-management)
- OS Native Keychain (for DB credentials) (001-devtool-management)
- Go 1.21+ + `spf13/cobra` (CLI framework), internal config/checker packages (002-devtool-add-remove)
- Local configuration file (JSON/YAML via internal config package) (002-devtool-add-remove)
- Go 1.25 + `github.com/spf13/cobra` (existing), `gopkg.in/yaml.v3` (existing), `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/bubbles` (new) (003-interactive-cli)
- `~/.devtools.yml` (YAML, via `yaml.v3`) (003-interactive-cli)

- (001-devtool-management)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for 

## Code Style

: Follow standard conventions

## Recent Changes
- 004-telepresence-management: Added [if applicable, e.g., PostgreSQL, CoreData, files or N/A]
- 003-interactive-cli: Added Go 1.25 + `github.com/spf13/cobra` (existing), `gopkg.in/yaml.v3` (existing), `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/bubbles` (new)
- 002-devtool-add-remove: Added Go 1.21+ + `spf13/cobra` (CLI framework), internal config/checker packages


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
