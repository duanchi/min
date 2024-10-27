package server

import (
	"github.com/duanchi/min/context"
	"github.com/duanchi/min/log"
	"github.com/duanchi/min/server/httpserver"
	"github.com/duanchi/min/server/middleware"
	"github.com/duanchi/min/server/route"
	"github.com/duanchi/min/server/static"
	"github.com/duanchi/min/server/validate"
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
