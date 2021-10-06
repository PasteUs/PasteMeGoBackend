package router

import (
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/common/flag"
	"github.com/PasteUs/PasteMeGoBackend/common/logging"
	"github.com/PasteUs/PasteMeGoBackend/handler/common"
	"github.com/PasteUs/PasteMeGoBackend/handler/paste"
	"github.com/PasteUs/PasteMeGoBackend/handler/token"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var router *gin.Engine

func init() {
	if !flag.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.Default()

	api := router.Group("/api")
	{
		v3 := api.Group("/v3")
		{
			v3.GET("/", common.Beat)

			s := v3.Group("/token")
			{
				s.POST("", token.AuthMiddleware.LoginHandler)    // 创建 Session（登陆）
				s.DELETE("", token.AuthMiddleware.LogoutHandler) // 销毁 Session（登出）
				s.GET("", token.AuthMiddleware.RefreshHandler)   // 刷新 Session
			}

			u := v3.Group("/user")
			{
				u.POST("")
				u.DELETE("")
				u.PUT("")
			}

			p := v3.Group("/paste")
			{
				p.POST("/", token.AuthMiddleware.MiddlewareFunc(true),
					paste.Create)         // 创建一个 Paste
				p.GET("/:key", paste.Get) // 读取 Paste
			}
		}
	}

	router.NoRoute(common.NotFoundHandler)
}

func Run(address string, port uint16) {
	if err := router.Run(fmt.Sprintf("%s:%d", address, port)); err != nil {
		logging.Panic("Run server failed", zap.Error(err))
	}
}
