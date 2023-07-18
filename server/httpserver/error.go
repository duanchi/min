package httpserver

import (
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/gofiber/fiber/v2"
)

type ErrorHandler = func(context.Context, error) error

// Error represents an error that occurred while handling a request.
type Error struct {
	instance *fiber.Error
	Code     int    `json:"code"`
	Message  string `json:"message"`
}

// Error makes it compatible with the `error` interface.
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new Error instance with an optional message
func NewError(code int, message ...string) *Error {
	err := fiber.NewError(code, message...)
	return &Error{
		instance: err,
		Code:     err.Code,
		Message:  err.Message,
	}
}
