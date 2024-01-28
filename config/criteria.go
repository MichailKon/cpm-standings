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

type ContestCriteria struct {
	ContestId int          `yaml:"ContestId"`
	Groups    []*TaskGroup `yaml:"Groups"`
}

type Criteria map[string]*ContestCriteria

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

func (c *Criteria) GetContestName(contestId int) string {
	for name, val := range *c {
		if val.ContestId == contestId {
			return name
		}
	}

	return ""
}
