package validate

import (
	"bytes"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"reflect"
	"strings"
)

const (
	fieldErrMsg = "Key: '%s' Error:Field validation for '%s' failed on the '%s' tag"
)

type ValidationErrors [] fieldError


func (ve ValidationErrors) Error() string {

	buff := bytes.NewBufferString("")

	for i := 0; i < len(ve); i++ {

		fe := ve[i]
		buff.WriteString(fe.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

// Translate translates all of the ValidationErrors
func (ve ValidationErrors) Translate(ut ut.Translator) validator.ValidationErrorsTranslations {

	trans := make(validator.ValidationErrorsTranslations)

	for i := 0; i < len(ve); i++ {
		fe := ve[i]

		// // in case an Anonymous struct was used, ensure that the key
		// // would be 'Username' instead of ".Username"
		// if len(fe.ns) > 0 && fe.ns[:1] == "." {
		// 	trans[fe.ns[1:]] = fe.Translate(ut)
		// 	continue
		// }

		trans[fe.ns] = fe.Translate(ut)
	}

	return trans
}

func (ve ValidationErrors) Has(field string) bool {
	for _, v := range ve {
		if v.structfield == field {
			return true
		}
	}

	return false
}

/*func (ve ValidationErrors) Translation() string {
	transStrings := []string{}
	for _, err := range ve.ValidationErrors {

		transString := err.Translate(trans)

		for k, v := range ve.translateMap {
			transString = strings.ReplaceAll(transString, k, v)
		}

		transStrings = append(transStrings, transString)
	}

	return strings.Join(transStrings, "\n")
}*/

/*func (ve ValidationErrors) FieldTranslation() map[string]string {
	transFields := map[string]string{}
	for _, err := range ve.ValidationErrors {
		transField := err.Translate(trans)

		for k, v := range ve.translateMap {
			transField = strings.ReplaceAll(transField, k, v)
		}

		transFields[err.StructField()] = transField
	}

	return transFields
}*/


type fieldError struct {
	tag            string
	actualTag      string
	ns             string
	structNs       string
	field          string
	structfield    string
	value          interface{}
	param          string
	kind           reflect.Kind
	typ            reflect.Type
	translate	   string
}

// Tag returns the validation tag that failed.
func (fe *fieldError) Tag() string {
	return fe.tag
}

// ActualTag returns the validation tag that failed, even if an
// alias the actual tag within the alias will be returned.
func (fe *fieldError) ActualTag() string {
	return fe.actualTag
}

// Namespace returns the namespace for the field error, with the tag
// name taking precedence over the fields actual name.
func (fe *fieldError) Namespace() string {
	return fe.ns
}

// StructNamespace returns the namespace for the field error, with the fields
// actual name.
func (fe *fieldError) StructNamespace() string {
	return fe.structNs
}

// Field returns the fields name with the tag name taking precedence over the
// fields actual name.
func (fe *fieldError) Field() string {

	return fe.field
	// // return fe.field
	// fld := fe.ns[len(fe.ns)-int(fe.fieldLen):]

	// log.Println("FLD:", fld)

	// if len(fld) > 0 && fld[:1] == "." {
	// 	return fld[1:]
	// }

	// return fld
}

// returns the fields actual name from the struct, when able to determine.
func (fe *fieldError) StructField() string {
	// return fe.structField
	return fe.structfield
}

// Value returns the actual fields value in case needed for creating the error
// message
func (fe *fieldError) Value() interface{} {
	return fe.value
}

// Param returns the param value, in string form for comparison; this will
// also help with generating an error message
func (fe *fieldError) Param() string {
	return fe.param
}

// Kind returns the Field's reflect Kind
func (fe *fieldError) Kind() reflect.Kind {
	return fe.kind
}

// Type returns the Field's reflect Type
func (fe *fieldError) Type() reflect.Type {
	return fe.typ
}

// Error returns the fieldError's error message
func (fe *fieldError) Error() string {
	if fe.translate != "" {
		return fe.translate
	}
	return fmt.Sprintf(fieldErrMsg, fe.ns, fe.Field(), fe.tag)
}

// Translate returns the FieldError's translated error
// from the provided 'ut.Translator' and registered 'TranslationFunc'
//
// NOTE: is not registered translation can be found it returns the same
// as calling fe.Error()
func (fe *fieldError) Translate(ut ut.Translator) string {

	if fe.translate == "" {
		return fe.Error()
	}
	return fe.translate
}
