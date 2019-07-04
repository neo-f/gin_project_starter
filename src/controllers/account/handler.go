package account

import (
	"gin_project_starter/src/services"
	"gin_project_starter/src/services/account"
	"gin_project_starter/src/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type Account struct {
	service services.AccountService
}

func Router(service services.AccountService) *Account {
	return &Account{service: service}
}

func (r Account) Delete(ctx *gin.Context) {
	acc := ctx.MustGet("account").(account.Account)
	if err := r.service.Delete(acc.ID); err != nil {
		_ = ctx.AbortWithError(400, err)
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (r Account) Update(ctx *gin.Context) {
	var req struct {
		Email    *string `json:"email" validate:"email"`
		Password *string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(400, err)
		return
	}
	acc := ctx.MustGet("user").(account.Account)
	fields := make([]string, 0)
	if req.Email != nil {
		acc.Email = *req.Email
		fields = append(fields, *req.Email)
	}
	if req.Password != nil {
		hash, _ := bcrypt.GenerateFromPassword([]byte(*req.Password), 10)
		acc.Email = string(hash)
		fields = append(fields, *req.Email)
	}

	if err := r.service.Update(&acc, fields...); err != nil {
		_ = ctx.AbortWithError(400, err)
		return
	}
	ctx.JSON(http.StatusOK, acc)
}

func (r Account) Retrieve(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ctx.MustGet("account"))
}

func (r Account) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" validate:"required, email"`
		Password string `json:"password" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(400, err)
		return
	}
	acc, token, err := r.service.Login(req.Email, req.Password)
	if err != nil {
		_ = ctx.AbortWithError(401, err)
		return
	}
	ctx.JSON(200, gin.H{
		"account": acc,
		"token":   token,
	})
}

func (r Account) Create(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(400, err)
		return
	}
	acc, err := r.service.Create(req.Username, req.Password, req.Email)
	if err != nil {
		_ = ctx.AbortWithError(400, err)
		return
	}
	ctx.JSON(201, acc)
}

func (r Account) List(ctx *gin.Context) {
	var req struct {
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}
	acc, total, err := r.service.List(req.Limit, req.Offset)
	if err != nil {
		_ = ctx.AbortWithError(400, err)
		return
	}
	ctx.JSON(200, utils.PaginatedResult(acc, total))
}

func (r Account) DetailContext(ctx *gin.Context) {
	var req struct {
		ID int64 `uri:"id" validate:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	acc, err := r.service.Retrieve(req.ID)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	ctx.Set("account", acc)
	ctx.Next()
}
