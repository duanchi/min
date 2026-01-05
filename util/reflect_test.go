package util

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type Method struct {
	Name string `@json:"name\"can dllkejlk\"" @json:"name2"`
	Kind string `json:" kind"    ste:""  l `
}

func TestTag(t *testing.T) {
	k := strings.Split("aaaaaaaaaaaaarestful:aaaaaaaaaaaaa", "restful:")
	fmt.Println(k[0], k[1])
	b := reflect.TypeOf(Method{})
	field, _ := b.FieldByName("Name")
	fmt.Println(GetTags("json", field.Tag))
}
