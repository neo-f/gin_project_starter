package account

import (
	"errors"
	"gin_project_starter/src/services"
	"gin_project_starter/src/utils"
	"net/http"

	"github.com/go-pg/pg/v9"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type Account struct {
	service services.AccountService
	key     string
}

func Router(service services.AccountService) *Account {
	return &Account{
		service: service,
		key:     "account",
	}
}

func (r Account) Delete(ctx *gin.Context) {
	acc := ctx.MustGet(r.key).(services.Account)
	if err := r.service.Delete(acc.ID); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (r Account) Update(ctx *gin.Context) {
	var req struct {
		Email    *string `json:"email" validate:"email"`
		Password *string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	acc := ctx.MustGet(r.key).(services.Account)
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
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, acc)
}

func (r Account) Retrieve(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ctx.MustGet(r.key))
}

func (r Account) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" validate:"required, email"`
		Password string `json:"password" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	acc, token, err := r.service.Login(req.Email, req.Password)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"postgres": acc,
		"token":    token,
	})
}

func (r Account) Create(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	acc, err := r.service.Create(req.Username, req.Password, req.Email)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, acc)
}

func (r Account) List(ctx *gin.Context) {
	var req struct {
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	acc, total, err := r.service.List(req.Limit, req.Offset)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, utils.PaginatedResult(acc, total))
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
		if errors.Is(err, pg.ErrNoRows) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.Set(r.key, acc)
	ctx.Next()
}
