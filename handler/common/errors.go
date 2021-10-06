package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrZeroExpireSecond               = New(http.StatusBadRequest, 1, "zero expire time")
	ErrZeroExpireCount                = New(http.StatusBadRequest, 2, "zero expire count")
	ErrExpireSecondGreaterThanMonth   = New(http.StatusBadRequest, 3, "expire minute greater than a month")
	ErrExpireCountGreaterThanMaxCount = New(http.StatusBadRequest, 4, "expire count greater than max count")
	ErrEmptyContent                   = New(http.StatusBadRequest, 5, "empty content")
	ErrEmptyLang                      = New(http.StatusBadRequest, 6, "empty lang")
	ErrInvalidLang                    = New(http.StatusBadRequest, 7, "invalid lang")
	ErrWrongParamType                 = New(http.StatusBadRequest, 8, "wrong param type")
	ErrInvalidKeyLength               = New(http.StatusBadRequest, 9, "invalid key length")
	ErrInvalidKeyFormat               = New(http.StatusBadRequest, 10, "invalid key format")

	ErrUnauthorized = New(http.StatusUnauthorized, 1, "unauthorized")

	ErrWrongPassword = New(http.StatusForbidden, 1, "wrong password")

	ErrNoRouterFounded = New(http.StatusNotFound, 1, "no router founded")
	ErrRecordNotFound  = New(http.StatusNotFound, 2, "record not found")

	ErrQueryDBFailed = New(http.StatusInternalServerError, 1, "query from db failed")
	ErrSaveFailed    = New(http.StatusInternalServerError, 2, "save failed")
)

type ErrorResponse struct {
	*Response
	Message string `json:"message" example:"ok"`
}

func (response *ErrorResponse) Error() string {
	return response.Message
}

func (response *ErrorResponse) Abort(context *gin.Context) {
	context.AbortWithStatusJSON(response.GetHttpStatusCode(), response)
}

func New(code int, index int, message string) *ErrorResponse {
	return &ErrorResponse{
		Response: &Response{code*100 + index},
		Message:  message,
	}
}
