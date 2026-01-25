package config

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

func NewEnv[E any](path string) (*E, error) {
	if path == "" {
		path = ".env"
	}
	// We don't check for error here because it's possible that the user
	// provides environment variables directly without a .env file.
	// However, to maintain backward compatibility with previous strict behavior (as enforced by tests),
	// we will return error if godotenv fails, unless we decide to relax it.
	// The previous implementation returned error.
	if err := godotenv.Load(path); err != nil {
		return nil, err
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
