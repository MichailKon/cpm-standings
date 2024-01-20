package config

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Task string

type TaskGroup struct {
	Tasks    []string `yaml:"Tasks"`
	Norm     int      `yaml:"Norm"`
	Required bool     `yaml:"Required"`
	MarkDate string   `yaml:"MarkDate"`
	Name     string   `yaml:"Name"`
}

type Criteria map[string][]*TaskGroup

type Config struct {
	Host         string `yaml:"Host"`
	CriteriaPath string `yaml:"CriteriaPath"`
	SubmitsLink  string `yaml:"SubmitsLink"`
}

func ParseCriteria(filepath string) (criteria Criteria) {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Warn("Can't open criteria file:", err.Error())
		return nil
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&criteria); err != nil {
		slog.Warn("Can't read/decode criteria file:", err.Error())
		return nil
	}
	return
}

// ParseConfig The same codes... Don't like that
func ParseConfig(filepath string) (config *Config) {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Warn("Can't open config file:", err.Error())
		return nil
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		slog.Warn("Can't read/decode config file:", err.Error())
		return nil
	}
	return
}
