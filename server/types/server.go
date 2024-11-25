package types

import "github.com/duanchi/min/server/httpserver/context"

type Route struct {
	Url     string
	Method  string
	Handler ServerHandleFunc
}

type ServerHandleFunc func(*context.Context)
