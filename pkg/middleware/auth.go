// Package middleware provides gin middleware.
package middleware

import (
	"evaframe/pkg/jwt"
	"evaframe/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth is a factory function to create a JWT authentication middleware.
func NewAuthMiddleware(jwt *jwt.JWT) AuthMiddleware {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.Unauthorized(c, "未授权")
			c.Abort()
			return
		}

		// Real token starts after "Bearer "
		token, err := jwt.ParseToken(tokenStr[7:])
		if err != nil {
			response.Unauthorized(c, "令牌无效或已过期")
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("claims", token)
		c.Next()
	}
}
