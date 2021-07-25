package session

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func Create(requests *gin.Context) {
    requests.JSON(http.StatusCreated, gin.H{
        "status": http.StatusCreated,
        "token":  "Hello World!",
    })
}

func Destroy(requests *gin.Context) {
    requests.JSON(http.StatusNoContent, gin.H{
        "status": http.StatusNoContent,
    })
}
