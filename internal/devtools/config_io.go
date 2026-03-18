package devtools

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	configFileName = ".devtools.yml"
)

// getConfigPath returns the absolute path to ~/.devtools.yml
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, configFileName), nil
}

// Load reads the devtools configuration from ~/.devtools.yml.
// If the file does not exist, it returns an empty config.
func Load() (*DevtoolsConfig, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &DevtoolsConfig{Devtools: []DevtoolProfile{}}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg DevtoolsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml config: %w", err)
	}

	if cfg.Devtools == nil {
		cfg.Devtools = []DevtoolProfile{}
	}

	return &cfg, nil
}

// Save writes the devtools configuration to ~/.devtools.yml atomically.
func Save(cfg *DevtoolsConfig) error {
	path, err := getConfigPath()
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

// Add appends a new devtool profile to the config and validates it.
func Add(cfg *DevtoolsConfig, profile DevtoolProfile) error {
	if profile.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if profile.Executable == "" {
		return fmt.Errorf("executable path cannot be empty")
	}

	for _, p := range cfg.Devtools {
		if p.Name == profile.Name {
			return fmt.Errorf("devtool with name %q already exists", profile.Name)
		}
	}

	cfg.Devtools = append(cfg.Devtools, profile)
	return nil
}

// Remove deletes a devtool profile by name from the config.
func Remove(cfg *DevtoolsConfig, name string) error {
	for i, p := range cfg.Devtools {
		if p.Name == name {
			cfg.Devtools = append(cfg.Devtools[:i], cfg.Devtools[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("devtool with name %q not found", name)
}
