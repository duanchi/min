package abstract

import "reflect"

type BeanParser struct {
}

func (parser BeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type) {}
