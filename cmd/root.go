package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
