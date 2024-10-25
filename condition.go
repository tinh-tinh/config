package config

import (
	"os"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/core"
)

type RegisterWhenOptions interface {
	string | func() bool
}

func RegisterWhen[opt RegisterWhenOptions](module core.Module, env opt) core.Module {
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
