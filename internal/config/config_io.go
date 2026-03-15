package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const managedConfigFileName = "managed.json"
const managedConfigDirName = ".devtool"

// managedConfigPath returns the absolute path to the managed config file.
func managedConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, managedConfigDirName, managedConfigFileName), nil
}

// LoadConfig reads the managed config from disk.
// If the file does not yet exist, an empty ManagedConfig is returned (not an error).
func LoadConfig() (*ManagedConfig, error) {
	path, err := managedConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &ManagedConfig{Instances: []ManagedInstance{}}, nil
		}
		return nil, err
	}

	var cfg ManagedConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Instances == nil {
		cfg.Instances = []ManagedInstance{}
	}
	return &cfg, nil
}

// SaveConfig writes the managed config atomically to disk.
// The directory is created if it does not exist.
func SaveConfig(cfg *ManagedConfig) error {
	path, err := managedConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Write to a temp file in the same directory, then rename for atomicity.
	tmpFile, err := os.CreateTemp(filepath.Dir(path), "managed-*.json.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, path)
}
