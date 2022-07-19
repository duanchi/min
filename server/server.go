package server

import (
	"github.com/duanchi/min/context"
	"github.com/duanchi/min/log"
	"github.com/duanchi/min/server/middleware"
	"github.com/duanchi/min/server/route"
	"github.com/duanchi/min/server/static"
	"github.com/duanchi/min/server/validate"
	"github.com/gin-gonic/gin"
)

var HttpServer *gin.Engine

func Init(err chan error) {
	HttpServer = gin.Default()

	if context.GetApplicationContext().GetConfig("Env").(string) == "production" {
		gin.SetMode("release")
	} else {
		gin.SetMode("debug")
	}

	validate.Init()
	middleware.Init(HttpServer, middleware.BeforeRoute)
	static.Init(HttpServer)
	route.Init(HttpServer)

	serverError := HttpServer.Run(":" + context.GetApplicationContext().GetConfig("HttpServer.ServerPort").(string))

	if serverError != nil {
		log.Log.Fatal(serverError)
	}

	err <- serverError
	return
}
