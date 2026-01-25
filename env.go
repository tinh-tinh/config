package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

var (
	typeDuration = reflect.TypeOf(time.Duration(0))
	typeTime     = reflect.TypeOf(time.Time{})
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

	v := validator.Validator{}
	err = v.Validate(&env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func GetRaw(key string) string {
	return os.Getenv(key)
}

func Scan(env any) {
	if env == nil {
		log.Fatal("env must not be nil")
	}

	v := reflect.ValueOf(env)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		log.Fatal("env must be a non-nil pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		log.Fatal("env must point to a struct")
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if fieldType.PkgPath != "" {
			continue
		}

		tag := fieldType.Tag.Get("env")
		if tag == "" || tag == "-" {
			continue
		}

		// Handle tags like: env:"FOO,omitempty"
		key := strings.Split(tag, ",")[0]
		if key == "" {
			continue
		}

		val, exists := os.LookupEnv(key)

		// Apply default only if env var is missing
		if !exists {
			if def := fieldType.Tag.Get("default"); def != "" {
				val = def
			}
		}

		if val == "" {
			continue
		}

		if err := setField(fieldVal, fieldType, val); err != nil {
			log.Printf("error setting field %s: %s", fieldType.Name, err.Error())
		}
	}

}

func setField(fieldVal reflect.Value, fieldStruct reflect.StructField, val string) error {
	if !fieldVal.IsValid() {
		return nil
	}

	// Handle pointers (allocate if nil)
	if fieldVal.Kind() == reflect.Pointer {
		if fieldVal.IsNil() {
			fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
		}
		return setField(fieldVal.Elem(), fieldStruct, val)
	}

	if !fieldVal.CanSet() {
		return nil
	}

	// Handle concrete well-known types first
	switch fieldVal.Type() {
	case typeDuration:
		d, err := time.ParseDuration(val)
		if err != nil {
			return fmt.Errorf("field %s: invalid duration: %w", fieldStruct.Name, err)
		}
		fieldVal.Set(reflect.ValueOf(d))
		return nil

	case typeTime:
		// RFC3339 is the safest default
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return fmt.Errorf("field %s: invalid time: %w", fieldStruct.Name, err)
		}
		fieldVal.Set(reflect.ValueOf(t))
		return nil
	}

	// Handle by kind (works for aliases like `type Age int`)
	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(val)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bits := fieldVal.Type().Bits()
		i, err := strconv.ParseInt(val, 10, bits)
		if err != nil {
			return fmt.Errorf("field %s: invalid int: %w", fieldStruct.Name, err)
		}
		fieldVal.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bits := fieldVal.Type().Bits()
		u, err := strconv.ParseUint(val, 10, bits)
		if err != nil {
			return fmt.Errorf("field %s: invalid uint: %w", fieldStruct.Name, err)
		}
		fieldVal.SetUint(u)

	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("field %s: invalid bool: %w", fieldStruct.Name, err)
		}
		fieldVal.SetBool(b)

	case reflect.Float32, reflect.Float64:
		bits := fieldVal.Type().Bits()
		f, err := strconv.ParseFloat(val, bits)
		if err != nil {
			return fmt.Errorf("field %s: invalid float: %w", fieldStruct.Name, err)
		}
		fieldVal.SetFloat(f)

	default:
		return fmt.Errorf(
			"field %s: unsupported type %s",
			fieldStruct.Name,
			fieldVal.Type().String(),
		)
	}

	return nil
}
