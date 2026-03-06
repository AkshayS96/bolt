package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	hashCmd := &cobra.Command{
		Use:   "hash",
		Short: "Hashing utilities (md5, sha1, sha256)",
	}

	hashMD5Cmd := &cobra.Command{
		Use:   "md5 <text>",
		Short: "Generate MD5 hash",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Printf("%x\n", md5.Sum([]byte(input)))
		},
	}

	hashSHA1Cmd := &cobra.Command{
		Use:   "sha1 <text>",
		Short: "Generate SHA1 hash",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Printf("%x\n", sha1.Sum([]byte(input)))
		},
	}

	hashSHA256Cmd := &cobra.Command{
		Use:   "sha256 <text>",
		Short: "Generate SHA256 hash",
		Run: func(cmd *cobra.Command, args []string) {
			input, ok := utils.GetInputFromArgsOrStdin(args, 0)
			if !ok {
				utils.PrintError("No input provided")
				return
			}
			fmt.Printf("%x\n", sha256.Sum256([]byte(input)))
		},
	}

	hashFileCmd := &cobra.Command{
		Use:   "file <path>",
		Short: "Generate SHA256 hash of a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.Open(args[0])
			if err != nil {
				utils.PrintError("Cannot open file: " + err.Error())
				return
			}
			defer f.Close()

			h := sha256.New()
			if _, err := io.Copy(h, f); err != nil {
				utils.PrintError("Failed to read file: " + err.Error())
				return
			}
			utils.PrintKeyValue("SHA256:", fmt.Sprintf("%x", h.Sum(nil)))
			utils.PrintKeyValue("File:", args[0])
		},
	}

	hashCmd.AddCommand(hashMD5Cmd, hashSHA1Cmd, hashSHA256Cmd, hashFileCmd)
	rootCmd.AddCommand(hashCmd)
}
