package utils

import "github.com/gin-gonic/gin"

func PaginatedResult(objs interface{}, total int) gin.H {
	return gin.H{
		"results": objs,
		"total":   total,
	}
}
