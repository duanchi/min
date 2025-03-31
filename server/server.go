package server

import (
	"github.com/duanchi/min/v2/context"
	"github.com/duanchi/min/v2/log"
	"github.com/duanchi/min/v2/server/httpserver"
	"github.com/duanchi/min/v2/server/middleware"
	"github.com/duanchi/min/v2/server/route"
	"github.com/duanchi/min/v2/server/static"
	"github.com/duanchi/min/v2/server/validate"
)

var HttpServer *httpserver.Httpserver

func Init(err chan error) {
	HttpServer = httpserver.New(struct{}{})

	if context.GetApplicationContext().GetConfig("Env").(string) == "production" {
		HttpServer.SetLogLevel(httpserver.LOG_ERROR)
	} else {
		HttpServer.SetLogLevel(httpserver.LOG_TRACE)
	}

	validate.Init()
	middleware.Init(HttpServer)
	middleware.InitBeforeRoute(HttpServer)
	static.Init(HttpServer)
	route.Init(HttpServer)

	serverError := HttpServer.Listen(
		context.GetApplicationContext().GetConfig("HttpServer.ServerHost").(string),
		context.GetApplicationContext().GetConfig("HttpServer.ServerPort").(string),
	)

	if serverError != nil {
		log.Log.Fatal(serverError)
	}

	err <- serverError
	return
}
