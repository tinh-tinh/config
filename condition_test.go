package config

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Condition(t *testing.T) {
	godotenv.Load(".env.example")

	const USER core.Provide = "TinhTinh"
	userModule := func(module *core.DynamicModule) *core.DynamicModule {
		mod := module.New(core.NewModuleOptions{})

		mod.NewProvider(core.ProviderOptions{
			Name:  USER,
			Value: "TinhTinh",
		})
		mod.Export(USER)

		return mod
	}

	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{
			RegisterWhen(userModule, "NODE_ENV"),
		},
	})

	value := module.Ref(USER)
	require.Equal(t, "TinhTinh", value)
}

func Test_ConditionFailed(t *testing.T) {
	godotenv.Load(".env.example")

	const USER core.Provide = "TinhTinh"
	userModule := func(module *core.DynamicModule) *core.DynamicModule {
		mod := module.New(core.NewModuleOptions{})

		mod.NewProvider(core.ProviderOptions{
			Name:  USER,
			Value: "TinhTinh",
		})
		mod.Export(USER)

		return mod
	}

	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{
			RegisterWhen(userModule, "HAHA"),
		},
	})

	value := module.Ref(USER)
	require.Nil(t, value)
}

func Test_ConditionFnc(t *testing.T) {
	godotenv.Load(".env.example")

	const USER core.Provide = "TinhTinh"
	userModule := func(module *core.DynamicModule) *core.DynamicModule {
		mod := module.New(core.NewModuleOptions{})

		mod.NewProvider(core.ProviderOptions{
			Name:  USER,
			Value: "TinhTinh",
		})
		mod.Export(USER)

		return mod
	}

	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{
			RegisterWhen(userModule, func() bool {
				return os.Getenv("NODE_ENV") == "development"
			}),
		},
	})

	value := module.Ref(USER)
	require.Equal(t, "TinhTinh", value)
}

func Test_ConditionFncFailed(t *testing.T) {
	godotenv.Load(".env.example")

	const USER core.Provide = "TinhTinh"
	userModule := func(module *core.DynamicModule) *core.DynamicModule {
		mod := module.New(core.NewModuleOptions{})

		mod.NewProvider(core.ProviderOptions{
			Name:  USER,
			Value: "TinhTinh",
		})
		mod.Export(USER)

		return mod
	}

	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{
			RegisterWhen(userModule, func() bool {
				return os.Getenv("NODE_ENV") == "HAHA"
			}),
		},
	})

	value := module.Ref(USER)
	require.Nil(t, value)
}
