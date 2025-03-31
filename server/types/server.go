package types

import "github.com/duanchi/min/v2/server/httpserver/context"

type Route struct {
	Url     string
	Method  string
	Handler ServerHandleFunc
}

type ServerHandleFunc func(*context.Context)
