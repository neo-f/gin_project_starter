package controllers

import (
	"gin_project_starter/src/controllers/account"
	"gin_project_starter/src/services"

	"github.com/gin-gonic/gin"
)

// register Handlers
func Register(engine *gin.Engine) {
	accountG := engine.Group("account")
	accountR := account.Router(services.NewAccountService())
	accountG.GET("/", accountR.List)
	accountG.POST("/", accountR.Create)
	accountG.POST("/login", accountR.Login)
	{
		accountGDetail := accountG.Group(":id", accountR.DetailContext)
		accountGDetail.GET("/", accountR.Retrieve)
		accountGDetail.PUT("/", accountR.Update)
		accountGDetail.DELETE("/", accountR.Delete)
	}
}
