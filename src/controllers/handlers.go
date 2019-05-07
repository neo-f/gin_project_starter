package controllers

import (
	"gin_project_starter/src/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloWorldController(ctx *gin.Context) {
	var query struct {
		Format string `form:"format"`
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if query.Format == "json" {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello world!",
		})
	} else {
		ctx.String(http.StatusOK, "hello world!")
	}
}

func KVSetController(ctx *gin.Context) {
	var reqJson struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&reqJson); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	services.SetKeyValue(reqJson.Key, reqJson.Value)
	ctx.JSON(http.StatusCreated, gin.H{
		"key":   reqJson.Key,
		"value": reqJson.Value,
	})
}

func KVGetController(ctx *gin.Context) {
	key := ctx.Param("key")
	value := services.GetValue(key)
	ctx.JSON(http.StatusOK, gin.H{
		"value": value,
	})
}
