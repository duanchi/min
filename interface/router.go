package _interface

import "github.com/gin-gonic/gin"

type RouterInterface interface {
	Handle(path string, method string, params gin.Params, ctx *gin.Context)
}