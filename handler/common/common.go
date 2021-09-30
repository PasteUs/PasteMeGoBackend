package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IResponse interface {
	GetHttpStatusCode() int
}

type Response struct {
	Code int `json:"code" example:"200"`
}

func (response *Response) GetHttpStatusCode() int {
	if response.Code > 1000 {
		return response.Code / 100
	}
	return response.Code
}

func JSON(context *gin.Context, response IResponse) {
	context.JSON(response.GetHttpStatusCode(), response)
}

func NotFoundHandler(context *gin.Context) {
	ErrNoRouterFounded.Abort(context)
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
		JSON(context, &Response{Code: http.StatusOK})
	} else {
		NotFoundHandler(context)
	}
}
