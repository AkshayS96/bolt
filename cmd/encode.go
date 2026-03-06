package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"

	"bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	// --- Base64 ---
	base64Cmd := &cobra.Command{
		Use:   "base64",
		Short: "Base64 encode/decode",
	}

	base64EncodeCmd := &cobra.Command{
		Use:   "encode <text>",
		Short: "Encode text to Base64",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Println(base64.StdEncoding.EncodeToString([]byte(input)))
		},
	}

	base64DecodeCmd := &cobra.Command{
		Use:   "decode <text>",
		Short: "Decode Base64 text",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			decoded, err := base64.StdEncoding.DecodeString(input)
			if err != nil {
				// Try URL-safe encoding
				decoded, err = base64.URLEncoding.DecodeString(input)
				if err != nil {
					utils.PrintError("Invalid Base64: " + err.Error())
					return
				}
			}
			fmt.Println(string(decoded))
		},
	}

	base64Cmd.AddCommand(base64EncodeCmd, base64DecodeCmd)

	// --- URL ---
	urlCmd := &cobra.Command{
		Use:   "url",
		Short: "URL encode/decode",
	}

	urlEncodeCmd := &cobra.Command{
		Use:   "encode <text>",
		Short: "URL encode text",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Println(url.QueryEscape(input))
		},
	}

	urlDecodeCmd := &cobra.Command{
		Use:   "decode <text>",
		Short: "Decode URL-encoded text",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			decoded, err := url.QueryUnescape(input)
			if err != nil {
				utils.PrintError("Invalid URL encoding: " + err.Error())
				return
			}
			fmt.Println(decoded)
		},
	}

	urlCmd.AddCommand(urlEncodeCmd, urlDecodeCmd)

	// --- Hex ---
	hexCmd := &cobra.Command{
		Use:   "hex",
		Short: "Hex encode/decode",
	}

	hexEncodeCmd := &cobra.Command{
		Use:   "encode <text>",
		Short: "Hex encode text",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Println(hex.EncodeToString([]byte(input)))
		},
	}

	hexDecodeCmd := &cobra.Command{
		Use:   "decode <text>",
		Short: "Decode hex string",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			decoded, err := hex.DecodeString(input)
			if err != nil {
				utils.PrintError("Invalid hex: " + err.Error())
				return
			}
			fmt.Println(string(decoded))
		},
	}

	hexCmd.AddCommand(hexEncodeCmd, hexDecodeCmd)

	rootCmd.AddCommand(base64Cmd, urlCmd, hexCmd)
}
