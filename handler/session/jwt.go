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

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var (
	identityKey    = "username"
	AuthMiddleware *jwt.GinJWTMiddleware
)

func init() {
	var err error
	if AuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "pasteme",
		Key:         []byte(config.Config.Secret),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*user.User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &user.User{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var body login
			if err := c.ShouldBind(&body); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := body.Username
			password := body.Password

			if (username == "admin" && password == "admin") || (username == "test" && password == "test") {
				c.Set("username", username)
				return &user.User{
					Username: username,
				}, nil
			}

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
	}); err != nil {
		logging.Panic("jwt initializer create failed", zap.String("err", err.Error()))
	}

	if err = AuthMiddleware.MiddlewareInit(); err != nil {
		logging.Panic("jwt middleware init failed", zap.String("err", err.Error()))
	}
}
