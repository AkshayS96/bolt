package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/AkshayS96/bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	clipCmd := &cobra.Command{
		Use:   "clip",
		Short: "Clipboard utilities (copy, paste)",
	}

	clipCopyCmd := &cobra.Command{
		Use:   "copy [text]",
		Short: "Copy text to clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			var text string
			if len(args) > 0 {
				text = strings.Join(args, " ")
			} else {
				input, ok := utils.ReadStdin()
				if !ok {
					utils.PrintError("No input provided. Use: bolt clip copy <text> or pipe input")
					return
				}
				text = input
			}

			if err := copyToClipboard(text); err != nil {
				utils.PrintError("Failed to copy to clipboard: " + err.Error())
				return
			}
			utils.PrintSuccess(fmt.Sprintf("Copied to clipboard (%d chars)", len(text)))
		},
	}

	clipPasteCmd := &cobra.Command{
		Use:   "paste",
		Short: "Paste text from clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			text, err := pasteFromClipboard()
			if err != nil {
				utils.PrintError("Failed to read clipboard: " + err.Error())
				return
			}
			fmt.Print(text)
		},
	}

	clipCmd.AddCommand(clipCopyCmd, clipPasteCmd)
	rootCmd.AddCommand(clipCmd)
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func pasteFromClipboard() (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbpaste")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
	case "windows":
		cmd = exec.Command("powershell", "-command", "Get-Clipboard")
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
