package config_test

import (
	"testing"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/config/v2"
)

type Config struct {
	NodeEnv   string        `env:"NODE_ENV"`
	Port      int           `env:"PORT"`
	ExpiresIn time.Duration `env:"EXPIRES_IN"`
	Log       bool          `env:"LOG"`
	Special   interface{}   `env:"-"`
	Secret    string        `env:"SECRET"`
}

func Test_Scan(t *testing.T) {
	err := godotenv.Load(".env.example")
	require.Nil(t, err)
	var cfg Config
	err = env.Parse(&cfg)
	require.Nil(t, err)

	require.Equal(t, "development", cfg.NodeEnv)
	require.Equal(t, 5000, cfg.Port)
	require.Equal(t, 5*time.Minute, cfg.ExpiresIn)
	require.Equal(t, false, cfg.Log)
	require.Equal(t, "", cfg.Secret)
}

func Test_New_Env(t *testing.T) {
	cfg, err := config.NewEnv[Config](".env.example")
	require.Nil(t, err)

	require.Equal(t, "development", cfg.NodeEnv)
	require.Equal(t, 5000, cfg.Port)
	require.Equal(t, 5*time.Minute, cfg.ExpiresIn)
}

func Test_GetRaw(t *testing.T) {
	_, err := config.NewEnv[Config](".env.example")
	require.Nil(t, err)

	dev := config.GetRaw("NODE_ENV")
	require.Equal(t, "development", dev)
}

func Test_Invalid(t *testing.T) {
	_, err := config.New[Config]("doc.txt")
	require.NotNil(t, err)
}

func Test_NewEnv_NoPath(t *testing.T) {
	// When no .env path is provided, NewEnv should still work
	// by reading environment variables from OS
	type SimpleConfig struct {
		TestVar string `env:"TEST_VAR_FOR_EMPTY_PATH" default:"default_value"`
	}

	cfg, err := config.NewEnv[SimpleConfig]("")
	require.Nil(t, err)
	require.Equal(t, "default_value", cfg.TestVar)
}

func Test_Validate(t *testing.T) {
	type ConfigV struct {
		EmailAddress string `env:"NODE_ENV" validate:"required,isEmail"`
	}

	_, err := config.NewEnv[ConfigV](".env.example")
	require.NotNil(t, err)
}

func Test_Default(t *testing.T) {
	type ConfigDefault struct {
		NodeEnv string `env:"NODEENV" default:"production"`
	}

	cfg, err := config.NewEnv[ConfigDefault](".env.example")
	require.Nil(t, err)

	require.Equal(t, "production", cfg.NodeEnv)
}

func Test_Scan_ExtraTypes(t *testing.T) {
	type ConfigExtra struct {
		Count       uint      `env:"COUNT"`
		Ratio       float64   `env:"RATIO"`
		Established time.Time `env:"ESTABLISHED"`
	}

	t.Setenv("COUNT", "123")
	t.Setenv("RATIO", "3.14")
	t.Setenv("ESTABLISHED", "2023-10-26T12:00:00Z")

	var cfg ConfigExtra
	err := env.Parse(&cfg)
	require.Nil(t, err)

	require.Equal(t, uint(123), cfg.Count)
	require.Equal(t, 3.14, cfg.Ratio)

	expectedTime, _ := time.Parse(time.RFC3339, "2023-10-26T12:00:00Z")
	require.Equal(t, expectedTime, cfg.Established)
}

func Test_Scan_ErrorCases(t *testing.T) {
	type ConfigError struct {
		IntVal      int           `env:"INT_VAL"`
		UintVal     uint          `env:"UINT_VAL"`
		FloatVal    float64       `env:"FLOAT_VAL"`
		BoolVal     bool          `env:"BOOL_VAL"`
		DurationVal time.Duration `env:"DURATION_VAL"`
		TimeVal     time.Time     `env:"TIME_VAL"`
	}

	// Set invalid values
	t.Setenv("INT_VAL", "invalid_int")
	t.Setenv("UINT_VAL", "-10") // uint cannot be negative
	t.Setenv("FLOAT_VAL", "invalid_float")
	t.Setenv("BOOL_VAL", "not_a_bool")
	t.Setenv("DURATION_VAL", "invalid_duration")
	t.Setenv("TIME_VAL", "invalid_time")

	var cfg ConfigError
	// env.Parse returns error on invalid input types usually
	err := env.Parse(&cfg)
	require.Error(t, err)
}

func Test_Scan_PrivateField(t *testing.T) {
	type ConfigPrivate struct {
		Public  string `env:"PUBLIC"`
		private string `env:"PRIVATE"`
	}

	t.Setenv("PUBLIC", "public_val")
	t.Setenv("PRIVATE", "private_val")

	var cfg ConfigPrivate
	err := env.Parse(&cfg)
	require.Nil(t, err)

	require.Equal(t, "public_val", cfg.Public)
	// private field should remain empty/default
	require.Equal(t, "", cfg.private)
}

