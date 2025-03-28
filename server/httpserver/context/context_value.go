package context

import (
	"reflect"
)

func GetWithDefaults[T any](key string, ctx *Context, defaultValue T) T {
	value := ctx.Get(key)
	if value.IsNil() {
		return defaultValue
	} else {
		return value.Value().(T)
	}
}

type ContextValue struct {
	value interface{}
}

func (this ContextValue) Value() interface{} {
	if this.IsNil() {
		return interface{}(nil)
	}
	return this.value
}

func (this ContextValue) Int() int {
	if this.IsNil() {
		return 0
	}
	return this.value.(int)
}

func (this ContextValue) Int64() int64 {
	if this.IsNil() {
		return 0
	}
	return this.value.(int64)
}

func (this ContextValue) Int8() int8 {
	if this.IsNil() {
		return 0
	}
	return this.value.(int8)
}

func (this ContextValue) Int16() int16 {
	if this.IsNil() {
		return 0
	}
	return this.value.(int16)
}

func (this ContextValue) Float() float32 {
	if this.IsNil() {
		return 0
	}
	return this.value.(float32)
}

func (this ContextValue) Float64() float64 {
	if this.IsNil() {
		return 0
	}
	return this.value.(float64)
}

func (this ContextValue) Byte() byte {
	if this.IsNil() {
		return 0
	}
	return this.value.(byte)
}

func (this ContextValue) Bool() bool {
	if this.IsNil() {
		return false
	}
	return this.value.(bool)
}

func (this ContextValue) String() string {
	if this.IsNil() {
		return ""
	}
	return this.value.(string)
}

func (this ContextValue) IsNil() bool {
	return this.value == nil || reflect.ValueOf(this.value).IsZero()
}
