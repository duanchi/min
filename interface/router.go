package _interface

import (
	"github.com/gofiber/fiber/v2"
)

type RouterInterface interface {
	Handle(path string, method string, params fiber., ctx *fiber.App)
}
