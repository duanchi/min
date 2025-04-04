package handler

import (
	"fmt"
	"github.com/duanchi/min/v2/config"
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/server/httpserver"
	"github.com/duanchi/min/v2/server/httpserver/context"
	"github.com/duanchi/min/v2/server/middleware"
	serverTypes "github.com/duanchi/min/v2/server/types"
	"github.com/duanchi/min/v2/server/websocket"
	"github.com/duanchi/min/v2/types"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"
)

func RestfulHandle(resource string, controller serverTypes.RestfulRoute, ctx *context.Context, engine *httpserver.Httpserver) error {
	isCustomResponse := config.Get("HttpServer.Restful.CustomResponse").(bool)
	params := ctx.Params()
	id := ctx.Param(controller.ResourceKey)
	method := ctx.Request().Method()
	requestId := ctx.Request().Header("Request-Id")
	beforeResponseHandlers := middleware.GetHandlersBeforeResponse()
	if requestId == "" {
		requestId = uuid.NewV4().String()
	}
	response := types.Response{
		RequestId: requestId,
		Status:    false,
		Code:      -1,
		Data:      nil,
		Message:   "Ok",
	}

	defer func() {
		statusCode := http.StatusInternalServerError

		if exception := recover(); exception != nil {
			defer func() {
				ctx.JSONWithStatus(statusCode, response)
			}()

			_, implemented := exception.(types.Error)

			if implemented {
				runtimeError := reflect.ValueOf(exception).Interface().(types.Error)

				statusCode = runtimeError.Status()
				response.Message = runtimeError.Error()
				response.Data = runtimeError.Data()
				response.Code = runtimeError.Code()

			} else {
				commonError := reflect.ValueOf(exception).Interface().(error)
				response.Message = commonError.Error()
			}
			fmt.Println(exception.(error).Error())
			fmt.Println(string(debug.Stack()))
		}
	}()

	var data interface{}
	var err error
	executor := controller.Value.Interface().(_interface.RestfulControllerInterface)
	// Upgrade Protocol to Websocket
	if method == "GET" {
		upgradeRequest := ctx.Request().Header("Connection")
		upgradeProtocol := ctx.Request().Header("Upgrade")

		if upgradeRequest == "Upgrade" && strings.ToLower(upgradeProtocol) == "websocket" {
			return websocket.Handle(id, resource, params, ctx, executor.Connect)
		}
	}

	switch method {
	case "GET":
		if id == "" {
			data, err = executor.FetchList(id, resource, params, ctx)
		} else {
			data, err = executor.Fetch(id, resource, params, ctx)
		}
	case "POST":
		data, err = executor.Create(id, resource, params, ctx)
	case "PUT":
		data, err = executor.Update(id, resource, params, ctx)
	case "DELETE":
		data, err = executor.Remove(id, resource, params, ctx)
	case "HEAD":
		data, err = executor.Fetch(id, resource, params, ctx)
	case "OPTIONS":
		data, err = executor.Fetch(id, resource, params, ctx)
	}

	if err == nil {
		status := http.StatusOK
		switch method {
		case "GET":
			status = 200
		case "POST":
			status = 201
		case "PUT":
			status = 201
		case "DELETE":
			status = 204
		}
		for _, handler := range beforeResponseHandlers {
			handler(ctx)
		}
		if isCustomResponse {
			return ctx.JSONWithStatus(status, data)
		} else {
			response.Status = true
			response.Data = data
			response.Code = 0
			return ctx.JSONWithStatus(status, response)
		}
	} else {
		for _, handler := range beforeResponseHandlers {
			handler(ctx)
		}
		status := http.StatusInternalServerError
		response.Status = false
		if _, implemented := err.(types.Error); implemented {
			runtimeError := reflect.ValueOf(err).Interface().(types.Error)

			response.Data = runtimeError.Data()
			response.Code = runtimeError.Code()
			response.Message = runtimeError.Error()
			status = runtimeError.Status()
		} else {
			response.Data = nil
			response.Code = -1
			response.Message = err.Error()
		}

		return ctx.JSONWithStatus(status, response)
		// panic(err)
	}
}
