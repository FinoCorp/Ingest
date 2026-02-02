/*
Copyright © 2026 FinoCorp (FinochioMatias)
*/
package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var short_description = "Normalize a file with a fixed set of rules"

var long_description = `The normalize command takes a CSV input file and transforms its raw data into a
clean, predictable format using a fixed set of normalization rules.

Arguments:
  <file>    Path to the input CSV file to normalize. The file must contain a header row.

The command reads the input file, applies normalization rules to each row, and outputs
the normalized data in CSV format.`

var correct_headers = []string{"date", "amount", "description"}

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
	validated_headers := make(map[string]int)
	
	if len(args) == 0 {
		return errors.New("<file> argument missing.")
	}

	records, err := get_data_from_input(args[0])

	if err != nil {
		return err
	}

	if len(records) == 0 {
		return errors.New("<file> provided is empty.")
	}

	if len(records[0]) == 0 {
		return errors.New("<file> has no header row.")
	}

	headers := records[0] // header should always be first index of the csv data

	validated_headers, err = validate_header(headers)

	if err != nil {
		return err
	}

	fmt.Println(validated_headers)

	return nil
}

func get_data_from_input(path string) ([][]string, error ){
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

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

func validate_header(csv_records []string) (map[string]int, error) {
	if len(csv_records) == 0 {
		return nil, errors.New("Header row cannot be empty")
	}
	
	header_check_map := make(map[string]int)

	for i, value := range csv_records {
		value = strings.TrimSpace(value)
		value = strings.ToLower(value)

		header_check_map[value] = i
	}

	for _, cr := range correct_headers {
		cr = strings.TrimSpace(cr)
		cr = strings.ToLower(cr)

		_, found := header_check_map[cr]
		if !found {
			return nil, errors.New("Headers were not found in the CSV file")
		}
	}

	return header_check_map, nil
}
