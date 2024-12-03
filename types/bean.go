package types

const (
	BEAN_TYPE_STRUCT = iota
	BEAN_TYPE_FIELD
	BEAN_TYPE_FUNCTION
	BEAN_TYPE_METHOD
)

type BeanDefinition struct {
	PackageName  string
	Type         int
	ReceiverName string
	Parameters   map[string]interface{}
	Value        interface{}
}

type BeanMapper map[string]BeanDefinition
