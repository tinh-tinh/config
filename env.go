package config

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
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
			switch field.Type.Kind() {
			case reflect.String:
				ct.Field(i).SetString(val)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if field.Type == reflect.TypeOf(time.Duration(0)) {
					valDate, err := time.ParseDuration(val)
					if err != nil {
						log.Default().Printf("Error parsing duration: %v for field %s\n", err, field.Name)
						continue
					}
					ct.Field(i).Set(reflect.ValueOf(valDate))
				} else {
					valInt, err := strconv.ParseInt(val, 10, 64)
					if err != nil {
						log.Default().Printf("Error parsing int: %v for field %s\n", err, field.Name)
						continue
					}
					ct.Field(i).SetInt(valInt)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				valUint, err := strconv.ParseUint(val, 10, 64)
				if err != nil {
					log.Default().Printf("Error parsing uint: %v for field %s\n", err, field.Name)
					continue
				}
				ct.Field(i).SetUint(valUint)
			case reflect.Bool:
				valBool, err := strconv.ParseBool(val)
				if err != nil {
					log.Default().Printf("Error parsing bool: %v for field %s\n", err, field.Name)
					continue
				}
				ct.Field(i).SetBool(valBool)
			default:
				log.Default().Println("Unsupported type: ", field.Type.Name())
			}
		}
	}
}
