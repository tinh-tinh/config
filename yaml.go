package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func NewYaml[E any](path string) (*E, error) {
	var e E

	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile get error: %v\n", err)
		return nil, err
	}

	err = yaml.Unmarshal(file, &e)
	if err != nil {
		log.Printf("get error: %v\n", err)
		return nil, err
	}
	return &e, nil
}
