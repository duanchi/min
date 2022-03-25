package util

import "github.com/gin-gonic/gin"

func CtxGet(key string, ctx *gin.Context) interface{} {
	value, _ := ctx.Get(key)
	return value
}
