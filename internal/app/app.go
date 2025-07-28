package app

import (
	"evaframe/internal/handler"
	"evaframe/pkg/config"
	"evaframe/pkg/jwt"
	"evaframe/pkg/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Application struct {
	Config      *config.Config
	Router      *gin.Engine
	UserHandler *handler.UserHandler
	Logger      *zap.Logger
}

func NewApplication(
	cfg *config.Config,
	userHandler *handler.UserHandler,
	jwtService *jwt.JWT,
	logger *zap.Logger,
) *Application {
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由器
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 注册路由
	authMiddleware := middleware.JWTAuth(jwtService)
	userHandler.RegisterRoutes(router, authMiddleware)

	return &Application{
		Config:      cfg,
		Router:      router,
		UserHandler: userHandler,
		Logger:      logger,
	}
}
