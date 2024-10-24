package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type ConfigYaml struct {
	NodeEnv   string        `yaml:"node_env"`
	Port      int           `yaml:"port"`
	ExpiresIn time.Duration `yaml:"expires_in"`
	Log       bool          `yaml:"log"`
	Secret    interface{}   `yaml:"secret"`
}

func Test_New_Yml(t *testing.T) {
	cfg, err := NewYaml[ConfigYaml]("env.yaml")
	require.Nil(t, err)

	require.Equal(t, "development", cfg.NodeEnv)
	require.Equal(t, 3000, cfg.Port)
	require.Equal(t, 5*time.Minute, cfg.ExpiresIn)
	require.Equal(t, true, cfg.Log)
	require.Equal(t, "secret", cfg.Secret)

	_, err = NewYaml[ConfigYaml]("")
	require.NotNil(t, err)

	type Cfg struct {
		NodeEnv   time.Duration `yaml:"node_env"`
		Port      string        `yaml:"port"`
		ExpiresIn bool          `yaml:"expires_in"`
		Log       float32       `yaml:"log"`
		Secret    int           `yaml:"secret"`
	}

	_, err = NewYaml[Cfg]("env.yaml")
	require.NotNil(t, err)
}
