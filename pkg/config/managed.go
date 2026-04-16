package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ManagedInstance represents a single installation/registration of a tool.
type ManagedInstance struct {
	ToolName   string    `json:"tool_name"`
	Identifier string    `json:"identifier"` // e.g., "postgres-5432" or "docker"
	CreatedAt  time.Time `json:"created_at"`
}

// ManagedConfig holds the list of all managed tool instances.
type ManagedConfig struct {
	Instances []ManagedInstance `json:"instances"`
}

func getManagedConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".devtool")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "managed.json"), nil
}

// LoadManagedConfig reads the managed tools state.
func LoadManagedConfig() (*ManagedConfig, error) {
	path, err := getManagedConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ManagedConfig{Instances: []ManagedInstance{}}, nil
		}
		return nil, err
	}

	var cfg ManagedConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveManagedConfig writes the managed tools state.
func SaveManagedConfig(cfg *ManagedConfig) error {
	path, err := getManagedConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// AddInstance validates and appends a new instance.
func (c *ManagedConfig) AddInstance(instance *ManagedInstance) error {
	for _, existing := range c.Instances {
		if existing.Identifier == instance.Identifier {
			return fmt.Errorf("identifier %s is already in use", instance.Identifier)
		}
	}
	instance.CreatedAt = time.Now()
	c.Instances = append(c.Instances, *instance)
	return nil
}

// RemoveInstance removes the instance with the given identifier.
func (c *ManagedConfig) RemoveInstance(identifier string) error {
	for i, inst := range c.Instances {
		if inst.Identifier == identifier {
			c.Instances = append(c.Instances[:i], c.Instances[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("identifier %s not found", identifier)
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
