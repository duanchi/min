package yaml

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)


func getRaw(key string) reflect.Value {

	keyStack := strings.Split(key, ".")
	value := reflect.ValueOf(configInstance).Elem()

	for i := 0; i < len(keyStack); i++ {

		//fmt.Printf("key: %s, kind %s",keyStack[i], reflect.TypeOf(value).Kind())

		// 调用栈不是末尾, 并且value是可用值, 并且value是基础类型
		if i < len(keyStack) - 1 && value.IsValid() && value.Kind() != reflect.Ptr && value.Kind() != reflect.Struct {
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

func parseConfig (config interface{}, configMap interface{}) {
	configType := reflect.TypeOf(config).Elem()
	configValue := reflect.ValueOf(config).Elem()
	configRaw := reflect.ValueOf(configMap)
	for i := 0; i < configValue.NumField(); i++ {
		configKey := configType.Field(i).Tag.Get("yaml")
		if configKey == "" {
			configKey = configType.Field(i).Name
		}
		if configValue.Field(i).CanInterface() && configValue.Field(i).CanSet() {
			if configValue.Field(i).Kind() == reflect.Struct {
				if configRaw.IsValid() {
					if configKey == ",inline" {
						parseConfig(configValue.Field(i).Addr().Interface(),  configRaw.Interface())
					} else {
						if configRaw.MapIndex(reflect.ValueOf(configKey)).IsValid() {
							parseConfig(configValue.Field(i).Addr().Interface(), configRaw.MapIndex(reflect.ValueOf(configKey)).Interface())
						} else {
							parseConfig(configValue.Field(i).Addr().Interface(),  nil)
						}
					}
				} else {
					parseConfig(configValue.Field(i).Addr().Interface(),  nil)
				}


			} else if configValue.Field(i).Kind() == reflect.Slice {
				fmt.Println(configType.Field(i).Type)
				fmt.Println(configRaw.Kind())
				if configRaw.Kind() == reflect.Slice || configRaw.Kind() == reflect.Map {

					for j:=0; j < configRaw.Len(); j++ {
						// reflect.Append(configValue.Field(i), reflect.New(configType.Field(i).Type.))
						parseConfig(configValue.Field(i).Index(j).Addr().Interface(), configRaw.Index(i))
						// parseConfig(configValue.Field(i).Index(j).Addr().Interface(), configRaw.Index(i))
					}

					fmt.Printf("len - %d", configValue.Field(i).Len())
				}

			} else {
				v := ""
				defaultValue := configType.Field(i).Tag.Get("default")
				yamlValue := configValue.Field(i).String()
				//未配置则取默认值
				if yamlValue == "" && defaultValue != "" {
					// configValue.Field(i).SetString(defaultValue)
					v = defaultValue
				} else {
					v = yamlValue
				}
				envValue := ""
				envKey := ""
				envDefaultValue := ""
				if strings.Index(yamlValue, "$") != -1 {
					start := strings.Index(yamlValue, "{")
					end := strings.LastIndex(yamlValue, "}")
					elContent := yamlValue[start:end]
					index := strings.Index(elContent, ":")
					if index != -1 {
						envKey = elContent[1:index]
						envDefaultValue = elContent[index+1:]
					} else {//不存在配置值
						envKey = elContent[1:]
					}

					envValue = os.Getenv(envKey)
					if envValue != "" {
						log.Println(envKey + ": " + envValue)
					}
					//fmt.Println(envKey + ": " + envValue)
					//环境变量不存在则获取缺省值
					if envValue == "" && envDefaultValue != "" {
						envValue = envDefaultValue
					}
					if envValue != "" {
						// configValue.Field(i).SetString(envValue)
						v = yamlValue[0:start - 1] + envValue + yamlValue[end:len(yamlValue) - 1]
					}
				}


				class := configType.Field(i).Type.Kind()
				switch class {
				case reflect.String:
					configValue.Field(i).SetString(v)

				case reflect.Int, reflect.Int64:
					value, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						configValue.Field(i).SetInt(0)
					} else {
						configValue.Field(i).SetInt(value)
					}

				case reflect.Bool:
					value, err := strconv.ParseBool(v)
					if err != nil {
						configValue.Field(i).SetBool(false)
					} else {
						configValue.Field(i).SetBool(value)
					}

				case reflect.Float64:
					value, err := strconv.ParseFloat(v, 10)
					if err != nil {
						configValue.Field(i).SetFloat(0)
					} else {
						configValue.Field(i).SetFloat(value)
					}
				}
			}
		}
	}
}
