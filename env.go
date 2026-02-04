package config

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

func NewEnv[E any](path string) (*E, error) {
	// Load .env file if path is provided (optional - won't fail if file doesn't exist)
	// This allows users to provide environment variables directly via OS without a .env file
	if path != "" {
		// Ignore error - .env file is optional, env.ParseWithOptions will read from OS env
		_ = godotenv.Load(path)
	}

	var cfg E
	opts := env.Options{
		DefaultValueTagName: "default",
	}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		return nil, err
	}

	v := validator.Validator{}
	if err := v.Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetRaw(key string) string {
	return os.Getenv(key)
}
