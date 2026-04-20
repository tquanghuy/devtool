package config

import (
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

// ToolDefinition represents a managed developer tool's configuration and operations.
type ToolDefinition struct {
	Name        string   `yaml:"name" json:"name"`
	DefaultPort int      `yaml:"default_port" json:"default_port"`
	CheckCmd    string   `yaml:"check_cmd" json:"check_cmd"`
	StartCmd    string   `yaml:"start_cmd" json:"start_cmd"`
	StopCmd     string   `yaml:"stop_cmd" json:"stop_cmd"`
}

// ManagedInstance represents a single installation/registration of a tool.
type ManagedInstance struct {
	ToolName   string    `yaml:"tool_name" json:"tool_name"`
	Identifier string    `yaml:"identifier" json:"identifier"`
	Port       int       `yaml:"port" json:"port"`
	CreatedAt  time.Time `yaml:"created_at" json:"created_at"`
}

// AppConfig is the top-level configuration for devtool
type AppConfig struct {
	PostgresConns map[string]DatabaseConfig `yaml:"postgres_conns" json:"postgres_conns"`
	MySQLConns    map[string]DatabaseConfig `yaml:"mysql_conns" json:"mysql_conns"`
	Tools         map[string]ToolDefinition `yaml:"tools" json:"tools"`
	Managed       []ManagedInstance         `yaml:"managed" json:"managed"`
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
		PostgresConns: make(map[string]DatabaseConfig),
		MySQLConns:    make(map[string]DatabaseConfig),
		Tools:         GetDefaultTools(),
		Managed:       []ManagedInstance{},
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

	// 2. Ensure core tools are in Managed list (migration)
	cfg.ensureDefaults()

	return cfg, nil
}

func (cfg *AppConfig) ensureDefaults() {
	coreInitial := []string{"docker", "telepresence"}
	modified := false

	for _, name := range coreInitial {
		exists := false
		for _, inst := range cfg.Managed {
			if inst.ToolName == name {
				exists = true
				break
			}
		}

		if !exists {
			if def, ok := cfg.Tools[name]; ok {
				cfg.Managed = append(cfg.Managed, ManagedInstance{
					ToolName:   name,
					Identifier: name,
					Port:       def.DefaultPort,
					CreatedAt:  time.Now(),
				})
				modified = true
			}
		}
	}

	if modified {
		_ = Save(cfg)
	}
}

func mergeConfigs(base, override *AppConfig) {
	if base.PostgresConns == nil {
		base.PostgresConns = make(map[string]DatabaseConfig)
	}
	for k, v := range override.PostgresConns {
		base.PostgresConns[k] = v
	}

	if base.MySQLConns == nil {
		base.MySQLConns = make(map[string]DatabaseConfig)
	}
	for k, v := range override.MySQLConns {
		base.MySQLConns[k] = v
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
