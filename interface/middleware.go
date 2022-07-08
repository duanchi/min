package _interface

import "github.com/gin-gonic/gin"

type MiddlewareInterface interface {
	Includes() (includes map[string][]string)
	Excludes() (excludes map[string][]string)
	BeforeRoute(ctx *gin.Context)
	AfterRoute(ctx *gin.Context)
	BeforeResponse(ctx *gin.Context)
	AfterResponse(ctx *gin.Context)
	AfterPanic(ctx *gin.Context)
}
