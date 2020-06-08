package middlewares

import (
	"github.com/gin-gonic/gin"
)

func ErrorResponder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) == 0 {
			return
		}
		code := ctx.Writer.Status()
		ctx.JSON(code, ctx.Errors.JSON())
	}
}
