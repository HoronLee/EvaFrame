package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show application version",
	Long:  `Display the current version of the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version: %s\n", Name, Version)
	},
}
