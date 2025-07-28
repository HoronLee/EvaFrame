package cmd

import (
	"fmt"
	"os"

	"evaframe/internal/models"
	"evaframe/pkg/config"
	"evaframe/pkg/database"

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

		// 连接数据库
		db, err := database.NewDB(cfg)
		if err != nil {
			fmt.Printf("Failed to connect database: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Starting database migration...")

		// 自动迁移表结构
		err = db.AutoMigrate(
			&models.User{},
			// 在这里添加其他模型
		)
		if err != nil {
			fmt.Printf("Migration failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Database migration completed successfully!")
	},
}
