package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/config/v2"
)

type ConfigYaml struct {
	NodeEnv   string        `yaml:"node_env"`
	Port      int           `yaml:"port"`
	ExpiresIn time.Duration `yaml:"expires_in"`
	Log       bool          `yaml:"log"`
	Secret    interface{}   `yaml:"secret"`
}

func Test_New_Yml(t *testing.T) {
	cfg, err := config.NewYaml[ConfigYaml]("env.yaml")
	require.Nil(t, err)

	require.Equal(t, "development", cfg.NodeEnv)
	require.Equal(t, 3000, cfg.Port)
	require.Equal(t, 5*time.Minute, cfg.ExpiresIn)
	require.Equal(t, true, cfg.Log)
	require.Equal(t, "secret", cfg.Secret)

	_, err = config.NewYaml[ConfigYaml]("")
	require.NotNil(t, err)

	type Cfg struct {
		NodeEnv   time.Duration `yaml:"node_env"`
		Port      string        `yaml:"port"`
		ExpiresIn bool          `yaml:"expires_in"`
		Log       float32       `yaml:"log"`
		Secret    int           `yaml:"secret"`
	}

	_, err = config.NewYaml[Cfg]("env.yaml")
	require.NotNil(t, err)
}
