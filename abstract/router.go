package abstract

import "github.com/gin-gonic/gin"

type Router struct {
	Bean
}

func (this *Router) Handle (path string, method string, params gin.Params, ctx *gin.Context) {}
