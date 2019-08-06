package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/go-playground/validator.v9/translations/en"
	"gopkg.in/go-playground/validator.v9/translations/zh"
)

func Translator() gin.HandlerFunc {
	universalTranslators := ut.New(en_US.New(), zh_Hans.New(), en_US.New())
	return func(ctx *gin.Context) {
		v := binding.Validator.Engine().(*validator.Validate)
		locale := ctx.DefaultQuery("locale", "en")
		trans, _ := universalTranslators.GetTranslator(locale)
		switch locale {
		case "zh_Hans":
			_ = zh.RegisterDefaultTranslations(v, trans)
		case "en_US":
			_ = en.RegisterDefaultTranslations(v, trans)
		default:
			_ = en.RegisterDefaultTranslations(v, trans)
		}
		ctx.Set("trans", trans)
	}
}
