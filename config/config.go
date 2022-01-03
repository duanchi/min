package config

import (
	"fmt"
	"github.com/duanchi/min/config/yaml"
	"github.com/duanchi/min/util"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ConfigInstance interface{}

func Init(config interface{}) {
	fmt.Println("[min-framework] Init Config...")
	err := yaml.GetConfig(config)

	if err == nil {
		parseConfig(config, "")
		ConfigInstance = config
	} else {
		panic(err)
	}
}

func Get(key string) interface{} {
	return GetRaw(key).Interface()
}

func GetRaw(key string) reflect.Value {

	keyStack := strings.Split(key, ".")
	value := reflect.ValueOf(ConfigInstance)

	if value.IsValid() {
		value = value.Elem()
	} else {
		return value
	}

	for i := 0; i < len(keyStack); i++ {

		//fmt.Printf("key: %s, kind %s",keyStack[i], reflect.TypeOf(value).Kind())

		// 调用栈不是末尾, 并且value是可用值, 并且value是基础类型
		if i < len(keyStack)-1 && value.IsValid() && value.Kind() != reflect.Ptr && value.Kind() != reflect.Struct {
			return reflect.New(value.Type())
		} else {
			if value.Kind() == reflect.Struct {
				if value.FieldByName(keyStack[i]).IsValid() {
					value = value.FieldByName(keyStack[i])
				} else {
					value = reflect.New(value.Type())
				}
			} else if value.Kind() == reflect.Ptr {
				if value.Elem().FieldByName(keyStack[i]).IsValid() {
					value = value.Elem().FieldByName(keyStack[i])
				} else {
					value = reflect.New(value.Elem().Type())
				}
			} else {
				if value.Elem().FieldByName(keyStack[i]).IsZero() || value.Elem().FieldByName(keyStack[i]).IsNil() {
					value = reflect.New(value.Elem().Type())
				} else {
					value = value.Elem().FieldByName(keyStack[i])
				}
			}
		}
	}

	return value
}

func SetConfigFile(configFile string) {
	yaml.SetConfigFile(configFile)
}

func parseConfig(config interface{}, defaults string) {

	configType := reflect.TypeOf(config)
	configValue := reflect.ValueOf(config)

	if reflect.TypeOf(config).Kind() == reflect.Ptr {
		configType = configType.Elem()
		configValue = configValue.Elem()
	}

	switch configType.Kind() {
	case reflect.Struct:
		for i := 0; i < configValue.NumField(); i++ {
			if configValue.Field(i).CanInterface() {
				switch configValue.Field(i).Type().Kind() {
				case reflect.Ptr, reflect.Struct, reflect.Slice:
					parseConfig(configValue.Field(i).Addr().Interface(), "")

				case reflect.Map:
					for _, key := range configValue.Field(i).MapKeys() {
						parseConfig(configValue.Field(i).MapIndex(key).Interface(), "")
					}

				default:
					if configValue.Field(i).CanAddr() {
						parseConfig(configValue.Field(i).Addr().Interface(), configType.Field(i).Tag.Get("default"))
					} else {
						parseConfig(configValue.Field(i).Interface(), configType.Field(i).Tag.Get("default"))
					}

				}
			}
		}
	case reflect.Map:
		for _, key := range configValue.MapKeys() {
			parseConfig(configValue.MapIndex(key).Addr().Interface(), "")
		}
	case reflect.Slice:
		for index := 0; index < configValue.Len(); index++ {
			parseConfig(configValue.Index(index).Addr().Interface(), "")
		}

	case reflect.Ptr:
		parseConfig(configValue.Elem(), "")

	default:

		if configValue.CanSet() {
			v := ""

			if configValue.IsZero() {
				v = defaults

				pattern, _ := regexp.Compile(`\${.+?}`)
				v = pattern.ReplaceAllStringFunc(defaults, func(s string) string {
					value := strings.SplitN(s[2:len(s)-1], ":", 2)
					if len(value) > 1 {
						return util.Getenv(value[0], value[1])
					} else {
						return util.Getenv(value[0], "")
					}
				})

				class := configType.Kind()
				switch class {
				case reflect.String:
					configValue.SetString(v)

				case reflect.Int, reflect.Int64:
					value, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						configValue.SetInt(0)
					} else {
						configValue.SetInt(value)
					}

				case reflect.Bool:
					value, err := strconv.ParseBool(v)
					if err != nil {
						configValue.SetBool(false)
					} else {
						configValue.SetBool(value)
					}

				case reflect.Float64:
					value, err := strconv.ParseFloat(v, 10)
					if err != nil {
						configValue.SetFloat(0)
					} else {
						configValue.SetFloat(value)
					}
				}
			}
		}
	}
}