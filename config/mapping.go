package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type StudentsHandlesMapping map[string]string

func ParseStudentsHandlesMapping(filepath string) (mapping StudentsHandlesMapping) {
	file, err := os.Open(filepath)
	if err != nil {
		return make(StudentsHandlesMapping)
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&mapping); err != nil {
		return make(StudentsHandlesMapping)
	}
	return
}
