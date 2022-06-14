package handler

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/server/middleware"
	"github.com/duanchi/min/server/websocket"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"
)

func RestfulHandle(resource string, controller reflect.Value, ctx *gin.Context, engine *gin.Engine) {
	params := ctx.Params
	id := ctx.Param("id")
	method := ctx.Request.Method
	requestId := ctx.Request.Header.Get("Request-Id")
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
				/*if ctx.Writer.Status() != http.StatusOK {
					statusCode = ctx.Writer.Status()
				}*/
				ctx.JSON(statusCode, response)
				debug.PrintStack()
			}()

			_, implemented := exception.(types.Error)

			if implemented {
				runtimeError := reflect.ValueOf(exception).Interface().(types.Error)
				/*switch exception.(type) {

				case error2.RequestError:
					statusCode = http.StatusBadRequest

				case error2.ResponseError:
					statusCode = http.StatusInternalServerError

				case error2.AuthorizeError:
					statusCode = http.StatusUnauthorized

				case error2.ForbiddenError:
					statusCode = http.StatusForbidden

				case error2.NotFoundError:
					statusCode = http.StatusNotFound

				default:
					if runtimeError.Code() < 600 {
						statusCode = runtimeError.Code()
					} else {
						statusCode = 500
					}
				}*/

				statusCode = runtimeError.Status()
				response.Message = runtimeError.Error()
				response.Data = runtimeError.Data()
				response.Code = runtimeError.Code()

			} else {
				commonError := reflect.ValueOf(exception).Interface().(error)
				response.Message = commonError.Error()
			}
		}
	}()

	var data interface{}
	var err error

	executor := controller.Interface().(_interface.RestControllerInterface)

	// Upgrade Protocol to Websocket
	if method == "GET" {
		upgradeRequest := ctx.Request.Header.Get("Connection")
		upgradeProtocol := ctx.Request.Header.Get("Upgrade")

		if upgradeRequest == "Upgrade" && strings.ToLower(upgradeProtocol) == "websocket" {
			websocket.Handle(id, resource, &params, ctx, executor.Connect)
			return
		}
	}

	switch method {
	case "GET":
		data, err = executor.Fetch(id, resource, &params, ctx)
	case "POST":
		data, err = executor.Create(id, resource, &params, ctx)
	case "PUT":
		data, err = executor.Update(id, resource, &params, ctx)
	case "DELETE":
		data, err = executor.Remove(id, resource, &params, ctx)
	case "HEAD":
		data, err = executor.Fetch(id, resource, &params, ctx)
	case "OPTIONS":
		data, err = executor.Fetch(id, resource, &params, ctx)
	}

	if err == nil {
		response.Status = true
		response.Data = data
		response.Code = 0
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
		ctx.JSON(status, response)
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

		ctx.JSON(status, response)
		// panic(err)
	}

	return
}
