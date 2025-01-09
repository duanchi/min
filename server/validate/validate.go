package validate

import (
	"errors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"reflect"
	"strings"
)

func Validate(obj interface{}) (err error) {
	err = engine.Struct(obj)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&validator.InvalidValidationError{}) {
			return err
		} else {
			errs := err.(validator.ValidationErrors)
			es := errs.Translate(trans)
			if len(es) > 0 {
				eMessages := []string{}
				for _, e := range es {
					eMessages = append(eMessages, e)
				}
				return errors.New(strings.Join(eMessages, ", "))
			}
		}
	}
	return nil
}

var trans ut.Translator
var engine = validator.New()
var Validators map[string]struct {
	validateFunction  validator.Func
	validateTranslate string
} = map[string]struct {
	validateFunction  validator.Func
	validateTranslate string
}{}

func Init() {
	zh := zh.New()
	en := en.New()
	uni := ut.New(en, zh)
	trans, _ = uni.GetTranslator("zh")

	zh_translations.RegisterDefaultTranslations(engine, trans)

	for tag, validate := range Validators {
		engine.RegisterValidation(tag, validate.validateFunction)
		engine.RegisterTranslation(tag, trans,
			func(ut ut.Translator) error {
				return ut.Add(tag, validate.validateTranslate, false)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field())
				if err != nil {
					log.Printf("警告: 字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		)
	}
}
