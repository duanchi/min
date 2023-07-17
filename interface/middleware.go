package _interface

import (
	"github.com/duanchi/min/types/middleware"
	"github.com/gofiber/fiber/v2"
)

type MiddlewareInterface interface {
	Includes() (includes middleware.Includes)
	Excludes() (excludes middleware.Excludes)
	BeforeRoute(ctx *fiber.Ctx)
	AfterRoute(ctx *fiber.Ctx)
	BeforeResponse(ctx *fiber.Ctx)
	AfterResponse(ctx *fiber.Ctx)
	AfterPanic(ctx *fiber.Ctx)
}
