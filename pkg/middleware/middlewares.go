// Package middleware provides gin middleware.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ProviderSet is a provider set for middlewares.
var ProviderSet = wire.NewSet(
	NewMiddlewares,
	NewLoggerMiddleware,
	NewRecoveryMiddleware,
	NewAuthMiddleware,
)

// AuthMiddleware is a custom type for auth middleware.
type AuthMiddleware gin.HandlerFunc

// LoggerMiddleware is a custom type for logger middleware.
type LoggerMiddleware gin.HandlerFunc

// RecoveryMiddleware is a custom type for recovery middleware.
type RecoveryMiddleware gin.HandlerFunc

// Middlewares contains all middlewares.
type Middlewares struct {
	Logger   LoggerMiddleware
	Auth     AuthMiddleware
	Recovery RecoveryMiddleware
}

// NewMiddlewares creates a new Middlewares container.
func NewMiddlewares(logger LoggerMiddleware, recovery RecoveryMiddleware, auth AuthMiddleware) *Middlewares {
	return &Middlewares{
		Logger:   logger,
		Auth:     auth,
		Recovery: recovery,
	}
}
