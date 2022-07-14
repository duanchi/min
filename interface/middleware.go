package _interface

import (
	"github.com/duanchi/min/types/middleware"
	"github.com/gin-gonic/gin"
)

type MiddlewareInterface interface {
	Includes() (includes middleware.Includes)
	Excludes() (excludes middleware.Excludes)
	BeforeRoute(ctx *gin.Context)
	AfterRoute(ctx *gin.Context)
	BeforeResponse(ctx *gin.Context)
	AfterResponse(ctx *gin.Context)
	AfterPanic(ctx *gin.Context)
}
