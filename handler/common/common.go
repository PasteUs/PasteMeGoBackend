package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Status int `json:"status" example:"200"`
}

func NotFoundHandler(context *gin.Context) {
	Error(context, http.StatusNotFound, ErrNoRouterFounded)
}

// Beat godoc
// @Summary 心跳检测
// @Description 心跳检测
// @Tags Common
// @Produce json
// @Param method query string true "方法" Enums("beat")
// @Success 200 {object} common.Response
// @Failure default {object} common.ErrorResponse
// @Router / [get]
func Beat(context *gin.Context) {
	method := context.DefaultQuery("method", "none")
	if method == "beat" {
		context.JSON(http.StatusOK, Response{http.StatusOK})
	} else {
		NotFoundHandler(context)
	}
}
