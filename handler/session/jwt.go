package session

import (
	"github.com/PasteUs/PasteMeGoBackend/config"
	"github.com/PasteUs/PasteMeGoBackend/logging"
	"github.com/PasteUs/PasteMeGoBackend/model/user"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var (
	IdentityKey    = "username"
	AuthMiddleware *jwt.GinJWTMiddleware
)

func init() {
	var err error
	if AuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "pasteme",
		Key:         []byte(config.Config.Secret),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*user.User); ok {
				return jwt.MapClaims{
					IdentityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			if username, ok := claims[IdentityKey]; ok {
				return username
			}
			return "nobody"
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var body user.User
			if err := c.ShouldBind(&body); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			// TODO login authenticator

			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"status":  code,
				"message": message,
			})
		},
		TokenLookup:   "cookie: token",
		TokenHeadName: "PasteMe",
		TimeFunc:      time.Now,
		DisabledAbort: true,
	}); err != nil {
		logging.Panic("jwt initializer create failed", zap.String("err", err.Error()))
	}

	if err = AuthMiddleware.MiddlewareInit(); err != nil {
		logging.Panic("jwt middleware init failed", zap.String("err", err.Error()))
	}
}
