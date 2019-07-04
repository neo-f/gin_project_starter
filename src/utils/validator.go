package utils

import (
	"reflect"
	"sync"

	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	zh_trans "gopkg.in/go-playground/validator.v9/translations/zh"
)

type ValidatorV9 struct {
	once     sync.Once
	validate *validator.Validate
}

func (v *ValidatorV9) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *ValidatorV9) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *ValidatorV9) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("validate")
		// add any custom validations etc. here
		zh := zhongwen.New()
		uni := ut.New(zh, zh)
		trans, _ := uni.GetTranslator("zh")
		_ = zh_trans.RegisterDefaultTranslations(v.validate, trans)
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
