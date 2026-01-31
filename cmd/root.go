/* Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"
	"fmt"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
cfgFile string

rootCmd = &cobra.Command{
	Use: "ingest",
	Short: "Clean data from files",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
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

func initializeConfig(cmd *cobra.Command) error {
	viper.SetEnvPrefix("INGEST-CLI")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()

		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "./ingest-cli")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError

		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}
	
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	fmt.Println("Config initialized. Using config file: ", viper.ConfigFileUsed())
	return nil
}
