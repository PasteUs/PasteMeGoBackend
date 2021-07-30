package router

import (
    "fmt"
    "github.com/PasteUs/PasteMeGoBackend/flag"
    "github.com/PasteUs/PasteMeGoBackend/handler/paste"
    "github.com/PasteUs/PasteMeGoBackend/handler/session"
    "github.com/PasteUs/PasteMeGoBackend/logging"
    v2Handler "github.com/PasteUs/PasteMeGoBackend/v2/handler"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

var router *gin.Engine

func init() {
    if !flag.GetArgv().Debug {
        gin.SetMode(gin.ReleaseMode)
    }
    router = gin.Default()
    api := router.Group("/api")
    {
        v2 := api.Group("/v2")
        {
            v2.GET("/", v2Handler.Beat) // 心跳检测
            // 访问未加密的 Paste，token 为 <Paste ID>
            // 访问加密的 Paste，token 为 <Paste ID>,<Password>
            v2.GET("/:token", v2Handler.Query)
            v2.POST("/", v2Handler.PermanentCreator)    // 创建一个永久的 Paste, key 是自增键
            v2.POST("/once", v2Handler.ReadOnceCreator) // 创建一个阅后即焚的 Paste, key 是随机的
            v2.PUT("/:key", v2Handler.TemporaryCreator) // 创建一个阅后即焚的 Paste, key 是指定的
        }

        v3 := api.Group("/v3")
        {
            v3.POST("/session", session.Create)    // 创建 Session（登陆）
            v3.DELETE("/session", session.Destroy) // 销毁 Session（登出）
            v3.POST("/:namespace", paste.Create)   // 创建一个 Paste
            v3.GET("/:namespace/:key", paste.Get)  // 读取 Paste
        }
    }

    router.NoRoute(v2Handler.NotFoundHandler)
}

func Run(address string, port uint16) {
    if err := router.Run(fmt.Sprintf("%s:%d", address, port)); err != nil {
        logging.Panic("Run server failed", zap.String("err", err.Error()))
    }
}
