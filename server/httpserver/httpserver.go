package httpserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const (
	LOG_TRACE = iota
	LOG_DEBUG
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_FATAL
	LOG_PANIC
)

type Httpserver struct {
	instance *fiber.App
}

func New(config interface{}) *Httpserver {
	return &Httpserver{
		instance: fiber.New(config.(fiber.Config)),
	}
}

func (this *Httpserver) Listen(host string, port string) error {
	return this.instance.Listen(host + ":" + port)
}

func (this *Httpserver) SetLogLevel(level int) {
	switch level {
	case LOG_TRACE:
		log.SetLevel(log.LevelTrace)
	case LOG_DEBUG:
		log.SetLevel(log.LevelDebug)
	case LOG_INFO:
		log.SetLevel(log.LevelInfo)
	case LOG_WARN:
		log.SetLevel(log.LevelWarn)
	case LOG_ERROR:
		log.SetLevel(log.LevelError)
	case LOG_FATAL:
		log.SetLevel(log.LevelFatal)
	case LOG_PANIC:
		log.SetLevel(log.LevelPanic)
	}
}
