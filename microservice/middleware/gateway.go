package middleware

import (
	"encoding/base64"
	"encoding/json"
	"github.com/duanchi/min/v2/abstract"
	"github.com/duanchi/min/v2/server/httpserver/context"
	"github.com/duanchi/min/v2/types/gateway"
)

type GatewayMiddleware struct {
	abstract.Middleware
}

func (this *GatewayMiddleware) AfterRoute(ctx *context.Context) {

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
