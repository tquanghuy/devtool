package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// DatabaseConfig represents configuration for a database connection
type DatabaseConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Database string `yaml:"database" json:"database"`
}

type ToolKind string

const (
	Singleton ToolKind = "singleton"
	PortBound ToolKind = "portbound"
)

// ToolDefinition represents a managed developer tool's configuration and operations.
type ToolDefinition struct {
	Name        string   `yaml:"name" json:"name"`
	Kind        ToolKind `yaml:"kind" json:"kind"`
	DefaultPort int      `yaml:"default_port" json:"default_port"`
	CheckCmd    string   `yaml:"check_cmd" json:"check_cmd"`
	StartCmd    string   `yaml:"start_cmd" json:"start_cmd"`
	StopCmd     string   `yaml:"stop_cmd" json:"stop_cmd"`
}

// ManagedInstance represents a single installation/registration of a tool.
type ManagedInstance struct {
	ToolName   string    `yaml:"tool_name" json:"tool_name"`
	Identifier string    `yaml:"identifier" json:"identifier"`
	CreatedAt  time.Time `yaml:"created_at" json:"created_at"`
}

// AppConfig is the top-level configuration for devtool
type AppConfig struct {
	Postgres    DatabaseConfig            `yaml:"postgres" json:"postgres"`
	MySQL       DatabaseConfig            `yaml:"mysql" json:"mysql"`
	Connections map[string]DatabaseConfig `yaml:"connections" json:"connections"`
	Tools       map[string]ToolDefinition `yaml:"tools" json:"tools"`
	Managed     []ManagedInstance         `yaml:"managed" json:"managed"`
}

const (
	configFileName = ".devtool.yml"
)

// GetConfigPath returns the absolute path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, configFileName), nil
}

// Load reads the application configuration by merging defaults and global config.
func Load() (*AppConfig, error) {
	cfg := &AppConfig{
		Postgres:    DatabaseConfig{Host: "localhost", Port: 5432, User: "postgres", Database: "postgres"},
		MySQL:       DatabaseConfig{Host: "localhost", Port: 3306, User: "root", Database: "mysql"},
		Connections: make(map[string]DatabaseConfig),
		Tools:       GetDefaultTools(),
		Managed:     []ManagedInstance{},
	}

	// 1. Load Global Config
	globalPath, err := GetConfigPath()
	if err == nil {
		if data, err := os.ReadFile(globalPath); err == nil {
			var globalCfg AppConfig
			if err := yaml.Unmarshal(data, &globalCfg); err == nil {
				mergeConfigs(cfg, &globalCfg)
			}
		}
	}

	// 2. Migration Check: if Managed is empty, try migrating from managed.json
	if len(cfg.Managed) == 0 {
		if migrated, err := migrateManagedConfig(); err == nil && len(migrated) > 0 {
			cfg.Managed = migrated
			// Save immediately to persist migration
			_ = Save(cfg)
		}
	}

	return cfg, nil
}

func mergeConfigs(base, override *AppConfig) {
	if override.Postgres.Host != "" {
		base.Postgres = override.Postgres
	}
	if override.MySQL.Host != "" {
		base.MySQL = override.MySQL
	}
	for k, v := range override.Connections {
		base.Connections[k] = v
	}
	if override.Tools != nil {
		for k, v := range override.Tools {
			base.Tools[k] = v
		}
	}
	if len(override.Managed) > 0 {
		base.Managed = override.Managed
	}
}

func migrateManagedConfig() ([]ManagedInstance, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ".devtool", "managed.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var oldCfg struct {
		Instances []ManagedInstance `json:"instances"`
	}
	if err := json.Unmarshal(data, &oldCfg); err != nil {
		return nil, err
	}

	return oldCfg.Instances, nil
}


// Save writes the application configuration atomically
func Save(cfg *AppConfig) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config to yaml: %w", err)
	}

	// Use temporary file for atomic write
	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary config file: %w", err)
	}

	if err := os.Rename(tmpFile, path); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed to rename temporary config file: %w", err)
	}

	return nil
}
