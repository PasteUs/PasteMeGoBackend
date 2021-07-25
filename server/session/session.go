package session

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func Create(context *gin.Context) {
    context.JSON(http.StatusCreated, gin.H{
        "status": http.StatusCreated,
        "token":  "Hello World!",
    })
}

func Destroy(context *gin.Context) {
    context.JSON(http.StatusNoContent, gin.H{
        "status": http.StatusNoContent,
    })
}
