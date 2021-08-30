package service

import "reflect"

type ServiceBeanMap map[string]reflect.Value

var ServiceBeans = ServiceBeanMap{}
