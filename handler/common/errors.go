package common

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ErrNoRouterFounded = errors.New("no router founded")
)

type ErrorResponse struct {
	*Response
	Message string `json:"message" example:"ok"`
}

func Error(context *gin.Context, status int, err error) {
	httpStatusCode := 0

	if status < 1000 {
		httpStatusCode = status
	} else {
		httpStatusCode = status / 100
	}

	response := ErrorResponse{&Response{Status: status}, err.Error()}
	context.JSON(httpStatusCode, response)
}
