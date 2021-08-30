package static

import (
	"github.com/gin-gonic/gin"
	"github.com/duanchi/min/config"
	"strings"
)

func Init (httpServer *gin.Engine) {


	if staticPath := config.Get("Application.StaticPath").(string); staticPath != "" {
		staticPathStack := strings.Split(staticPath, ",")

		for _, path := range staticPathStack {
			if pathStack := strings.SplitN(path, ":", 2); len(pathStack) > 1 {
				httpServer.Static(pathStack[0], pathStack[1])
			}
		}
	}
}
