package config

import (
	"os"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type RegisterWhenOptions interface {
	string | func() bool
}

func RegisterWhen[opt RegisterWhenOptions](module core.Modules, env opt) core.Modules {
	condition := false

	if reflect.TypeOf(env).Kind() == reflect.String {
		condition = os.Getenv(any(env).(string)) != ""
	} else if reflect.TypeOf(env).Kind() == reflect.Func {
		fnc := any(env).(func() bool)
		condition = fnc()
	}

	if condition {
		return module
	}

	return nil
}
