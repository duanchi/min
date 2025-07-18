package util

import (
	"reflect"
	"strconv"
	"strings"
)

func GetTags(key string, tag reflect.StructTag) (tags []string) {
	tagString := string(tag)
	//fmt.Println(tagString)
	idx := strings.Index(tagString, key+":\"")
	if idx == -1 {
		return
	}
	tagStack := strings.Split(tagString[idx+len(key)+2:], key+":\"")
	tags = []string{}

	for _, tagItem := range tagStack {
		for i := 0; i < len(tagItem); i++ {
			if tagItem[i] == '"' && tagItem[i-1] != '\\' {
				t, err := strconv.Unquote(tagItem[0:i])
				if err != nil {
					t = tagItem[0:i]
				}
				tags = append(tags, strings.TrimSpace(t))
				break
			}
		}
	}
	return
}
