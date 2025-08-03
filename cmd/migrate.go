package cmd

import (
	"fmt"
	"os"

	"evaframe/internal/models"
	"evaframe/pkg/config"
	"evaframe/pkg/database"
	"evaframe/pkg/logger"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run GORM auto-migration to create/update database schema.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置
		cfg, err := config.NewConfig(configFile)
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// 初始化日志记录器
		appLogger, err := logger.NewLogger(cfg)
		if err != nil {
			fmt.Printf("Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}

		// 连接数据库
		db, err := database.NewDB(cfg, appLogger)
		if err != nil {
			fmt.Printf("Failed to connect database: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Starting database migration...")

		// 自动迁移表结构
		err = db.AutoMigrate(
			&models.User{},
		)
		if err != nil {
			fmt.Printf("Migration failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Database migration completed successfully!")
	},
}
