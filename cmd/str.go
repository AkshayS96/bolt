package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	slugCmd := &cobra.Command{
		Use:   "slug <text>",
		Short: "Convert text to URL-friendly slug",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Println(slugify(input))
		},
	}

	caseCmd := &cobra.Command{
		Use:   "case",
		Short: "Convert text between cases (camel, snake, kebab, pascal)",
	}

	caseCamelCmd := &cobra.Command{
		Use:   "camel <text>",
		Short: "Convert to camelCase",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Println(toCamelCase(input, false))
		},
	}

	casePascalCmd := &cobra.Command{
		Use:   "pascal <text>",
		Short: "Convert to PascalCase",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Println(toCamelCase(input, true))
		},
	}

	caseSnakeCmd := &cobra.Command{
		Use:   "snake <text>",
		Short: "Convert to snake_case",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Println(toSnakeCase(input))
		},
	}

	caseKebabCmd := &cobra.Command{
		Use:   "kebab <text>",
		Short: "Convert to kebab-case",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Println(toKebabCase(input))
		},
	}

	caseCmd.AddCommand(caseCamelCmd, casePascalCmd, caseSnakeCmd, caseKebabCmd)

	trimCmd := &cobra.Command{
		Use:   "trim <text>",
		Short: "Trim whitespace from text",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Println(strings.TrimSpace(input))
		},
	}

	lengthCmd := &cobra.Command{
		Use:   "length <text>",
		Short: "Count characters in text",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			runes := []rune(input)
			utils.PrintKeyValue("Characters:", fmt.Sprintf("%d", len(runes)))
			utils.PrintKeyValue("Bytes:", fmt.Sprintf("%d", len(input)))
			utils.PrintKeyValue("Words:", fmt.Sprintf("%d", len(strings.Fields(input))))
		},
	}

	rootCmd.AddCommand(slugCmd, caseCmd, trimCmd, lengthCmd)
}

func splitWords(s string) []string {
	// Replace common delimiters with spaces
	re := regexp.MustCompile(`[-_\s]+`)
	s = re.ReplaceAllString(s, " ")
	// Insert space before uppercase letters (for camelCase splitting)
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) && (i+1 < len(s) && unicode.IsLower(rune(s[i+1])) || unicode.IsLower(rune(s[i-1]))) {
			result = append(result, ' ')
		}
		result = append(result, r)
	}
	return strings.Fields(strings.ToLower(string(result)))
}

func slugify(s string) string {
	words := splitWords(s)
	return strings.Join(words, "-")
}

func toCamelCase(s string, pascal bool) string {
	words := splitWords(s)
	for i, w := range words {
		if i == 0 && !pascal {
			words[i] = strings.ToLower(w)
		} else {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, "")
}

func toSnakeCase(s string) string {
	words := splitWords(s)
	return strings.Join(words, "_")
}

func toKebabCase(s string) string {
	words := splitWords(s)
	return strings.Join(words, "-")
}
