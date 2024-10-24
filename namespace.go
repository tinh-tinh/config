package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tinh-tinh/tinhtinh/core"
)

// Namespace only available for env file

func ForRootRaw(path ...string) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		err := godotenv.Load(path...)
		if err != nil {
			log.Printf("can read env file because %s\n", err.Error())
		}
		return module
	}
}

func ForFeature[E any](name string, fncs ...func() *E) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		var value E
		if len(fncs) > 0 {
			value = *fncs[0]()
		} else {
			Scan(&value)
		}

		cfgModule := module.New(core.NewModuleOptions{})
		cfgModule.NewProvider(core.ProviderOptions{
			Name:  GetNamespace(name),
			Value: &value,
		})
		cfgModule.Export(GetNamespace(name))

		return cfgModule
	}
}

func GetNamespace(name string) core.Provide {
	return core.Provide(name)
}

func InjectNamespace[E any](module *core.DynamicModule, name string) *E {
	cfg, ok := module.Ref(GetNamespace(name)).(*E)
	if !ok {
		return nil
	}
	return cfg
}
