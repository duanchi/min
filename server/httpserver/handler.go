package httpserver

import "github.com/duanchi/min/v2/server/httpserver/context"

type Handler = func(*context.Context) error
