package static

import (
	"github.com/duanchi/min/v2/context"
	"github.com/duanchi/min/v2/server/httpserver"
	"strings"
)

func Init(httpServer *httpserver.Httpserver) {

	if staticPath := context.GetApplicationContext().GetConfig("HttpServer.StaticPath").(string); staticPath != "" {
		staticPathStack := strings.Split(staticPath, ",")

		for _, path := range staticPathStack {
			if pathStack := strings.SplitN(path, ":", 2); len(pathStack) > 1 {
				httpServer.Static(pathStack[0], pathStack[1])
			}
		}
	}
}
