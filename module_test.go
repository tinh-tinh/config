package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/config/v2"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_ForRootNil(t *testing.T) {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			config.ForRoot[Config, string](),
		},
	})

	cfg, ok := appModule.Ref(config.ENV).(*Config)
	require.False(t, ok)
	require.Nil(t, cfg)
}

func Test_ForRoot(t *testing.T) {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			config.ForRoot[Config](".env.example", ".env.local"),
		},
	})

	cfg, ok := appModule.Ref(config.ENV).(*Config)
	require.True(t, ok)
	require.NotNil(t, cfg)
	require.Equal(t, "development", cfg.NodeEnv)
}

func Test_LoadConfig(t *testing.T) {
	type Database struct {
		Host string
		Port string
	}
	type Cfg struct {
		NodeEnv  string
		Database Database
	}

	load := func() *Cfg {
		return &Cfg{
			NodeEnv: os.Getenv("NODE_ENV"),
			Database: Database{
				Host: os.Getenv("DB_HOST"),
				Port: os.Getenv("DB_PORT"),
			},
		}
	}

	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			config.ForRoot[Cfg](config.Options[Cfg]{
				EnvPath: ".env.example",
				Load:    load,
			}),
		},
	})

	cfg := config.Inject[Cfg](appModule)
	require.NotNil(t, cfg)
	require.Equal(t, "development", cfg.NodeEnv)
	require.Equal(t, "localhost", cfg.Database.Host)
	require.Equal(t, "5432", cfg.Database.Port)
}

func Test_Yaml(t *testing.T) {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			config.ForRoot[ConfigYaml]("env.yaml"),
		},
	})

	cfg, ok := appModule.Ref(config.ENV).(*ConfigYaml)
	require.True(t, ok)

	require.Equal(t, "development", cfg.NodeEnv)
	require.Equal(t, 3000, cfg.Port)
	require.Equal(t, 5*time.Minute, cfg.ExpiresIn)
	require.Equal(t, true, cfg.Log)
	require.Equal(t, "secret", cfg.Secret)
}

func Test_Hybrid(t *testing.T) {
	type HybridEnv struct {
		NodeEnv   string        `yaml:"node_env"`
		Port      int           `yaml:"port"`
		ExpiresIn time.Duration `yaml:"expires_in"`
		Log       bool          `yaml:"log"`
		Special   string        `mapstructure:"SPECIAL"`
	}
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			config.ForRoot[HybridEnv]("env.yaml", ".env.example"),
		},
	})

	cfg, ok := appModule.Ref(config.ENV).(*HybridEnv)
	require.True(t, ok)
	require.NotNil(t, cfg)

	// Test that values from both YAML and env files are properly merged
	require.Equal(t, "development", cfg.NodeEnv) // From YAML file
	require.Equal(t, 3000, cfg.Port)             // From YAML file
	require.True(t, cfg.Log)                     // From YAML file
	require.NotEmpty(t, cfg.Special)             // From env file (if SPECIAL is set)

	// Optional: Print for debugging
	t.Logf("Merged config: %+v", cfg)
}
