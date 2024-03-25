package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Config struct {
	Host         string `yaml:"Host"`
	CriteriaPath string `yaml:"CriteriaPath"`
	ApiKey       string `yaml:"ApiKey"`
	Secret       string `yaml:"Secret"`
}

// ParseConfig The same codes... Don't like that
func ParseConfig(filepath string) (config *Config) {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Warn(fmt.Sprintf("Can't open config file: %v", err.Error()))
		return nil
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		slog.Warn(fmt.Sprintf("Can't read/decode config file: %v", err.Error()))
		return nil
	}
	return
}
