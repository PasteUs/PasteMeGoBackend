package paste

import (
    _ "github.com/PasteUs/PasteMeGoBackend/model/paste"
    "github.com/PasteUs/PasteMeGoBackend/util/logging"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "net/http"
)

func Create(requests *gin.Context) {
    requests.JSON(http.StatusCreated, gin.H{
        "status": http.StatusCreated,
        "token":  "Hello World!",
    })
}

func Get(requests *gin.Context) {
    namespace, key := requests.Param("namespace"), requests.Param("key")
    logging.Info("test", requests, zap.String("namespace", namespace), zap.String("key", key))
    requests.JSON(http.StatusOK, gin.H{
        "status": http.StatusOK,
        "token":  "Hello World!",
    })
}
