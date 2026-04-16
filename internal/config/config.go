package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AppConfig holds the configuration for a single managed application.
type AppConfig struct {
	Name       string   `yaml:"name"`
	Executable string   `yaml:"executable"`
	Args       []string `yaml:"args"`
	Port       int      `yaml:"port"`
	LogFile    string   `yaml:"log_file"`
}

type fileSchema struct {
	Apps []AppConfig `yaml:"apps"`
}

// Load reads and parses the YAML config file at path.
func Load(path string) ([]AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	var schema fileSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return schema.Apps, nil
}
