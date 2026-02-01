/*
Copyright © 2026 FinoCorp (FinochioMatias)
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"strings"
//	"io"
	"log"

	"github.com/spf13/cobra"
)

var in_string = `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
` 

// normalizeCmd is the command for normalizing the values of the input file.
var normalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalize a file with a fixed set of rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := csv.NewReader(strings.NewReader(in_string))

		r.Comma = ','
		r.Comment = '#'

		records, err := r.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(records)

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
