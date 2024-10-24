package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/tinh-tinh/tinhtinh/dto/transform"
	"github.com/tinh-tinh/tinhtinh/dto/validator"
)

func NewEnv[E any](path string) (*E, error) {
	if path == "" {
		path = ".env"
	}
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	var env E
	Scan(&env)

	err = validator.Scanner(&env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func GetRaw(key string) string {
	return os.Getenv(key)
}

func Scan(env interface{}) {
	ct := reflect.ValueOf(env).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		tagVal := field.Tag.Get("mapstructure")
		if tagVal != "" {
			val := os.Getenv(tagVal)
			// Check default value
			defaultVal := field.Tag.Get("default")
			if val == "" {
				val = defaultVal
			}

			// Check empty
			if val == "" {
				continue
			}
			switch field.Type.Name() {
			case "string":
				ct.Field(i).SetString(val)
			case "int":
				ct.Field(i).SetInt(transform.StringToInt64(val))
			case "bool":
				ct.Field(i).SetBool(transform.StringToBool(val))
			case "Duration":
				ct.Field(i).Set(reflect.ValueOf(transform.StringToTimeDuration(val)))
			default:
				fmt.Println(field.Type.Name())
			}
		}
	}
}
