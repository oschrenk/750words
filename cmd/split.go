package cmd

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(splitCmd)
}

var formatFlag Format

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split exports",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0]
		outputDir := args[1]

		if _, err := os.Stat(inputPath); errors.Is(err, os.ErrNotExist) {
			log.Fatal("input file can't be found")
		}

		if _, err := os.Stat(outputDir); errors.Is(err, os.ErrNotExist) {
			log.Fatal("outputDir can't be found")
		}

		if formatFlag == OldFormat {
			parseOld(inputPath, outputDir)
		}

		if formatFlag == NewFormat {
			parseNew(inputPath, outputDir)
		}

	},
}

func parseOld(inputPath string, outputDir string) {
	const OldDelim = "------ ENTRY ------"
	const OldDatePrefix = "Date:"
	const EndMetaPrefix = "Minutes:"

	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := []string{}
	insideMeta := false
	ignoringEmpty := false
	date := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// trim whitespace for good measure
		line := strings.TrimSpace(scanner.Text())
		// 750 words export is broken, replace some german chars
		line = strings.Replace(line, "Ã¤", "ä", -1)
		line = strings.Replace(line, "Ã¼", "ü", -1)
		line = strings.Replace(line, "â‚¬", "€", -1)
		line = strings.Replace(line, "Ã¶", "ö", -1)

		// found a new entry
		if line == OldDelim {
			if lines != nil {
				write(date, lines, outputDir)
			}

			// reset
			lines = []string{}
			insideMeta = true
			date = ""
			ignoringEmpty = false
			continue
		}

		// parse date
		if strings.HasPrefix(line, OldDatePrefix) {
			s := strings.Split(line, ":")
			date = strings.TrimSpace(s[1])
			continue
		}

		// end of meta data found
		if strings.HasPrefix(line, EndMetaPrefix) {
			insideMeta = false
			ignoringEmpty = true
			continue
		}

		// stop ignoring empty lines
		if line != "" && ignoringEmpty {
			ignoringEmpty = false
		}

		// ignore empty lines after meta data
		if line == "" && ignoringEmpty {
			continue
		}

		// add line if not in meta
		if !insideMeta && !ignoringEmpty {
			lines = append(lines, line)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	//  write last entry
	if len(lines) > 0 && date != "" {
		write(date, lines, outputDir)
	}
}

func parseNew(inputPath string, outputDir string) {
	const NewDelim = "===== ENTRY ====="
	const NewDatePrefix = "=== DATE:"
	const EndMetaPrefix = "=== BODY"

	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := []string{}
	insideMeta := false
	ignoringEmpty := false
	date := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// trim whitespace for good measure
		line := strings.TrimSpace(scanner.Text())
		// 750 words export is broken, replace some german chars
		line = strings.Replace(line, "Ã¤", "ä", -1)
		line = strings.Replace(line, "Ã¼", "ü", -1)
		line = strings.Replace(line, "â‚¬", "€", -1)
		line = strings.Replace(line, "Ã¶", "ö", -1)

		// found a new entry
		if line == NewDelim {
			if lines != nil {
				write(date, lines, outputDir)
			}

			// reset
			lines = []string{}
			insideMeta = true
			date = ""
			ignoringEmpty = false
			continue
		}

		// parse date
		if strings.HasPrefix(line, NewDatePrefix) {
			s := strings.Split(line, ":")
			date = strings.TrimSpace(strings.Trim(s[1], "="))
			continue
		}

		// end of meta data found
		if strings.HasPrefix(line, EndMetaPrefix) {
			insideMeta = false
			ignoringEmpty = true
			continue
		}

		// stop ignoring empty lines
		if line != "" && ignoringEmpty {
			ignoringEmpty = false
		}

		// ignore empty lines after meta data
		if line == "" && ignoringEmpty {
			continue
		}

		// add line if not in meta
		if !insideMeta && !ignoringEmpty {
			lines = append(lines, line)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	//  write last entry
	if len(lines) > 0 && date != "" {
		write(date, lines, outputDir)
	}
}

func write(date string, lines []string, directory string) {
	path := directory + "/" + date + ".md"

	// create file
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, line := range lines {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func init() {
	splitCmd.Flags().Var(&formatFlag, "format", `format. allowed: "old", "new"`)
	splitCmd.MarkFlagRequired("format")
}
