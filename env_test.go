package config_test

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/config"
)

type Config struct {
	NodeEnv   string        `mapstructure:"NODE_ENV"`
	Port      int           `mapstructure:"PORT"`
	ExpiresIn time.Duration `mapstructure:"EXPIRES_IN"`
	Log       bool          `mapstructure:"LOG"`
	Special   interface{}   `mapstructure:"SPECIAL"`
	Secret    string        `mapstructure:"SECRET"`
}

func Test_Scan(t *testing.T) {
	err := godotenv.Load(".env.example")
	require.Nil(t, err)
	var cfg Config
	config.Scan(&cfg)

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
		EmailAddress string `mapstructure:"NODE_ENV" validate:"required,isEmail"`
	}

	_, err := config.NewEnv[ConfigV](".env.example")
	require.NotNil(t, err)
}

func Test_Default(t *testing.T) {
	type ConfigDefault struct {
		NodeEnv string `mapstructure:"NODEENV" default:"production"`
	}

	cfg, err := config.NewEnv[ConfigDefault](".env.example")
	require.Nil(t, err)

	require.Equal(t, "production", cfg.NodeEnv)
}
