/* Copyright © 2026 FinoCorp (FinochioMatias)
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
cfgFile string

rootCmd = &cobra.Command{
	Use: "ingest",
	Short: "Clean data from files",
}

)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.myapp.yaml)")
}
