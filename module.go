package config

import (
	"log"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/tinh-tinh/tinhtinh/core"
)

const ENV core.Provide = "ConfigEnv"

type Options[E any] struct {
	EnvPath       string
	IgnoreEnvFile bool
	// Only for env file
	Load func() *E
}

type Param[V any] interface {
	string | Options[V]
}

func ForRoot[E any, param Param[E]](params ...param) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		var lastValue *E
		var err error

		if len(params) == 0 {
			lastValue, err = New[E]("")
			if err != nil {
				log.Println("env not found")
			}
		} else {
			for _, v := range params {
				if reflect.TypeOf(v).Kind() == reflect.String {
					lastValue, err = New[E](any(v).(string))
					if err != nil {
						continue
					}
				} else if reflect.TypeOf(v).Kind() == reflect.Struct {
					opt := any(v).(Options[E])
					err = godotenv.Load(opt.EnvPath)
					if err != nil {
						continue
					}
					lastValue = opt.Load()
				}
			}
		}

		configModule := module.New(core.NewModuleOptions{})

		if lastValue == nil {
			log.Println("env not found")
			return configModule
		}

		configModule.NewProvider(core.ProviderOptions{
			Name:  ENV,
			Value: lastValue,
		})
		configModule.Export(ENV)

		return configModule
	}
}
