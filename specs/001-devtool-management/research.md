# Phase 0: Outline & Research

## Decision 1: Target Language
- **Decision**: Go (1.21+)
- **Rationale**: The user has recently built CLI tools in Go (`ctx-aggregator`, `e2ee`). Go is the ideal language for cross-platform single-binary CLI tools that require high concurrency (for checking 4 background services simultaneously).
- **Alternatives considered**: Bash (hard to test and scale concurrency gracefully), Python (requires managing virtual environments or packaging runtimes for distribution).

## Decision 2: CLI Framework
- **Decision**: `spf13/cobra`
- **Rationale**: standard de-facto library for Go CLIs. It handles commands, flags, and help menus out-of-the-box perfectly matching the Constitution principle "Simple CLI Interface".
- **Alternatives considered**: Standard library `flag` (too basic for potential future subcommands like adding/removing tools from management), `urfavee/cli` (less standard than Cobra in modern Go ecosystem).

## Decision 3: Secure Credential Storage
- **Decision**: `zalando/go-keyring`
- **Rationale**: The specification requires storing database credentials securely in the OS Native Keychain (FR-007). `zalando/go-keyring` provides a simple cross-platform API for macOS and Linux secret services, fitting the local developer environment perfectly without complex CGO dependencies on macOS.
- **Alternatives considered**: `99designs/keyring` (more robust but more complex to configure, overkill for simple local DB credentials). Plain YAML file natively via Viper (rejected during the clarification phase due to security concerns with plaintext passwords).
