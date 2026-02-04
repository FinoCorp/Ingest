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
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"io"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	normalizeShort = "Normalize a file with a fixed set of rules."
	normalizeLong = `The normalize command takes a CSV file and transforms it into a
clean, predictable format.

Arguments:
  <file>    Path to the input CSV file to normalize. The file must contain a header row.

Default behaviour:
- ouput headers are the input headers normalized (lowercase, trimmed, no spaces and no special characters.
- values are cleaned (trim, remove tabs/newlines, no spaces.

You can change this behaviour with the '--schema' and '--map' flags.
`
	
	normalizeOut string
	normalizeStrict bool
	
	// MATI: some sort of generic data input and output control, still have to figure this out.
	normalizeSchema string 
	normalizeMap string
)

var normalizeCmd = &cobra.Command{
	Use:   "normalize <arguments>",
	Short: normalizeShort,
	Long: normalizeLong,
	Args: cobra.MaximumNArgs(1),
	RunE: normalizeFunc,
}

func init() {
	rootCmd.AddCommand(normalizeCmd)
	
	normalizeCmd.Flags().StringVarP(
		&normalizeOut,
		"out",
		"o",
		"",
		"output file path (default: <input>_normalized.csv",
	)

	normalizeCmd.Flags().BoolVar(
		&normalizeStrict,
		"strict",
		false,
		"fails when schema columns cannot be sourced from the input file (only works when the '--schema' flag is called)",
	)

	normalizeCmd.Flags().StringVar(
		&normalizeSchema,
		"schema",
		"",
		"comma-separated output headers (e.g. 'id, date, amount, description')",
	)

	normalizeCmd.Flags().StringVar(
		&normalizeMap,
		"map",
		"",
		"comma-separated renames (input=output). e.g. 'transaction_date = date'.",
	)
}

func normalizeFunc(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer inFile.Close()

	reader := csv.NewReader(inFile)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) < 2 {
		return errors.New("input file must contain a header and at least one data row")
	}

	rawHeader := records[0]

	inHeaders, inIndex := normalizeAndIndexHeaders(rawHeader)

	renameMap, err := parseRenameMap(normalizeMap)
	if err != nil {
		return err
	}

	// use --schema flag if user privaded it.
	var outHeaders []string
	if strings.TrimSpace(normalizeSchema) == "" {
		outHeaders = applyRenamesToHeaders(inHeaders, renameMap)
	} else {
		outHeaders = parseSchema(normalizeSchema)
		if len(outHeaders) == 0 {
			return errors.New("no valid headers we found in the '--schema' provided")
		}
	}

	outToInIndex := buildOutputToInputIndex(inHeaders, inIndex, renameMap)

	if normalizeStrict && strings.TrimSpace(normalizeSchema) != "" {
		if err := enforceSchemaSourced(outHeaders, outToInIndex); err != nil {
			return err
		}
	}

	outRecords := make([][]string, 0, len(records))
	outRecords = append(outRecords, outHeaders)

	for _, row := range records[1:] {
		outRow := make([]string, len(outHeaders))

		for i, outH := range outHeaders {
			if idx, ok := outToInIndex[outH]; ok && idx < len(row) {
				outRow[i] = cleanValue(row[idx])
			} else {
				outRow[i] = ""
			}
		}

		outRecords = append(outRecords, outRow)
	}

	if normalizeOut == "" {
		ext := filepath.Ext(inputPath)
		base := strings.TrimSuffix(inputPath, ext)
		normalizeOut = base + "_normalized.csv"
	}

	return writeCSV(normalizeOut, outRecords)
}

func normalizeAndIndexHeaders(headers []string) ([]string, map[string]int) {
	seen := make(map[string]int)
	out := make([]string, 0, len(headers))
	index := make(map[string]int)

	for i, h := range headers {
		n := normalizeHeaderName(h)
		if n == "" {
			n = "column"
		}

		if count, ok := seen[n]; ok {
			count++
			seen[n] = count
			n = fmt.Sprintf("%s_%d", n, count)
		} else {
			seen[n] = 1
		}

		out = append(out, n)

		index[n] = i
	}

	return out, index
}

func normalizeHeaderName(h string) string {
	h = strings.ToLower(h)
	h = strings.TrimSpace(h)
	h = strings.ReplaceAll(h, "_", " ")
	h = strings.Join(strings.Fields(h), " ")

	return h
}

func parseSchema(schema string) []string {
	parts := strings.Split(schema, ",")
	seen := make(map[string]int)
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		n := normalizeHeaderName(p)
		if n == "" {
			continue
		}

		if count, ok := seen[n]; ok {
			count++
			seen[n] = count
			n = fmt.Sprintf("%s_%d", n, count)
		} else {
			seen[n] = 1
		}

		out = append(out, n)
	}

	return out
}

func parseRenameMap(raw string) (map[string]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]string{}, nil
	}

	m := make(map[string]string)
	pairs := strings.Split(raw, ",") //lol

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid '--map' entry: %q (expected input=output)", pair)
		}

		from := normalizeHeaderName(parts[0])
		to := normalizeHeaderName(parts[1])

		if from == "" || to == "" {
			return nil, fmt.Errorf("invalid '--map' entry: %q (empty input/output)", pair)
		}

		m[from] = to
	}

	return m, nil
}

func applyRenamesToHeaders(inHeaders []string, rename map[string]string) []string {
	out := make([]string, len(inHeaders))
	
	for i, h := range inHeaders {
		if to, ok := rename[h]; ok {
			out[i] = to
		} else {
			out[i] = h
		}
	}
	
	// dedupe again in case there are collisions still
	return dedupeHeaders(out)
}

func dedupeHeaders(headers []string) []string {
	seen := make(map[string]int)
	out := make([]string, 0, len(headers))

	for _, h := range headers {
		n := h
		if n == "" {
			n = "column"
		}

		if count, ok := seen[n]; ok {
			count++
			seen[n] = count
			n = fmt.Sprintf("%s_%d", n, count)
		} else {
			seen[n] = 1
		}

		out = append(out, n)
	}

	return out
}

func buildOutputToInputIndex(inHeaders []string, inIndex map[string]int, rename map[string]string) map [string]int {
	outToIn := make(map[string]int)

	for _, inH := range inHeaders {
		outH := inH
		if to, ok := rename[inH]; ok {
			outH = to
		}

		if _, exists := outToIn[outH]; !exists {
			if idx, ok := inIndex[inH]; ok {
				outToIn[outH] = idx
			}
		}
	}

	return outToIn
}

func enforceSchemaSourced(outHeaders []string, outToInIndex map[string]int) error {
	missing := make([]string, 0)
	for _, h := range outHeaders {
		if _, ok := outToInIndex[h]; !ok {
			missing = append(missing, h)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("strict mode: schema columns not found in input: %s", strings.Join(missing, ", "))
	}

	return nil
}

func cleanValue(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, "\t", " ")
	v = strings.ReplaceAll(v, "\n", " ")
	v = strings.ReplaceAll(v, "\r", " ")
	v = strings.Join(strings.Fields(v), " ")

	return v
}

func writeCSV(path string, records [][]string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}

	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	for _, r := range records {
		if err := writer.Write(r); err != nil {
			return err
		}
	}

	if err := writer.Error(); err != nil && err != io.EOF {
		return err
	}

	return nil
}
