package config

import (
	"fmt"
	"time"
)

// ManagedConfig is the root object serialised to ~/.devtool/managed.json.
type ManagedConfig struct {
	Instances []ManagedInstance `json:"instances"`
}

// AddInstance validates and appends a new instance.
//
//   - For Singleton tools: rejects if ANY instance with the same ToolName already exists.
//   - For all tools: rejects if an instance with the same Identifier already exists.
//
// It does NOT set CreatedAt — callers are expected to set it before calling.
func (c *ManagedConfig) AddInstance(instance *ManagedInstance) error {
	for _, existing := range c.Instances {
		// Singleton duplicate check (ToolName-level).
		if existing.ToolName == instance.ToolName && existing.Identifier == instance.ToolName {
			// The existing instance uses the tool name as its identifier: it is a singleton entry.
			return fmt.Errorf("Tool %s is a singleton and is already managed.", instance.ToolName)
		}
		// Identifier uniqueness check (all types).
		if existing.Identifier == instance.Identifier {
			return fmt.Errorf("identifier %s is already in use", instance.Identifier)
		}
	}
	instance.CreatedAt = time.Now()
	c.Instances = append(c.Instances, *instance)
	return nil
}

// RemoveInstance removes the instance with the given identifier.
// Returns an error if no such instance exists.
func (c *ManagedConfig) RemoveInstance(identifier string) error {
	for i, inst := range c.Instances {
		if inst.Identifier == identifier {
			c.Instances = append(c.Instances[:i], c.Instances[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("identifier %s not found in managed list", identifier)
}

// InstancesByTool returns all managed instances for the given tool name.
func (c *ManagedConfig) InstancesByTool(toolName string) []ManagedInstance {
	var result []ManagedInstance
	for _, inst := range c.Instances {
		if inst.ToolName == toolName {
			result = append(result, inst)
		}
	}
	return result
}
