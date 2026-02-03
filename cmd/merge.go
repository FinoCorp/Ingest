/*
Copyright © 2026 FinoCorp (FinochioMatias)
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"sort"
	"encoding/csv"
	"io"

	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)
var (
	mergeOutName string
	mergeSheet string
	mergeFileType string
	
	mergeShort = "Merge all the files in a given folder into one master file."
	mergeLong = `The merge command takes a folder as an input, takes all the files inside of it
and merges all of those files into one master file.

Arguments:
<folder>    Path to the folder containing all the excel files.
	`
)

var mergeCmd = &cobra.Command{
	Use:   "merge <folder>",
	Short: mergeShort,
	Long: mergeLong,
	Args: cobra.ExactArgs(1),
	RunE: mergeFunc,
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	mergeCmd.Flags().StringVarP(&mergeFileType, "filetype", "f", "xlsx", "input file type (currently supporting csv and xlsx)")
	mergeCmd.Flags().StringVarP(&mergeOutName, "out", "o", "master_file.xlsx", "output .xlsx filename.")
	mergeCmd.Flags().StringVarP(&mergeSheet, "sheet", "s", "Sheet1", "name of the output sheet.")
}

func mergeFunc(cmd *cobra.Command, args []string) error {
	folder := args[0]

	info, err := os.Stat(folder)
	
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New("<folder> must be a directory.")
	}

	ft := strings.ToLower(strings.TrimSpace(mergeFileType))
	if ft == "" {
		ft = "xlsx"
	}

	if ft != "xlsx" && ft != "csv" {
		return errors.New("--filetype currently only support 'xlsx' or 'csv'")
	}
	
	files, err := listFilesByExt(folder, ft)

	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("no '.xlsx' files found in the provided folder.")
	}

	outPath := strings.TrimSpace(mergeOutName)
	if outPath == "" {
		return errors.New("--out flag cannot be empty.")
	}
	
	if !strings.HasSuffix(strings.ToLower(outPath), ".xlsx") {
		outPath += ".xlsx" // MATI: check if this is actually ok or if its better to just skip all files that are not xl files.
	}

	master := excelize.NewFile()
	defaultSheet := master.GetSheetName(0)

	if strings.TrimSpace(mergeSheet) == "" {
		mergeSheet = "Sheet1"
	}
	if defaultSheet != mergeSheet {
		master.SetSheetName(defaultSheet, mergeSheet)
	}

	currentRow := 1

	for fileIndex, path := range files {
		var rows [][]string

		if ft == "xlsx" {
			src, err := excelize.OpenFile(path)
			if err != nil {
				return fmt.Errorf("open %s: %w", path, err)
			}

			sheetName := src.GetSheetName(0)
			rows, err = src.GetRows(sheetName)
			_ = src.Close()

			if err != nil {
				return fmt.Errorf("read rows %s: %w", path, err)
			}
		} else {
			var err error
			rows, err = readCsvRows(path)
			if err != nil {
				return fmt.Errorf("read csv %s: %w", path, err)
			}
		}

		if len(rows) == 0 {
			continue
		}

		start := 0
		if fileIndex > 0 {
			start = 1 // skip header
		}

		for i := start; i < len(rows); i++ {
			row := rows[i]
			if len(row) == 0 {
				continue
			}

			cell, err := excelize.CoordinatesToCellName(1, currentRow)
			if err != nil {
				return err
			}

			values := make([]interface{}, len(row))
			for j := range row {
				values[j] = row[j]
			}

			if err := master.SetSheetRow(mergeSheet, cell, &values); err != nil {
				return fmt.Errorf("write row %d: %w", currentRow, err)
			}

			currentRow++
		}
	}
	

	if err := master.SaveAs(outPath); err != nil {
		return err
	}

	fmt.Printf("Merged %d file(s) into %s\n", len(files), outPath)
	return nil
}

func listXLSXFiles(folder string) ([]string, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()

		if strings.HasSuffix(strings.ToLower(name), ".xlsx") {
			paths = append(paths, filepath.Join(folder, name))
		}
	}

	sort.Strings(paths)
	return paths, nil
}

func listFilesByExt(folder string, ext string) ([]string, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	paths := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()

		if strings.HasPrefix(name, "~$") {
			continue // MATI: this ignores the files already locked by excel because they start with '~$"
		}

		if strings.HasSuffix(strings.ToLower(name), ext) {
			paths = append(paths, filepath.Join(folder, name))
		}
	}

	sort.Strings(paths)
	return paths, nil
}

func readCsvRows(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	rows := make([][]string, 0)

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		rows = append(rows, rec)
	}

	return rows, nil
}
