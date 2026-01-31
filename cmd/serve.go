/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",

	RunE: func(cmd *cobra.Command, args []string) error {
		port := viper.GetInt("port")

		fmt.Println("Server starting on port: %d\\n", port)
		// server would start here

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Int("port", 8080, "Port to run the server on.")
}
