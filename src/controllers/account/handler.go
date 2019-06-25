package account

import (
	"gin_project_starter/src/services/account"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func Delete(ctx *gin.Context) {
	service := ctx.MustGet("service").(account.Service)
	user := ctx.MustGet("user").(account.Account)
	if err := service.Delete(user.ID); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func Update(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"email"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	service := ctx.MustGet("service").(account.Service)
	a := &account.Account{
		Email: req.Email,
	}
	if err := service.Update(a, "email"); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, a)
}

func Retrieve(ctx *gin.Context) {

}

func Login(ctx *gin.Context) {

}

func Create(ctx *gin.Context) {

}

func List(ctx *gin.Context) {

}

func DetailContext(ctx *gin.Context) {
	id, err := cast.ToInt64E(ctx.Param("id"))
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	service := ctx.MustGet("service").(account.Service)
	user, err := service.Retrieve(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	ctx.Set("user", user)
	ctx.Next()
}
