package validate

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"strings"
	"sync"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyInit()

		if errs := v.validate.Struct(obj); errs != nil {

			objType := reflect.TypeOf(reflect.ValueOf(obj).Elem().Interface())
			replaceMap := map[string]struct{
				replace bool
				comment string
			}{}

			for n := 0; n < objType.NumField(); n++ {
				comment := objType.Field(n).Tag.Get("comment")
				if comment != "" {
					replace := false
					if len(comment) > 9 && comment[len(comment) - 9:len(comment) - 1] == ",replace" {
						replace = true
						comment = comment[0:len(comment) - 9]
 					}
					replaceMap[objType.Field(n).Name] = struct {
						replace bool
						comment string
					}{replace: replace, comment: comment}
				}
			}

			validateErrors := ValidationErrors{}

			for _, err := range errs.(validator.ValidationErrors) {
				transField := err.Translate(trans)

				for k, v := range replaceMap {
					if v.replace {
						transField = v.comment
					} else {
						transField = strings.ReplaceAll(transField, k, v.comment)
					}
				}
				validateError := fieldError{
					tag:         err.Tag(),
					actualTag:   err.ActualTag(),
					ns:          err.Namespace(),
					structNs:    err.StructNamespace(),
					field:       err.Field(),
					structfield: err.StructField(),
					value:       err.Value(),
					param:       err.Param(),
					kind:        err.Kind(),
					typ:         err.Type(),
					translate:   transField,
				}
				validateErrors = append(validateErrors, validateError)
			}

			return validateErrors
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		v.validate.SetTagName("validate")

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}