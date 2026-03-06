package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AkshayS96/bolt/internal/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	diffCmd := &cobra.Command{
		Use:   "diff <file1> <file2>",
		Short: "Show differences between two files",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			lines1, err := readLines(args[0])
			if err != nil {
				utils.PrintError("Cannot read " + args[0] + ": " + err.Error())
				return
			}
			lines2, err := readLines(args[1])
			if err != nil {
				utils.PrintError("Cannot read " + args[1] + ": " + err.Error())
				return
			}

			green := color.New(color.FgGreen)
			red := color.New(color.FgRed)
			dimf := color.New(color.Faint)
			accent := color.New(color.FgCyan, color.Bold)

			accent.Printf("  --- %s\n", args[0])
			accent.Printf("  +++ %s\n", args[1])

			// Find which lines differ
			maxLen := len(lines1)
			if len(lines2) > maxLen {
				maxLen = len(lines2)
			}

			contextLines := 3
			differs := make([]bool, maxLen)
			hasDiff := false
			for i := 0; i < maxLen; i++ {
				l1, l2 := "", ""
				if i < len(lines1) {
					l1 = lines1[i]
				}
				if i < len(lines2) {
					l2 = lines2[i]
				}
				if l1 != l2 {
					differs[i] = true
					hasDiff = true
				}
			}

			if !hasDiff {
				fmt.Println()
				utils.PrintSuccess("Files are identical")
				return
			}

			// Show only changed lines with context
			fmt.Println()
			lastPrinted := -1
			for i := 0; i < maxLen; i++ {
				if !isNearDiff(i, differs, contextLines) {
					continue
				}

				// Print separator if there's a gap
				if lastPrinted >= 0 && i > lastPrinted+1 {
					dimf.Println("  ···")
				}
				lastPrinted = i

				l1, l2 := "", ""
				if i < len(lines1) {
					l1 = lines1[i]
				}
				if i < len(lines2) {
					l2 = lines2[i]
				}

				if !differs[i] {
					dimf.Printf("  %4d  %s\n", i+1, l1)
				} else {
					if i < len(lines1) {
						red.Printf("  %4d - %s\n", i+1, l1)
					}
					if i < len(lines2) {
						green.Printf("  %4d + %s\n", i+1, l2)
					}
				}
			}
			fmt.Println()
		},
	}

	rootCmd.AddCommand(diffCmd)
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func isNearDiff(index int, differs []bool, context int) bool {
	start := index - context
	if start < 0 {
		start = 0
	}
	end := index + context
	if end >= len(differs) {
		end = len(differs) - 1
	}
	for i := start; i <= end; i++ {
		if differs[i] {
			return true
		}
	}
	return false
}
