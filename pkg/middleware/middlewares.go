// Package middleware provides gin middleware.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// AuthMiddleware is a custom type for auth middleware.
type AuthMiddleware gin.HandlerFunc

// LoggerMiddleware is a custom type for logger middleware.
type LoggerMiddleware gin.HandlerFunc

// Middlewares contains all middlewares.
type Middlewares struct {
	Logger LoggerMiddleware
	Auth   AuthMiddleware
}

// NewMiddlewares creates a new Middlewares container.
func NewMiddlewares(logger LoggerMiddleware, auth AuthMiddleware) *Middlewares {
	return &Middlewares{
		Logger: logger,
		Auth:   auth,
	}
}

// ProviderSet is a provider set for middlewares.
var ProviderSet = wire.NewSet(
	NewMiddlewares,
	NewLoggerMiddleware,
	NewAuthMiddleware,
)