package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var allowedBool = map[string]bool{
	"1":     true,
	"0":     false,
	"true":  true,
	"false": false,
	"True":  true,
	"False": false,
	"":      false,
}

func allowedBoolKeys() []string {
	res := make([]string, 0)
	for key := range allowedBool {
		res = append(res, key)
	}
	return res
}

//osEnvToMap transforms os.Environ() to map
func osEnvToMap() map[string]string {
	res := make(map[string]string, len(os.Environ()))
	for _, kv := range os.Environ() {
		splitted := strings.Split(kv, "=")
		res[splitted[0]] = splitted[1]
	}
	return res
}

func readEnv(field interface{}, envKey string) error {
	envVal, exists := osEnvToMap()[envKey]
	if !exists {
		return nil
	}
	value := reflect.ValueOf(field).Elem()
	switch value.Kind() {
	case reflect.String:
		value.SetString(envVal)
	case reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Int:
		v, err := strconv.ParseUint(envVal, 10, 16)
		if err != nil {
			return err
		}
		value.SetUint(v)
	case reflect.Bool:
		boolVal, exists := allowedBool[envVal]
		if !exists {
			return fmt.Errorf("not supported bool value. should be one of: %s", strings.Join(allowedBoolKeys(), ", "))
		}
		value.SetBool(boolVal)
	default:
		return fmt.Errorf("not support type of struct: %s", value.Kind())
	}
	return nil
}

func mapMapToStruct(target interface{}, valuesMap map[string]string) error {
	fields := reflect.TypeOf(target).Elem()
	values := reflect.ValueOf(target).Elem()
	num := fields.NumField()
	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		tag := field.Tag.Get("env")
		envVal, exists := valuesMap[tag]
		if !exists {
			continue
		}
		switch value.Kind() {
		case reflect.String:
			values.FieldByName(field.Name).SetString(envVal)
		case reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Int:
			v, err := strconv.ParseUint(envVal, 10, 16)
			if err != nil {
				return err
			}
			value.SetUint(v)
		case reflect.Bool:
			boolVal, exists := allowedBool[envVal]
			if !exists {
				return fmt.Errorf("not supported bool value. should be one of: %s", strings.Join(allowedBoolKeys(), ", "))
			}
			value.SetBool(boolVal)
		default:
			return fmt.Errorf("not support type of struct: %s", value.Kind())
		}
	}
	return nil
}

//mapEnvToStruct Maps os.Env to target struct corresponding to `env` tag
func mapEnvToStruct(target interface{}) error {
	err := mapMapToStruct(target, osEnvToMap())
	return err
}

//InitStructFromEnv Initialize target struct from env with corresponding to `env` tag
func InitStructFromEnv(target interface{}) error {
	err := mapEnvToStruct(target)
	if err != nil {
		return err
	}
	return nil
}
