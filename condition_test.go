package config_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/config/v2"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Condition(t *testing.T) {
	err := godotenv.Load(".env.example")
	require.Nil(t, err)

	const USER core.Provide = "TinhTinh"
	userModule := func(module core.Module) core.Module {
		mod := module.New(core.NewModuleOptions{})

		mod.NewProvider(core.ProviderOptions{
			Name:  USER,
			Value: "TinhTinh",
		})
		mod.Export(USER)

		return mod
	}

	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			config.RegisterWhen(userModule, "NODE_ENV"),
		},
	})

	value := module.Ref(USER)
	require.Equal(t, "TinhTinh", value)
}

func Test_ConditionFailed(t *testing.T) {
	err := godotenv.Load(".env.example")
	require.Nil(t, err)

	const USER core.Provide = "TinhTinh"
	userModule := func(module core.Module) core.Module {
		mod := module.New(core.NewModuleOptions{})

		mod.NewProvider(core.ProviderOptions{
			Name:  USER,
			Value: "TinhTinh",
		})
		mod.Export(USER)

		return mod
	}

	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			config.RegisterWhen(userModule, "HAHA"),
		},
	})

	value := module.Ref(USER)
	require.Nil(t, value)
}

func Test_ConditionFnc(t *testing.T) {
	err := godotenv.Load(".env.example")
	require.Nil(t, err)

	const USER core.Provide = "TinhTinh"
	userModule := func(module core.Module) core.Module {
		mod := module.New(core.NewModuleOptions{})

		mod.NewProvider(core.ProviderOptions{
			Name:  USER,
			Value: "TinhTinh",
		})
		mod.Export(USER)

		return mod
	}

	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			config.RegisterWhen(userModule, func() bool {
				return os.Getenv("NODE_ENV") == "development"
			}),
		},
	})

	value := module.Ref(USER)
	require.Equal(t, "TinhTinh", value)
}

func Test_ConditionFncFailed(t *testing.T) {
	err := godotenv.Load(".env.example")
	require.Nil(t, err)

	const USER core.Provide = "TinhTinh"
	userModule := func(module core.Module) core.Module {
		mod := module.New(core.NewModuleOptions{})

		mod.NewProvider(core.ProviderOptions{
			Name:  USER,
			Value: "TinhTinh",
		})
		mod.Export(USER)

		return mod
	}

	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			config.RegisterWhen(userModule, func() bool {
				return os.Getenv("NODE_ENV") == "HAHA"
			}),
		},
	})

	value := module.Ref(USER)
	require.Nil(t, value)
}
