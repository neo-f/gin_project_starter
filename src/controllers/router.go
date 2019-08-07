package controllers

import (
	"gin_project_starter/src/controllers/account"
	"gin_project_starter/src/services"

	"github.com/gin-gonic/gin"
)

// register Handlers
func Register(engine *gin.Engine) {
	accountRouter := account.Router(services.NewAccountService())
	{
		accountG := engine.Group("account")
		accountG.GET("/", accountRouter.List)
		accountG.POST("/", accountRouter.Create)
		accountG.POST("/login", accountRouter.Login)
		{
			accountGDetail := accountG.Group("/:id", accountRouter.DetailContext)
			accountGDetail.GET("/", accountRouter.Retrieve)
			accountGDetail.PUT("/", accountRouter.Update)
			accountGDetail.DELETE("/", accountRouter.Delete)
		}
	}
}
