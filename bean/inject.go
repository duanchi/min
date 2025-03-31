package bean

import (
	"errors"
	"fmt"
	"github.com/duanchi/min/v2/config"
	"github.com/duanchi/min/v2/util"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func Inject(rawBean reflect.Value, beanName string, beanMap map[string]reflect.Value) {

	beanType := rawBean.Type()

	for i := 0; i < beanType.NumField(); i++ {
		if rawBean.Field(i).CanSet() {
			fieldTag := beanType.Field(i).Tag

			if err := ParseValueFromConfig(fieldTag.Get("value"), rawBean.Field(i)); err != nil {
				fmt.Println("[min-framework] Bean " + beanName + " Injected error, " + err.Error())
			}
			if util.IsBeanKind(fieldTag, "autowired") || util.IsBeanKind(fieldTag, "resource") {
				parseTagNamedAutowired(rawBean.Field(i))
			}
		}
	}
}

func getCurrentGoroutineStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func ParseValueFromConfig(value string, field reflect.Value) error {
	if value != "" {
		class := field.Kind()
		regex, _ := regexp.Compile("^" + regexp.QuoteMeta("${") + "(.+)" + regexp.QuoteMeta("}") + "$")

		if regex.MatchString(value) {
			value = string(regex.ReplaceAllFunc([]byte(value), func(match []byte) []byte {
				return match[2 : len(match)-1]
			})[:])

			configField := strings.Split(value, ",")
			configValue := config.GetRaw(configField[0])

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
				if configValue.Type() == field.Type() {
					field.Set(configValue)
				} else {
					return errors.New("value " + value + " error while parsing config.")
				}

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

	return nil
}

func parseTagNamedAutowired(field reflect.Value) {
	beanType := field.Type()
	if beanType.Kind() == reflect.Ptr {
		beanPointer, ok := beanTypeMaps[beanType]
		if ok {
			field.Set(beanPointer)
		}
	}
}
