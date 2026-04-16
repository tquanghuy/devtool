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

// DevtoolProfile represents a user-configured developer tool entry.
type DevtoolProfile struct {
	Name       string `yaml:"name"`
	Executable string `yaml:"executable"`
	Args       string `yaml:"args"`
}

// AppConfig is the top-level configuration for devtool
type AppConfig struct {
	Postgres    DatabaseConfig            `yaml:"postgres" json:"postgres"`
	MySQL       DatabaseConfig            `yaml:"mysql" json:"mysql"`
	Connections map[string]DatabaseConfig `yaml:"connections" json:"connections"`
	Devtools    []DevtoolProfile         `yaml:"devtools" json:"devtools"`
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

// Load reads the application configuration
func Load() (*AppConfig, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &AppConfig{
				Postgres:    DatabaseConfig{Host: "localhost", Port: 5432, User: "postgres", Database: "postgres"},
				MySQL:       DatabaseConfig{Host: "localhost", Port: 3306, User: "root", Database: "mysql"},
				Connections: make(map[string]DatabaseConfig),
				Devtools:    []DevtoolProfile{},
			}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml config: %w", err)
	}

	if cfg.Connections == nil {
		cfg.Connections = make(map[string]DatabaseConfig)
	}
	if cfg.Devtools == nil {
		cfg.Devtools = []DevtoolProfile{}
	}

	return &cfg, nil
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
