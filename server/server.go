package server

import (
    "fmt"
    "github.com/PasteUs/PasteMeGoBackend/flag"
    "github.com/gin-gonic/gin"
    "github.com/wonderivan/logger"
)

var router *gin.Engine

func init() {
    if !flag.Debug {
        gin.SetMode(gin.ReleaseMode)
    }
    router = gin.Default()
    api := router.Group("/api")
    v1 := api.Group("/v1")
    {
        v1.GET("/", beat) // 心跳检测
        // 访问未加密的 Paste，token 为 <Paste ID>
        // 访问加密的 Paste，token 为 <Paste ID>,<Password>
        v1.GET("/:token", query)
        v1.POST("/", permanentCreator)    // 创建一个永久的 Paste, key 是自增键
        v1.POST("/once", readOnceCreator) // 创建一个阅后即焚的 Paste, key 是随机的
        v1.PUT("/:key", temporaryCreator) // 创建一个阅后即焚的 Paste, key 是指定的
    }
    router.NoRoute(notFoundHandler)
}

func Run(address string, port uint16, logToFile bool) {
    if err := router.Run(fmt.Sprintf("%s:%d", address, port)); err != nil {
        logger.Painc("Run server failed: " + err.Error())
    }
}
