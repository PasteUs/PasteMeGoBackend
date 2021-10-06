package token

import (
	"github.com/PasteUs/PasteMeGoBackend/common/config"
	"github.com/PasteUs/PasteMeGoBackend/common/logging"
	"github.com/PasteUs/PasteMeGoBackend/model/user"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var (
	IdentityKey    = "username"
	Nobody         = "nobody"
	AuthMiddleware *JWTMiddleware
)

func authenticator(c *gin.Context) (interface{}, error) {
	var body user.User
	if err := c.ShouldBind(&body); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	// TODO login authenticator

	return nil, jwt.ErrFailedAuthentication
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if username, ok := data.(string); ok {
		return jwt.MapClaims{
			IdentityKey: username,
		}
	}
	return jwt.MapClaims{}
}

func init() {
	var err error
	AuthMiddleware = &JWTMiddleware{
		&jwt.GinJWTMiddleware{
			Realm:         "pasteme",
			Key:           []byte(config.Config.Secret),
			Timeout:       time.Hour,
			MaxRefresh:    time.Hour,
			IdentityKey:   IdentityKey,
			Authenticator: authenticator,
			PayloadFunc:   payloadFunc,
			TokenLookup:   "cookie: token",
			TokenHeadName: "PasteMe",
		},
	}

	if err = AuthMiddleware.MiddlewareInit(); err != nil {
		logging.Panic("jwt middleware init failed", zap.Error(err))
	}
}
