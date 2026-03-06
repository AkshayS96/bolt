package cmd

import (
	"fmt"
	"strconv"

	"bolt/internal/utils"

	"github.com/google/uuid"
	gonanoid "github.com/jaevor/go-nanoid"
	"github.com/spf13/cobra"
)

func init() {
	uuidCmd := &cobra.Command{
		Use:   "uuid",
		Short: "Generate a UUID v4",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(uuid.New().String())
		},
	}

	uuidShortCmd := &cobra.Command{
		Use:   "short",
		Short: "Generate a short unique ID (8 chars)",
		Run: func(cmd *cobra.Command, args []string) {
			id, err := utils.RandomString(8, utils.AlphaNum)
			if err != nil {
				utils.PrintError("Failed to generate short ID: " + err.Error())
				return
			}
			fmt.Println(id)
		},
	}
	uuidCmd.AddCommand(uuidShortCmd)

	nanoidCmd := &cobra.Command{
		Use:   "nanoid",
		Short: "Generate a NanoID",
		Run: func(cmd *cobra.Command, args []string) {
			generate, err := gonanoid.Standard(21)
			if err != nil {
				utils.PrintError("Failed to generate NanoID: " + err.Error())
				return
			}
			fmt.Println(generate())
		},
	}

	cuidCmd := &cobra.Command{
		Use:   "cuid",
		Short: "Generate a CUID-like unique ID",
		Run: func(cmd *cobra.Command, args []string) {
			id, err := utils.RandomString(25, utils.AlphaLower+utils.Digits)
			if err != nil {
				utils.PrintError("Failed to generate CUID: " + err.Error())
				return
			}
			fmt.Println("c" + id)
		},
	}

	randomCmd := &cobra.Command{
		Use:   "random [length]",
		Short: "Generate a random alphanumeric string",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			length := 32
			if len(args) > 0 {
				n, err := strconv.Atoi(args[0])
				if err != nil || n <= 0 {
					utils.PrintError("Length must be a positive integer")
					return
				}
				length = n
			}
			s, err := utils.RandomString(length, utils.AlphaNum)
			if err != nil {
				utils.PrintError("Failed to generate random string: " + err.Error())
				return
			}
			fmt.Println(s)
		},
	}

	rootCmd.AddCommand(uuidCmd)
	rootCmd.AddCommand(nanoidCmd)
	rootCmd.AddCommand(cuidCmd)
	rootCmd.AddCommand(randomCmd)
}
