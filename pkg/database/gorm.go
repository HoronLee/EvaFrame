// Package database 提供 GORM 数据库实例的创建和配置
package database

import (
	"evaframe/pkg/config"
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ProviderSet = wire.NewSet(NewDB)

// NewDB GORM 数据库实例 Provider
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Type {
	case "mysql":
		dialector = mysql.Open(cfg.Database.DSN)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.DSN)
	case "postgres":
		dialector = postgres.Open(cfg.Database.DSN)
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Database.Type)
	}

	gcfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志模式
		// TODO: 实现自定义Gorm日志记录器
	}

	db, err := gorm.Open(dialector, gcfg)
	if err != nil {
		return nil, err
	}
	return db, nil
}
