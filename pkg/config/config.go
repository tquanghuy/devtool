package config

import (
	"fmt"
	"os"
	"path/filepath"

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

// AppConfig is the top-level configuration for devtool
type AppConfig struct {
	Postgres    DatabaseConfig            `yaml:"postgres" json:"postgres"`
	MySQL       DatabaseConfig            `yaml:"mysql" json:"mysql"`
	Connections map[string]DatabaseConfig `yaml:"connections" json:"connections"`
	Tools       map[string]ToolDefinition `yaml:"tools" json:"tools"`
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

// GetLocalConfigPath returns the path to a .devtool.yml in the current directory.
func GetLocalConfigPath() string {
	return configFileName
}


// Load reads the application configuration by merging defaults, global config, and local config.
func Load() (*AppConfig, error) {
	cfg := &AppConfig{
		Postgres:    DatabaseConfig{Host: "localhost", Port: 5432, User: "postgres", Database: "postgres"},
		MySQL:       DatabaseConfig{Host: "localhost", Port: 3306, User: "root", Database: "mysql"},
		Connections: make(map[string]DatabaseConfig),
		Tools:       GetDefaultTools(),
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

	// 2. Load Local Config
	localPath := GetLocalConfigPath()
	if data, err := os.ReadFile(localPath); err == nil {
		var localCfg AppConfig
		if err := yaml.Unmarshal(data, &localCfg); err == nil {
			mergeConfigs(cfg, &localCfg)
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
