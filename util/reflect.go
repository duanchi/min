package util

import (
	"reflect"
	"strconv"
	"strings"
)

func GetTags(key string, tag reflect.StructTag, fullMatch ...bool) (tags []string) {
	tagString := string(tag)
	idx := -1
	matchKey := key
	if val, _ := GetOptionalParameter(fullMatch); !val {
		matchKey = "@" + key
		idx = strings.Index(tagString, matchKey+":\"")
	}
	if idx == -1 {
		matchKey = key
		idx = strings.Index(tagString, matchKey+":\"")
		if idx == -1 {
			return
		}
	}

	tagStack := strings.Split(tagString[idx+len(matchKey)+2:], matchKey+":\"")
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
