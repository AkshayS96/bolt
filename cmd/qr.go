package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

func init() {
	qrCmd := &cobra.Command{
		Use:   "qr <text>",
		Short: "Generate a QR code in the terminal",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			text := strings.Join(args, " ")
			invert, _ := cmd.Flags().GetBool("invert")

			qr, err := qrcode.New(text, qrcode.Medium)
			if err != nil {
				fmt.Printf("Error generating QR code: %s\n", err)
				return
			}

			qr.DisableBorder = false
			art := qr.ToSmallString(invert)

			accent := color.New(color.FgCyan, color.Bold)
			dimC := color.New(color.Faint)

			fmt.Println()
			accent.Println("  ╔══════════════════════════════════════════╗")
			accent.Println("  ║              📱 QR Code                  ║")
			accent.Println("  ╚══════════════════════════════════════════╝")
			fmt.Println()

			for _, line := range strings.Split(art, "\n") {
				if line != "" {
					fmt.Println("    " + line)
				}
			}

			fmt.Println()
			dimC.Print("  Content: ")
			// Truncate long text for display
			display := text
			if len(display) > 50 {
				display = display[:47] + "..."
			}
			fmt.Println(display)
			fmt.Println()
			accent.Println("  Scan with your phone camera to open! 📷")
			fmt.Println()
		},
	}

	qrCmd.Flags().Bool("invert", false, "Invert colors (white on black)")
	rootCmd.AddCommand(qrCmd)
}
