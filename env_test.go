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

	_, err = config.NewEnv[Config]("")
	require.NotNil(t, err)
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

// Removed Test_Scan_Fatal tests as env.Parse returns error instead of fatal panic
