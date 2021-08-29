package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotFoundHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"status":  http.StatusNotFound,
		"message": ErrNoRouterFounded.Error(),
	})
}

func Beat(context *gin.Context) {
	method := context.DefaultQuery("method", "none")
	if method == "beat" {
		context.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
		})
	} else {
		NotFoundHandler(context)
	}
}
