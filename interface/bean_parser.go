package _interface

import "reflect"

type BeanParserInterface interface {
	Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string)
}