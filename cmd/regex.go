package cmd

import (
	"fmt"
	"regexp"

	"bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	regexCmd := &cobra.Command{
		Use:   "regex",
		Short: "Regex tools (match, test, extract)",
	}

	regexMatchCmd := &cobra.Command{
		Use:   "match <pattern> <text>",
		Short: "Find all regex matches in text",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			re, err := regexp.Compile(args[0])
			if err != nil {
				utils.PrintError("Invalid regex: " + err.Error())
				return
			}
			matches := re.FindAllString(args[1], -1)
			if len(matches) == 0 {
				utils.PrintWarning("No matches found")
				return
			}
			utils.PrintSuccess(fmt.Sprintf("Found %d match(es):", len(matches)))
			for i, m := range matches {
				fmt.Printf("  %d: %s\n", i+1, m)
			}
		},
	}

	regexTestCmd := &cobra.Command{
		Use:   "test <pattern> <text>",
		Short: "Test if text matches a regex pattern",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			re, err := regexp.Compile(args[0])
			if err != nil {
				utils.PrintError("Invalid regex: " + err.Error())
				return
			}
			if re.MatchString(args[1]) {
				utils.PrintSuccess("Match ✓")
			} else {
				utils.PrintError("No match ✗")
			}
		},
	}

	regexExtractCmd := &cobra.Command{
		Use:   "extract <pattern> <text>",
		Short: "Extract capture groups from text",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			re, err := regexp.Compile(args[0])
			if err != nil {
				utils.PrintError("Invalid regex: " + err.Error())
				return
			}
			matches := re.FindAllStringSubmatch(args[1], -1)
			if len(matches) == 0 {
				utils.PrintWarning("No matches found")
				return
			}
			names := re.SubexpNames()
			for i, match := range matches {
				fmt.Printf("Match %d:\n", i+1)
				for j, group := range match {
					label := fmt.Sprintf("  Group %d:", j)
					if j < len(names) && names[j] != "" {
						label = fmt.Sprintf("  %s:", names[j])
					}
					fmt.Printf("%-15s %s\n", label, group)
				}
			}
		},
	}

	regexReplaceCmd := &cobra.Command{
		Use:   "replace <pattern> <replacement> <text>",
		Short: "Replace regex matches in text",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			re, err := regexp.Compile(args[0])
			if err != nil {
				utils.PrintError("Invalid regex: " + err.Error())
				return
			}
			fmt.Println(re.ReplaceAllString(args[2], args[1]))
		},
	}

	regexCmd.AddCommand(regexMatchCmd, regexTestCmd, regexExtractCmd, regexReplaceCmd)
	rootCmd.AddCommand(regexCmd)
}
