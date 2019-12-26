package utils

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"sync"
)

type ValidatorV10 struct {
	once     sync.Once
	validate *validator.Validate
}

func (v *ValidatorV10) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *ValidatorV10) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *ValidatorV10) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
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
