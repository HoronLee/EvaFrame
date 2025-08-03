//go:build wireinject
// +build wireinject

package app

import (
	"evaframe/internal/dao/gorm"
	"evaframe/internal/handler"
	"evaframe/internal/service"
	"evaframe/pkg/config"
	"evaframe/pkg/database"
	"evaframe/pkg/jwt"
	"evaframe/pkg/logger"
	"evaframe/pkg/middleware"
	"evaframe/pkg/validator"

	"github.com/google/wire"
)

// NewConfigWithPath 创建一个包装函数来接收configPath参数
func NewConfigWithPath(configPath string) (*config.Config, error) {
	return config.NewConfig(configPath)
}

// InitializeApp 使用Wire进行依赖注入
func InitializeApp(configPath string) (*Application, func(), error) {
	panic(wire.Build(
		// 配置
		NewConfigWithPath,

		// 基础设施
		logger.ProviderSet,
		database.ProviderSet,
		jwt.ProviderSet,
		validator.ProviderSet,
		middleware.ProviderSet,

		// 数据访问层
		gorm.ProviderSet,

		// 服务层
		service.ProviderSet,

		// 控制器层
		handler.ProviderSet,

		// 应用程序
		NewApplication,
	))
}
