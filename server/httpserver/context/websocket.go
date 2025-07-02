package context

import (
	"github.com/gofiber/contrib/websocket"
)

type Websocket struct {
	connection *websocket.Conn
}

func NewWebsocket(connection *websocket.Conn) *Websocket {
	ctx := Websocket{connection: connection}
	return &ctx
}

func (this *Websocket) Locals(key string, value ...interface{}) interface{} {
	return this.connection.Locals(key, value...)
}

func (this *Websocket) Get(key string, defaults ...interface{}) (value ContextValue) {
	val := this.connection.Locals(key)
	if val == nil {
		if len(defaults) > 0 {
			return ContextValue{value: defaults[0]}
		} else {
			return ContextValue{}
		}
	}
	return val.(ContextValue)
}

func (this *Websocket) Set(key string, value interface{}) {
	this.connection.Locals(key, ContextValue{value: value})
}

func (this *Websocket) Has(key string) bool {
	return this.connection.Locals(key) != nil
}

func (this *Websocket) Params(key string, defaultValue ...string) string {
	return this.connection.Params(key, defaultValue...)
}
func (this *Websocket) Query(key string, defaultValue ...string) string {
	return this.connection.Query(key, defaultValue...)
}

func (this *Websocket) Cookies(key string, defaultValue ...string) string {
	return this.connection.Cookies(key, defaultValue...)
}
func (this *Websocket) Headers(key string, defaultValue ...string) string {
	return this.connection.Headers(key, defaultValue...)
}
func (this *Websocket) IP() string {
	return this.connection.IP()
}
func (this *Websocket) Close() {
	this.connection.Close()
}

func (this *Websocket) ReadMessage() (messageType int, p []byte, err error) {
	return this.connection.ReadMessage()
}
func (this *Websocket) WriteMessage(messageType int, p []byte) error {
	return this.connection.WriteMessage(messageType, p)
}
func (this *Websocket) Ctx() *websocket.Conn {
	return this.connection
}
