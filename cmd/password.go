package cmd

import (
	"fmt"
	"math"
	"strconv"

	"bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	passwordCmd := &cobra.Command{
		Use:   "password",
		Short: "Password generation tools",
	}

	var pwLength int

	passwordGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a random password (alphanumeric)",
		Run: func(cmd *cobra.Command, args []string) {
			pw, err := utils.RandomString(pwLength, utils.AlphaNum)
			if err != nil {
				utils.PrintError("Failed to generate password: " + err.Error())
				return
			}
			fmt.Println(pw)
		},
	}
	passwordGenerateCmd.Flags().IntVarP(&pwLength, "length", "l", 20, "Password length")

	var strongLength int

	passwordStrongCmd := &cobra.Command{
		Use:   "strong",
		Short: "Generate a strong password with symbols",
		Run: func(cmd *cobra.Command, args []string) {
			pw, err := utils.RandomString(strongLength, utils.AllChars)
			if err != nil {
				utils.PrintError("Failed to generate password: " + err.Error())
				return
			}
			fmt.Println(pw)
		},
	}
	passwordStrongCmd.Flags().IntVarP(&strongLength, "length", "l", 24, "Password length")

	passwordCmd.AddCommand(passwordGenerateCmd, passwordStrongCmd)

	entropyCmd := &cobra.Command{
		Use:   "entropy <text>",
		Short: "Calculate password entropy (bits)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pw := args[0]
			charsetSize := estimateCharsetSize(pw)
			entropy := float64(len(pw)) * math.Log2(float64(charsetSize))

			utils.PrintKeyValue("Password:", pw)
			utils.PrintKeyValue("Length:", strconv.Itoa(len(pw)))
			utils.PrintKeyValue("Charset size:", strconv.Itoa(charsetSize))
			utils.PrintKeyValue("Entropy:", fmt.Sprintf("%.1f bits", entropy))

			switch {
			case entropy < 28:
				utils.PrintError("Very weak ✗")
			case entropy < 36:
				utils.PrintWarning("Weak")
			case entropy < 60:
				utils.PrintWarning("Reasonable")
			case entropy < 128:
				utils.PrintSuccess("Strong")
			default:
				utils.PrintSuccess("Very strong ✓")
			}
		},
	}

	rootCmd.AddCommand(passwordCmd, entropyCmd)
}

func estimateCharsetSize(pw string) int {
	hasLower, hasUpper, hasDigit, hasSymbol := false, false, false, false
	for _, r := range pw {
		switch {
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= '0' && r <= '9':
			hasDigit = true
		default:
			hasSymbol = true
		}
	}
	size := 0
	if hasLower {
		size += 26
	}
	if hasUpper {
		size += 26
	}
	if hasDigit {
		size += 10
	}
	if hasSymbol {
		size += 32
	}
	if size == 0 {
		size = 26
	}
	return size
}
