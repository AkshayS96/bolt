package cmd

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	fileCmd := &cobra.Command{
		Use:   "file",
		Short: "File utilities (hash, size, lines, stats)",
	}

	fileHashCmd := &cobra.Command{
		Use:   "hash <path>",
		Short: "SHA256 hash of a file",
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
			fmt.Printf("%x\n", h.Sum(nil))
		},
	}

	fileSizeCmd := &cobra.Command{
		Use:   "size <path>",
		Short: "Show file size",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			info, err := os.Stat(args[0])
			if err != nil {
				utils.PrintError("Cannot stat file: " + err.Error())
				return
			}
			size := info.Size()
			utils.PrintKeyValue("Bytes:", fmt.Sprintf("%d", size))
			utils.PrintKeyValue("Human:", humanizeBytes(size))
		},
	}

	fileLinesCmd := &cobra.Command{
		Use:   "lines <path>",
		Short: "Count lines in a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.Open(args[0])
			if err != nil {
				utils.PrintError("Cannot open file: " + err.Error())
				return
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			count := 0
			for scanner.Scan() {
				count++
			}
			fmt.Println(count)
		},
	}

	fileStatsCmd := &cobra.Command{
		Use:   "stats <path>",
		Short: "Show file metadata",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			info, err := os.Stat(args[0])
			if err != nil {
				utils.PrintError("Cannot stat file: " + err.Error())
				return
			}
			utils.PrintKeyValue("Name:", info.Name())
			utils.PrintKeyValue("Size:", humanizeBytes(info.Size()))
			utils.PrintKeyValue("Permissions:", info.Mode().String())
			utils.PrintKeyValue("Modified:", info.ModTime().Format("2006-01-02 15:04:05"))
			utils.PrintKeyValue("Is Directory:", fmt.Sprintf("%v", info.IsDir()))
		},
	}

	fileCmd.AddCommand(fileHashCmd, fileSizeCmd, fileLinesCmd, fileStatsCmd)
	rootCmd.AddCommand(fileCmd)
}

func humanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
