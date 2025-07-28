package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"evaframe/internal/app"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Long:  `Start the web server with the specified configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化应用
		application, cleanup, err := app.InitializeApp(configFile)
		if err != nil {
			fmt.Printf("Failed to initialize app: %v\n", err)
			os.Exit(1)
		}
		defer cleanup()

		// 启动服务器
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", application.Config.Server.Port),
			Handler: application.Router,
		}

		// 优雅关闭
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Printf("Server failed to start: %v\n", err)
				os.Exit(1)
			}
		}()

		fmt.Printf("Server started on port %d\n", application.Config.Server.Port)

		// 等待中断信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		fmt.Println("Shutting down server...")

		// 优雅关闭，最多等待30秒
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			fmt.Printf("Server forced to shutdown: %v\n", err)
		}

		fmt.Println("Server exited")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
