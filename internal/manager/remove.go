package manager

import (
	"fmt"
	"os/exec"
	"strings"

	"devtool/internal/config"
)

// SelectFn is a callback that presents a list of identifiers to the user and
// returns the zero-based index of the selected one.
// CLI callers provide an interactive stdin reader; unit tests inject a stub.
type SelectFn func(identifiers []string) (int, error)

// RemoveTool removes a managed tool instance.
//
//   - toolName: the canonical tool name or explicit instance identifier to remove.
//   - force: if true, a running process is terminated without prompting.
//   - nonInteractive: suppresses all prompts; returns error instead.
//   - selectFn: called when multiple instances exist and the session is interactive.
func RemoveTool(toolName string, force bool, nonInteractive bool, selectFn SelectFn, confirmFn func(string) (bool, error)) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Find matching instances — either by tool name or by exact identifier.
	instances := cfg.InstancesByTool(toolName)
	if len(instances) == 0 {
		// Maybe the user passed an exact identifier (e.g. "postgres-5432").
		for _, inst := range cfg.Instances {
			if inst.Identifier == toolName {
				instances = []config.ManagedInstance{inst}
				break
			}
		}
	}

	if len(instances) == 0 {
		return fmt.Errorf("Tool %s is not currently managed.", toolName)
	}

	// Disambiguate when multiple instances exist.
	var target config.ManagedInstance
	if len(instances) == 1 {
		target = instances[0]
	} else {
		if nonInteractive {
			return fmt.Errorf("Multiple instances of %s found. Please specify exact instance ID or run interactively.", toolName)
		}
		if selectFn == nil {
			return fmt.Errorf("multiple instances found but no selection function provided")
		}
		ids := make([]string, len(instances))
		for i, inst := range instances {
			ids[i] = inst.Identifier
		}
		idx, err := selectFn(ids)
		if err != nil {
			return err
		}
		if idx < 0 || idx >= len(instances) {
			return fmt.Errorf("invalid selection")
		}
		target = instances[idx]
	}

	// Running-instance check.
	if isRunning(target.ToolName) {
		if force {
			// Best-effort termination; ignore errors.
			_ = killProcess(target.ToolName)
		} else if nonInteractive {
			return fmt.Errorf("Tool is currently running. Stop it first or run with --force.")
		} else {
			if confirmFn == nil {
				return fmt.Errorf("Tool is currently running. Stop it first or run with --force.")
			}
			ok, err := confirmFn(target.ToolName)
			if err != nil {
				return err
			}
			if ok {
				_ = killProcess(target.ToolName)
			} else {
				return fmt.Errorf("Tool is currently running. Stop it first or run with --force.")
			}
		}
	}

	if err := cfg.RemoveInstance(target.Identifier); err != nil {
		return err
	}
	if err := config.SaveConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("Successfully removed %s from managed list.\n", target.Identifier)
	return nil
}

// isRunning returns true if a process whose name contains toolName appears running.
// This is a best-effort check and may return false on platforms where pgrep is unavailable.
func isRunning(toolName string) bool {
	out, err := exec.Command("pgrep", "-i", toolName).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}

// killProcess attempts to kill all processes matching toolName with SIGTERM.
func killProcess(toolName string) error {
	return exec.Command("pkill", "-i", toolName).Run()
}
