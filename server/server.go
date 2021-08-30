package server

import (
	"github.com/gin-gonic/gin"
	"github.com/duanchi/min/config"
	"github.com/duanchi/min/server/middleware"
	"github.com/duanchi/min/server/route"
	"github.com/duanchi/min/server/static"
	"github.com/duanchi/min/server/validate"
	"log"
)

var HttpServer *gin.Engine

func Init (err chan error) {
	HttpServer = gin.Default()


	if config.Get("Env").(string) == "production" {
		gin.SetMode("release")
	} else {
		gin.SetMode("debug")
	}

	validate.Init()

	middleware.Init(HttpServer, middleware.BeforeRoute)

	static.Init(HttpServer)

	route.Init(HttpServer)

	serverError := HttpServer.Run(":" + config.Get("Application.ServerPort").(string))

	if serverError != nil {
		log.Fatal(serverError)
	}

	err <- serverError

	return
}