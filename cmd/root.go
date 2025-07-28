package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// go build -ldflags "-X cmd.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "EvaFrame"
	// Version is the version of the compiled software.
	Version string = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "evaframe",
	Short: "A modern Go web framework",
	Long:  `EvaFrame is a modern Go web framework built with Gin, GORM, and dependency injection.`,
}

var configFile string

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config/config.yaml", "config file path")
}

func Execute() {
	// 如果没有提供子命令，设置为 serve
	if len(os.Args) == 1 {
		args := append([]string{os.Args[0]}, "serve")
		os.Args = args
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
