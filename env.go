package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

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
				valInt, _ := strconv.Atoi(val)
				ct.Field(i).SetInt(int64(valInt))
			case "bool":
				ct.Field(i).SetBool(transform.ToBool(val))
			case "Duration":
				valDate, _ := time.ParseDuration(val)
				ct.Field(i).Set(reflect.ValueOf(valDate))
			default:
				fmt.Println(field.Type.Name())
			}
		}
	}
}
