package middleware

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types/gateway"
)

type GatewayMiddleware struct {
	abstract.Middleware
}

func (this *GatewayMiddleware) AfterRoute (ctx *gin.Context) {

	data := gateway.Data{}


	gatewayData := ctx.Request.Header.Get("X-Gateway-Data")

	if gatewayData == "" && ctx.Query("__X-GATEWAY-DATA__") != "" {
		gatewayData = ctx.Query("__X-GATEWAY-DATA__")
	}

	if decodeData, ok := base64.URLEncoding.DecodeString(gatewayData); ok == nil {
		json.Unmarshal(decodeData, &data)
	}

	ctx.Set("GATEWAY_DATA", data)
}
