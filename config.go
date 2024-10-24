package config

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func New[E any](path string) (*E, error) {
	if path == "" {
		path = ".env"
	}
	fmt.Println("path is ", path)

	if strings.Contains(path, ".env") {
		return NewEnv[E](path)
	} else if strings.Contains(path, ".yml") || strings.Contains(path, ".yaml") {
		return NewYaml[E](path)
	} else {
		log.Printf("not supported type: %v\n", path)
		return nil, errors.New("not support")
	}
}
