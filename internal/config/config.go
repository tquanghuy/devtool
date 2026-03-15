package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Database string `yaml:"database" json:"database"`
}

type AppConfig struct {
	Postgres DatabaseConfig `yaml:"postgres" json:"postgres"`
	MySQL    DatabaseConfig `yaml:"mysql" json:"mysql"`
}

func Load() (*AppConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfgPath := filepath.Join(home, ".devtool.yml")
	
	bytes, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config
			return &AppConfig{
				Postgres: DatabaseConfig{Host: "localhost", Port: 5432, User: "postgres", Database: "postgres"},
				MySQL:    DatabaseConfig{Host: "localhost", Port: 3306, User: "root", Database: "mysql"},
			}, nil
		}
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
