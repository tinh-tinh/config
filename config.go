package config

import (
	"errors"
	"log"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

func New[E any](path string) (*E, error) {
	var cfg E

	// Handle YAML file - load YAML first, then parse env tags
	if strings.Contains(path, ".yml") || strings.Contains(path, ".yaml") {
		yamlCfg, err := NewYaml[E](path)
		if err != nil {
			return nil, err
		}
		cfg = *yamlCfg
	} else if path != "" && !strings.Contains(path, ".env") {
		// Unsupported file type
		log.Printf("not supported type: %v\n", path)
		return nil, errors.New("not support")
	}

	// Load .env file if path contains .env (optional - won't fail if file doesn't exist)
	if strings.Contains(path, ".env") {
		_ = godotenv.Load(path)
	}

	// Always parse env tags from OS environment variables
	opts := env.Options{
		DefaultValueTagName: "default",
	}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		return nil, err
	}

	// Validate the config
	v := validator.Validator{}
	if err := v.Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
