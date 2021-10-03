package token

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JWTMiddleware struct {
	*jwt.GinJWTMiddleware
}

// MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
func (mw *JWTMiddleware) MiddlewareFunc(disableAboard bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.middlewareImpl(c, disableAboard)
	}
}

func (mw *JWTMiddleware) unauthorized(disableAboard bool, c *gin.Context, code int, message string) {
	if disableAboard {
		c.Set(mw.IdentityKey, Nobody)
		c.Next()
	} else {
		c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
		if !mw.DisabledAbort {
			c.Abort()
		}

		mw.Unauthorized(c, code, message)
	}
}

func (mw *JWTMiddleware) middlewareImpl(c *gin.Context, disableAboard bool) {
	claims, err := mw.GetClaimsFromJWT(c)
	if err != nil {
		mw.unauthorized(disableAboard, c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, c))
		return
	}

	if claims["exp"] == nil {
		mw.unauthorized(disableAboard, c, http.StatusBadRequest, mw.HTTPStatusMessageFunc(jwt.ErrMissingExpField, c))
		return
	}

	if _, ok := claims["exp"].(float64); !ok {
		mw.unauthorized(disableAboard, c, http.StatusBadRequest, mw.HTTPStatusMessageFunc(jwt.ErrWrongFormatOfExp, c))
		return
	}

	if int64(claims["exp"].(float64)) < mw.TimeFunc().Unix() {
		mw.unauthorized(disableAboard, c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(jwt.ErrExpiredToken, c))
		return
	}

	c.Set("JWT_PAYLOAD", claims)
	identity := mw.IdentityHandler(c)

	if identity != nil {
		c.Set(mw.IdentityKey, identity)
	}

	if !mw.Authorizator(identity, c) {
		mw.unauthorized(disableAboard, c, http.StatusForbidden, mw.HTTPStatusMessageFunc(jwt.ErrForbidden, c))
		return
	}

	c.Next()
}
