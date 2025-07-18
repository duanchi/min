package util

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func ParseValueFromConfigInstance(value string, field reflect.Value, configInstance interface{}) {
	if value != "" {
		class := field.Kind()
		regex, _ := regexp.Compile("^" + regexp.QuoteMeta("${") + "(.+)" + regexp.QuoteMeta("}") + "$")

		if regex.MatchString(value) {
			value = string(regex.ReplaceAllFunc([]byte(value), func(match []byte) []byte {
				return match[2 : len(match)-1]
			})[:])

			configField := strings.Split(value, ",")
			configValue := GetRawConfigFromInstance(configField[0], configInstance)

			if configValue.IsZero() && len(configField) > 1 {

				switch class {
				case reflect.String:
					field.SetString(configField[1])

				case reflect.Int, reflect.Int64:
					value, err := strconv.ParseInt(configField[1], 10, 64)
					if err != nil {
						field.SetInt(0)
					} else {
						field.SetInt(value)
					}

				case reflect.Bool:
					value, err := strconv.ParseBool(configField[1])
					if err != nil {
						field.SetBool(false)
					} else {
						field.SetBool(value)
					}

				case reflect.Float64:
					value, err := strconv.ParseFloat(configField[1], 10)
					if err != nil {
						field.SetFloat(0)
					} else {
						field.SetFloat(value)
					}
				}
			} else {
				field.Set(configValue)
			}
		} else {

			switch class {
			case reflect.String:
				field.SetString(value)

			case reflect.Int, reflect.Int64:
				value, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					field.SetInt(0)
				} else {
					field.SetInt(value)
				}

			case reflect.Bool:
				value, err := strconv.ParseBool(value)
				if err != nil {
					field.SetBool(false)
				} else {
					field.SetBool(value)
				}

			case reflect.Float64:
				value, err := strconv.ParseFloat(value, 10)
				if err != nil {
					field.SetFloat(0)
				} else {
					field.SetFloat(value)
				}
			}
		}
	}
}

func GetRawConfigFromInstance(key string, configInstance interface{}) reflect.Value {

	keyStack := strings.Split(key, ".")
	value := reflect.ValueOf(configInstance)

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

func ParseTag(name string, tag reflect.StructTag, defaultKey ...string) (tagMapList []map[string]string) {
	tags := GetTags(name, tag)
	tagMapList = []map[string]string{}
	key := "name"

	if len(defaultKey) > 0 && defaultKey[0] != "" {
		key = defaultKey[0]
	}

	for _, tagString := range tags {
		tagStack := strings.Split(tagString, ",")

		tagMap := map[string]string{}
		for _, tagItem := range tagStack {
			colons := strings.SplitN(strings.TrimSpace(tagItem), ":", 2)
			fmt.Println(colons)
			if len(colons) < 2 {
				tagMap[key] = colons[0]
			} else {
				tagMap[colons[0]] = colons[1]
			}
		}
		tagMapList = append(tagMapList, tagMap)
	}

	return
}
