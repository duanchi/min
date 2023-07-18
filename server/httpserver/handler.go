package httpserver

import "github.com/duanchi/min/server/httpserver/context"

type Handler = func(*context.Context) error
