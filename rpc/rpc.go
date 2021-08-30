package rpc

import (
	"github.com/gin-gonic/gin"
	"github.com/duanchi/min/config"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/server/middleware"
	"github.com/duanchi/min/types"
	"net/http"
	"reflect"
	"strings"
)

type RpcBeanMap map[string]struct{
	Package string
	Instance reflect.Value
}

var RpcBeans = RpcBeanMap{}

func (this RpcBeanMap) Init (httpServer *gin.Engine) {
	prefix := config.Get("Rpc.Server.Prefix").(string)

	httpServer.POST(prefix + "/rpc/*caller", middleware.HandleAfterRoute, func(ctx *gin.Context) {

		defer func() {
			runtimeErr := recover()

			errResponse := struct {
				Message string
				Code int
			}{Code: 500}

			if runtimeErr != nil {
				if !reflect.TypeOf(runtimeErr).
					Implements(
						reflect.TypeOf(
							(*_interface.Error)(nil)).
							Elem()) {
					errResponse.Message = runtimeErr.(error).Error()
					ctx.JSON(http.StatusInternalServerError, errResponse)
				} else {
					errResponse.Message = runtimeErr.(types.RuntimeError).Error()
					errResponse.Code = runtimeErr.(types.RuntimeError).Code()
					ctx.JSON(runtimeErr.(_interface.Error).Code(), errResponse)
				}
			}
		}()

		pathStack := strings.SplitN(ctx.Param("caller")[len(prefix) + 1:], "::", 2)

		beanName := pathStack[0]
		methodName := pathStack[1]

		if bean, ok := RpcBeans[beanName]; ok {

			method := bean.Instance.MethodByName(methodName)

			if method.IsValid() {

				methodType := method.Type()
				parameters := []interface{}{}
				arguments := []reflect.Value{}
				response := []interface{}{}

				for i := 0; i < methodType.NumIn(); i++ {
					parameters = append(parameters, reflect.New(methodType.In(i)).Elem().Interface())
				}
				ctx.BindJSON(&parameters)

				if methodType.NumIn() != len(parameters) {
					panic(types.RuntimeError{
						Message:   "Malformed arguments in Method \"" + methodName + "\" in Class \"" + beanName + "\"",
						ErrorCode: http.StatusBadRequest,
					})
				}

				for i := 0; i < methodType.NumIn(); i++ {
					arguments = append(arguments, reflect.ValueOf(parameters[i]))
				}

				returns := method.Call(arguments)

				for i := 0; i < methodType.NumOut(); i++ {
					response = append(response, returns[i].Interface())
				}

				ctx.JSON(http.StatusOK, response)

			} else {
				panic(types.RuntimeError{
					Message:   "No implement of Method \"" + methodName + "\" in Class \"" + beanName + "\"",
					ErrorCode: http.StatusBadRequest,
				})
			}
		} else {
			panic(types.RuntimeError{
				Message:   "No implement of Class \"" + beanName + "\"",
				ErrorCode: http.StatusBadRequest,
			})
		}
	})
}