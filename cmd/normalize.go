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

var short_description = "Normalize a file with a fixed set of rules"

var long_description = `The normalize command takes a CSV input file and transforms its raw data into a
clean, predictable format using a fixed set of normalization rules.

Arguments:
  <file>    Path to the input CSV file to normalize. The file must contain a header row.

The command reads the input file, applies normalization rules to each row, and outputs
the normalized data in CSV format.`

// normalizeCmd is the command for normalizing the values of the input file.
var normalizeCmd = &cobra.Command{
	Use:   "normalize <arguments>",
	Short: short_description,
	Long: long_description,
	Args: cobra.MaximumNArgs(1),
	RunE: process_data,
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

func process_data(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Println("<file> argument missing. See usage with the '--help' flag")
	} else {
		records, err := get_data_from_input(args[0])
		if err != nil {
			return err
		} else {
			data_ok := validate_header(records)

			if data_ok == true {
				fmt.Println(records)
			}
		}
	}

	return nil
}

func get_data_from_input(path string) ([][]string, error ){
	file, err := os.Open(path)
	data, err := os.ReadFile(file.Name())

	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(string(data))
	csv_reader := csv.NewReader(reader)

	csv_data, err := csv_reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return csv_data, nil
}

func validate_header(csv_records [][]string) bool {
	for i, headers := range csv_records {
		if i == 0 {
			//fmt.Println("Headers ROW")
			if len(headers) == 0 {
				return false
			} 
		}
	}

	return true
}
