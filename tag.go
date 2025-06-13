package goenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Struct populates a struct with values from environment variables.
// It uses reflection to iterate through struct fields and populates them
// based on the "goenv" struct tag.
//
// The function requires a pointer to a struct as input. It will recursively
// process nested structs if they are encountered.
//
// Supported field types:
//   - string
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
//   - bool
//   - time.Duration
//   - time.Time (uses Golang's time formats)
//   - nested structs (processed recursively)
//
// Struct tag format:
//   - Use `goenv:"ENV_VAR_NAME"` to specify the environment variable name
//   - Fields without the goenv tag are ignored
//   - Unexported fields are skipped automatically
//
// Example:
//
//	type DatabaseConfig struct {
//		Host     string `goenv:"DB_HOST"`
//		Port     int    `goenv:"DB_PORT"`
//		Username string `goenv:"DB_USER"`
//		Password string `goenv:"DB_PASS"`
//		SSL      bool   `goenv:"DB_SSL"`
//	}
//
//	var dbConfig DatabaseConfig
//	err := goenv.Struct(&dbConfig)
//	if err != nil {
//		return fmt.Errorf("failed to load database config: %w", err)
//	}
func Struct(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Pointer {
		return fmt.Errorf("value must be a pointer to a struct")
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("value must be a pointer to a struct")
	}

	for i := range val.NumField() {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		if field.Kind() == reflect.Struct && field.Type() != reflect.TypeOf(time.Time{}) {
			if err := Struct(field.Addr().Interface()); err != nil {
				return fmt.Errorf("failed to process nested struct %s: %w", val.Type().Field(i).Name, err)
			}
			continue
		}

		tag := val.Type().Field(i).Tag.Get("goenv")
		if tag == "" {
			continue
		}

		env, found := os.LookupEnv(tag)
		// TODO: support fallback values in future or if field is required
		if !found {
			return fmt.Errorf("environment variable %s is not defined", tag)
		}

		if err := setFieldValue(field, env); err != nil {
			return fmt.Errorf("failed to set %s environment variable: %w", tag, err)
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	if field.Type() == reflect.TypeOf(time.Time{}) {
		fmt.Printf("time value is: %s\n", value)
		timeVal, err := parseTimeValue(value)
		if err != nil {
			return fmt.Errorf("cannot parse %q as time.Time: %w", value, err)
		}
		field.Set(reflect.ValueOf(timeVal))
		return nil
	}

	switch field.Interface().(type) {
	case string:
		field.SetString(value)
	case int, int8, int16, int32, int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse %q as int: %w", value, err)
		}
		field.SetInt(intVal)
	case uint, uint8, uint16, uint32, uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse %q as uint: %w", value, err)
		}
		field.SetUint(uintVal)
	case float32, float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("cannot parse %q as float: %w", value, err)
		}
		field.SetFloat(floatVal)
	case bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot parse %q as bool: %w", value, err)
		}
		field.SetBool(boolVal)
	case time.Duration:
		durVal, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("cannot parse %q as time.Duration: %w", value, err)
		}
		field.Set(reflect.ValueOf(durVal))
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

func parseTimeValue(value string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unknown go format")
}
