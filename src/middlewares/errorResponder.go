package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

func ErrorResponder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) == 0 {
			return
		}
		code := ctx.Writer.Status()
		for idx, err := range ctx.Errors {
			switch errs := err.Err.(type) {
			case validator.ValidationErrors:
				if trans, ok := ctx.Get("trans"); ok {
					errMessages := make([]string, len(errs))
					for i, err := range errs {
						errMessages[i] = err.Translate(trans.(ut.Translator))
					}
					ctx.Errors[idx] = &gin.Error{
						Err:  errors.New(strings.Join(errMessages, "; ")),
						Type: gin.ErrorTypeBind,
						Meta: nil,
					}
				}
				ctx.JSON(code, ctx.Errors.JSON())
			default:
				ctx.JSON(code, ctx.Errors.JSON())
			}
		}
	}
}
