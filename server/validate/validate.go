package validate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	zh_translations "github.com/go-playground/validator/translations/zh"
	"gopkg.in/go-playground/validator.v9"
	"log"
)

var trans ut.Translator

var Validators map[string]struct {
	validateFunction validator.Func
	validateTranslate string
} = map[string]struct{
	validateFunction validator.Func
	validateTranslate string
}{}

func Init () {
	binding.Validator = new(defaultValidator)
	zh := zh.New()
	en := en.New()
	uni := ut.New(en, zh)
	trans, _ = uni.GetTranslator("zh")

	zh_translations.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), trans)

	for tag, validate := range Validators {
		binding.Validator.Engine().(*validator.Validate).RegisterValidation(tag, validate.validateFunction)
		binding.Validator.Engine().(*validator.Validate).RegisterTranslation(tag, trans,
			func(ut ut.Translator) error {
				return ut.Add(tag, validate.validateTranslate, false)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		)
	}

}