// Test_EnvOverridesYaml verifies that when a struct has both yaml and env tags,
// the env values override the yaml values after New() is called
func Test_EnvOverridesYaml(t *testing.T) {
	type HybridConfig struct {
		// This field has both yaml and env tags
		// yaml file has: node_env: development
		// env file has: NODE_ENV = 'development' (same value, but we'll set env to override)
		NodeEnv string `yaml:"node_env" env:"NODE_ENV"`

		// This field has both yaml and env tags
		// yaml file has: port: 3000
		// env file has: PORT = 5000 (different value - env should win)
		Port int `yaml:"port" env:"PORT"`

		// This field only has yaml tag - should come from yaml
		Secret string `yaml:"secret"`

		// This field only has env tag - should come from env
		Special string `env:"SPECIAL"`
	}

	// Test with yaml file first, then env parsing
	cfg, err := config.New[HybridConfig]("env.yaml")
	require.Nil(t, err)

	// Verify yaml-only field is loaded from yaml
	require.Equal(t, "secret", cfg.Secret)

	// Now set env variables and test again
	t.Setenv("NODE_ENV", "production")
	t.Setenv("PORT", "8080")
	t.Setenv("SPECIAL", "env_special")

	cfg2, err := config.New[HybridConfig]("env.yaml")
	require.Nil(t, err)

	// Verify env overrides yaml for fields with both tags
	require.Equal(t, "production", cfg2.NodeEnv) // env overrides yaml
	require.Equal(t, 8080, cfg2.Port)            // env overrides yaml

	// Verify yaml-only field is still from yaml
	require.Equal(t, "secret", cfg2.Secret)

	// Verify env-only field comes from env
	require.Equal(t, "env_special", cfg2.Special)
}

// Test_EnvSetByCode_OverridesYaml tests that environment variables set by code
// (not from .env file) properly override YAML values when using New() with a YAML file
func Test_EnvSetByCode_OverridesYaml(t *testing.T) {
	type AppConfig struct {
		// Fields with both yaml and env tags
		NodeEnv string `yaml:"node_env" env:"NODE_ENV"`
		Port    int    `yaml:"port" env:"PORT"`
		Secret  string `yaml:"secret" env:"SECRET"`

		// Nested struct with env tags
		Database struct {
			Host string `yaml:"host" env:"DB_HOST"`
			Port int    `yaml:"port" env:"DB_PORT"`
		} `yaml:"database"`
	}

	// YAML file env.yaml has:
	// node_env: development
	// port: 3000
	// secret: 'secret'
	// database:
	//   host: localhost
	//   port: 5432

	// Set environment variables by code (simulating docker/k8s env vars)
	t.Setenv("NODE_ENV", "production")    // Override yaml value "development"
	t.Setenv("PORT", "8080")              // Override yaml value 3000
	t.Setenv("SECRET", "super_secret")    // Override yaml value "secret"
	t.Setenv("DB_HOST", "db.example.com") // Override yaml value "localhost"
	t.Setenv("DB_PORT", "5433")           // Override yaml value 5432

	// Load config from YAML file - NO .env file involved
	cfg, err := config.New[AppConfig]("env.yaml")
	require.Nil(t, err)

	// Verify that code-set env vars override YAML values
	require.Equal(t, "production", cfg.NodeEnv, "NODE_ENV should be overridden by env")
	require.Equal(t, 8080, cfg.Port, "PORT should be overridden by env")
	require.Equal(t, "super_secret", cfg.Secret, "SECRET should be overridden by env")
	require.Equal(t, "db.example.com", cfg.Database.Host, "DB_HOST should be overridden by env")
	require.Equal(t, 5433, cfg.Database.Port, "DB_PORT should be overridden by env")
}

// Test_EnvSetByCode_PartialOverride tests that only some fields are overridden
// while others keep their YAML values
func Test_EnvSetByCode_PartialOverride(t *testing.T) {
	type Config struct {
		NodeEnv string `yaml:"node_env" env:"PARTIAL_NODE_ENV"`
		Port    int    `yaml:"port" env:"PARTIAL_PORT"`
		Secret  string `yaml:"secret"` // No env tag - should always come from YAML
	}

	// Set only some env vars (using unique names to avoid collision)
	t.Setenv("PARTIAL_NODE_ENV", "staging")
	// PARTIAL_PORT is NOT set - should come from YAML

	cfg, err := config.New[Config]("env.yaml")
	require.Nil(t, err)

	// NodeEnv is overridden by env var
	require.Equal(t, "staging", cfg.NodeEnv)

	// Port comes from YAML (3000) since env var is not set
	require.Equal(t, 3000, cfg.Port)

	// Secret has no env tag, so it always comes from YAML
	require.Equal(t, "secret", cfg.Secret)
}
