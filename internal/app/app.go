package app

import (
	"evaframe/internal/handler"
	"evaframe/pkg/config"
	"evaframe/pkg/logger"
	"evaframe/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Config      *config.Config
	Router      *gin.Engine
	UserHandler *handler.UserHandler
	Logger      *logger.Logger
}

func NewApplication(
	cfg *config.Config,
	userHandler *handler.UserHandler,
	mws *middleware.Middlewares,
	logger *logger.Logger,
) *Application {
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由器
	router := gin.New()
	router.Use(gin.HandlerFunc(mws.Logger))
	router.Use(gin.Recovery())

	// 注册路由
	apiV1 := router.Group("/api/v1")
	userHandler.RegisterRoutes(apiV1, gin.HandlerFunc(mws.Auth))

	return &Application{
		Config:      cfg,
		Router:      router,
		UserHandler: userHandler,
		Logger:      logger,
	}
}
