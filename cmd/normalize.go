/*
Copyright © 2026 FinoCorp (FinochioMatias)
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"strings"
	"os"

	"github.com/spf13/cobra"
)

var in_string = `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
` 

// normalizeCmd is the command for normalizing the values of the input file.
var normalizeCmd = &cobra.Command{
	Use:   "normalize <arguments>",
	Short: "Normalize a file with a fixed set of rules",
	Long: `The normalize command takes a CSV input file and transforms its raw data into a
clean, predictable format using a fixed set of normalization rules.

Arguments:
  <file>    Path to the input CSV file to normalize. The file must contain a header row.

The command reads the input file, applies normalization rules to each row, and outputs
the normalized data in CSV format.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("'normalize' command requires a <file> argument. Use the '--help' flag to see the usage")
		} else {
			file_path, err := os.Open(args[0])
			data, err := os.ReadFile(file_path.Name())
			if err != nil {
				return err
			}

			reader := strings.NewReader(string(data))
			csv_reader := csv.NewReader(reader)

			csv_reader.Comma = ','
			csv_reader.Comment = '#'

			records, err := csv_reader.ReadAll()
			if err != nil {
				return err
			}

			fmt.Println(records)		
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(normalizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// normalizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// normalizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
