package controllers

import (
	"github.com/gin-gonic/gin"
)

//register Handlers
func Register(engine *gin.Engine) {
	engine.GET("/", HelloWorldController)
	engine.POST("/kv", KVSetController)
	engine.GET("/kv/:key", KVGetController)
}
