/**********************************************************************************************
*
*   IngestCLI - Command-line tool for data processing.
*   LICENSE:
*       Mozilla Public License 2.0
*
*   Copyright © 2026 FinoCorp
*
*   This Source Code Form is subject to the terms of the Mozilla Public
*   License, v. 2.0. If a copy of the MPL was not distributed with this
*   file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
**********************************************************************************************/

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
