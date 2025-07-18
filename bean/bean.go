package bean

import (
	"fmt"
	"reflect"

	"github.com/duanchi/min/v2/bean/core_parsers"
	_interface "github.com/duanchi/min/v2/interface"
)

// var beanContainer interface{}
var beanMaps = map[string]reflect.Value{}
var beanNameMaps = map[string]reflect.Value{}
var beanTypeMaps = map[reflect.Type]reflect.Value{}

var customBeanParsers = []_interface.BeanParserInterface{}

type Container struct{}

func (bean *Container) Get(name string) reflect.Value {

	beanValue := reflect.ValueOf(bean).Elem()
	beanType := reflect.TypeOf(bean).Elem()

	value := reflect.ValueOf(bean).Elem().FieldByName(name)

	if reflect.Value.IsZero(value) {
		for i := 0; i < beanType.NumField(); i++ {
			if name == beanType.Field(i).Tag.Get("name") {
				return beanValue.Field(i)
			}
		}
	}

	return value
}

func InitBeans(beanContainerInstance interface{}, beanParsers interface{}) {

	if reflect.ValueOf(beanParsers).IsValid() {
		customBeanParsers = beanParsers.([]_interface.BeanParserInterface)
	}

	containerValue := reflect.ValueOf(beanContainerInstance)
	containerType := reflect.TypeOf(beanContainerInstance)

	if reflect.TypeOf(beanContainerInstance).Kind() == reflect.Ptr {
		containerValue = reflect.ValueOf(beanContainerInstance).Elem()
		containerType = reflect.TypeOf(beanContainerInstance).Elem()
	}

	// initBean(containerValue, containerType)

	// 保持先注册、再初始化、最后注入的步骤，所以执行三次完整循环

	for i := 0; i < containerValue.NumField(); i++ {
		Register(containerValue.Field(i), containerType.Field(i))
	}

	for name, bean := range beanMaps {
		Init(bean, name, beanMaps)
	}

	for name, bean := range beanMaps {
		Inject(bean, name, beanMaps)
	}
}

func Get(name string) interface{} {
	return beanNameMaps[name].Interface()
}

func GetAll() map[string]reflect.Value {
	return beanMaps
}

func Register(beanValue reflect.Value, beanDefinition reflect.StructField) {
	tag := beanDefinition.Tag
	// beanType := beanDefinition.Type
	name := tag.Get("name")
	if name == "" {
		name = beanDefinition.Name
	}
	beanMaps[name] = reflect.New(beanDefinition.Type).Elem()
	beanMaps[name].Addr().Interface().(_interface.BeanInterface).SetName(name)
	beanNameMaps[name] = beanMaps[name].Addr()
	beanTypeMaps[beanMaps[name].Addr().Type()] = beanMaps[name].Addr()

	parseBean(tag, beanMaps[name].Addr(), beanDefinition.Type, name)

	fmt.Println("[min-framework] Register Bean: " + name + " Ok!")
}

func parseBean(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	for i := 0; i < len(core_parsers.CoreBeanParsers); i++ {
		reflect.ValueOf(core_parsers.CoreBeanParsers[i]).Interface().(_interface.BeanParserInterface).Parse(tag, bean, definition, beanName)
	}

	for i := 0; i < len(customBeanParsers); i++ {
		reflect.ValueOf(customBeanParsers[i]).Interface().(_interface.BeanParserInterface).Parse(tag, bean, definition, beanName)
	}
}
