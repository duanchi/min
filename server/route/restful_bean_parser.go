package route

import (
	"reflect"
	"strings"

	serverTypes "github.com/duanchi/min/v2/server/types"
	"github.com/duanchi/min/v2/types"
	"github.com/duanchi/min/v2/util"
)

type RestfulBeanParser struct {
	types.BeanParser
}

func (parser RestfulBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {

	key := tag.Get("key")
	if key == "" {
		key = "id"
	}

	resourceList := util.ParseTag("@restful", tag, "path")

	if len(resourceList) == 0 {
		resourceList = util.ParseTag("restful", tag, "path")
	}

	for _, resource := range resourceList {
		if path, has := resource["path"]; has {
			pathKey := key
			if k, hasKey := resource["key"]; hasKey {
				pathKey = k
			}

			path = strings.ReplaceAll("/"+path, "//", "/")
			RestfulRoutes[path] = serverTypes.RestfulRoute{
				Value:       bean,
				ResourceKey: pathKey,
			}

			/*resources := strings.Split(resource, ",")
			for _, res := range resources {
				res = strings.TrimSpace(res)
				res = strings.ReplaceAll("/"+res, "//", "/")
				RestfulRoutes[res] = serverTypes.RestfulRoute{
					Value:       bean,
					ResourceKey: key,
				}
			}*/
		}
	}
}
