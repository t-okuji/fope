/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

func extractURLsFromFile(inputFile, outputFile string) error {
	// 入力ファイルを開く
	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	// 出力ファイルを作成
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	// urlを抜き出す正規表現
	re := regexp.MustCompile(`\[.+?\]\((.+?)\)`)

	for scanner.Scan() {
		line := scanner.Text()
		// 空白行をスキップする
		if line == "" {
			continue
		}
		match := re.FindStringSubmatch(line)
		matches := re.FindAllString(line, -1)
		if len(matches) > 1 {
			return fmt.Errorf("more links were found")
		}
		if len(match) > 1 {
			url := match[1]
			valid, err := isValidURL(url)
			if err != nil {
				return err // エラーを返す
			}
			if valid {
				// URLを出力ファイルに書き込む
				_, err := fmt.Fprintln(output, url)
				if err != nil {
					return fmt.Errorf("failed to write URL to output file: %w", err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading input file: %w", err)
	}

	return nil
}

// URLの形式を検証する
func isValidURL(url string) (bool, error) {
	pattern := `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(url) {
		return false, fmt.Errorf("invalid URL: %s", url)
	}

	return true, nil
}

// geturlCmd represents the geturl command
var geturlCmd = &cobra.Command{
	Use:   "geturl",
	Short: "extract url",
	Long:  `extract url from markdown format`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == args[1] {
			log.Fatal("same file")
		}
		err := extractURLsFromFile(args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("URLs extracted and saved successfully.")
	},
}

func init() {
	rootCmd.AddCommand(geturlCmd)
}
