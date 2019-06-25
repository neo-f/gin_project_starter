package controllers

import (
	"gin_project_starter/src/controllers/account"
	service "gin_project_starter/src/services/account"

	"github.com/gin-gonic/gin"
)

// register Handlers
func Register(engine *gin.Engine) {
	engine.Use(injectAccountService)
	accountG := engine.Group("account")
	accountG.GET("/", account.List)
	accountG.POST("/", account.Create)
	accountG.POST("/login", account.Login)
	{
		accountGDetail := accountG.Group(":id", account.DetailContext)
		accountGDetail.GET("/", account.Retrieve)
		accountGDetail.PUT("/", account.Update)
		accountGDetail.DELETE("/", account.Delete)
	}
}

func injectAccountService(ctx *gin.Context) {
	ctx.Set("service", service.NewPgService())
	ctx.Next()
}
