package main

import (
	"evaframe/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	db, err := gorm.Open(mysql.Open("user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移表结构
	db.AutoMigrate(&models.User{})

	// 创建 Gen 实例
	g := gen.NewGenerator(gen.Config{
		OutPath:        "./internal/dao/gormgen/query",
		Mode:           gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:  true,
		FieldCoverable: false,
		FieldSignable:  false,
	})

	// 使用现有的数据库连接
	g.UseDB(db)

	// 生成基础的 CRUD 方法
	g.ApplyBasic(models.User{})

	// 生成代码
	g.Execute()
}
