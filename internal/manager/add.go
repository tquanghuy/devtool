package manager

import (
	"fmt"
	"strings"

	"devtool/internal/checker"
	"devtool/internal/config"
)

// AddTool adds a tool to the managed list.
//
// toolName is the canonical name (e.g. "docker", "postgres").
// port overrides the default port for port-bound tools (0 means use default).
// nonInteractive suppresses prompts; failures are returned as errors.
func AddTool(toolName string, port int, nonInteractive bool) error {
	def, ok := checker.LookupTool(toolName)
	if !ok {
		supported := strings.Join(checker.SupportedToolNames(), ", ")
		fmt.Printf("Supported tools: %s\n", supported)
		return fmt.Errorf("Tool %s is not supported.", toolName)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var instance *config.ManagedInstance

	switch def.Kind {
	case checker.Singleton:
		instance = &config.ManagedInstance{
			ToolName:   def.Name,
			Identifier: def.Name,
		}

	case checker.PortBound:
		targetPort := def.DefaultPort
		if port != 0 {
			targetPort = port
		}

		identifier := fmt.Sprintf("%s-%d", def.Name, targetPort)

		// Check if this identifier is already taken.
		for _, inst := range cfg.Instances {
			if inst.Identifier == identifier {
				return fmt.Errorf("Port conflict for %s. Please specify a different port or run interactively.", def.Name)
			}
		}

		instance = &config.ManagedInstance{
			ToolName:   def.Name,
			Identifier: identifier,
		}
	}

	if err := cfg.AddInstance(instance); err != nil {
		return err
	}

	return config.SaveConfig(cfg)
}
