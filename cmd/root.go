package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const version = "4.4.6"

var (
	accent = color.New(color.FgCyan, color.Bold)
	bold   = color.New(color.Bold)
	dimC   = color.New(color.Faint)
	white  = color.New(color.FgWhite, color.Bold)
)

var rootCmd = &cobra.Command{
	Use:   "bolt",
	Short: "⚡ Bolt – A developer Swiss Army knife CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("bolt v{{.Version}}\n")
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Custom help only for the root command
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd == rootCmd {
			rootHelp(cmd)
		} else {
			subcommandHelp(cmd)
		}
	})
}

func rootHelp(cmd *cobra.Command) {
	fmt.Println()
	accent.Print("  ⚡ ")
	bold.Print("BOLT")
	dimC.Println("  —  Developer Swiss Army Knife")
	fmt.Println()
	dimC.Println("  A fast, single-binary toolbox for everyday developer tasks.")
	fmt.Println()

	groups := []struct {
		icon     string
		name     string
		commands []string
	}{
		{"🔑", "ID Generators", []string{"uuid", "nanoid", "cuid", "random"}},
		{"📦", "Data & Encoding", []string{"json", "base64", "url", "hex"}},
		{"🔒", "Security", []string{"jwt", "hash", "password", "entropy"}},
		{"⏰", "Time & Date", []string{"time"}},
		{"🌐", "HTTP & Network", []string{"http", "dns", "ping", "ip", "port"}},
		{"✏️ ", "Text & Strings", []string{"slug", "case", "trim", "length", "regex"}},
		{"📁", "Files & System", []string{"file", "diff", "clip"}},
		{"🧰", "Utilities", []string{"color", "lorem", "cron", "qr", "gitignore"}},
	}

	for _, g := range groups {
		fmt.Print("  " + g.icon + " ")
		bold.Println(g.name)
		for _, cmdName := range g.commands {
			if sub, _, err := cmd.Find([]string{cmdName}); err == nil && sub != cmd {
				fmt.Print("     ")
				accent.Printf("%-14s", sub.Name())
				dimC.Printf(" %s\n", sub.Short)
			}
		}
		fmt.Println()
	}

	bold.Println("  USAGE")
	fmt.Println()
	dimC.Print("    $ ")
	fmt.Println("bolt <command> [flags]")
	dimC.Print("    $ ")
	fmt.Println("bolt <command> --help")
	fmt.Println()

	bold.Println("  EXAMPLES")
	fmt.Println()
	printExample("bolt uuid", "Generate a UUID v4")
	printExample("bolt json format data.json", "Pretty print JSON")
	printExample("bolt hash sha256 hello", "SHA256 hash")
	printExample("bolt time now", "Current timestamp")
	printExample("bolt password strong", "Generate strong password")
	printExample("echo '{\"a\":1}' | bolt json format", "Pipe support")
	fmt.Println()

	dimC.Println("  ─────────────────────────────────────────────────────")
	dimC.Print("  v" + version + "  •  ")
	accent.Println("https://github.com/AkshayS96/bolt")
	dimC.Print("  Made with ❤️  by ")
	accent.Println("https://www.x.com/__akshaysolanki")
	dimC.Println("  ─────────────────────────────────────────────────────")
	fmt.Println()
}

func subcommandHelp(cmd *cobra.Command) {
	fmt.Println()
	accent.Print("  ⚡ bolt ")
	bold.Println(cmd.Name())
	fmt.Println()

	if cmd.Long != "" {
		dimC.Println("  " + strings.ReplaceAll(strings.TrimSpace(cmd.Long), "\n", "\n  "))
	} else if cmd.Short != "" {
		dimC.Println("  " + cmd.Short)
	}
	fmt.Println()

	// Show subcommands if any
	if cmd.HasAvailableSubCommands() {
		bold.Println("  COMMANDS")
		fmt.Println()
		for _, sub := range cmd.Commands() {
			if sub.IsAvailableCommand() {
				fmt.Print("     ")
				accent.Printf("%-18s", sub.Name())
				dimC.Printf(" %s\n", sub.Short)
			}
		}
		fmt.Println()
	}

	// Usage
	bold.Println("  USAGE")
	fmt.Println()
	if cmd.HasAvailableSubCommands() {
		dimC.Print("    $ ")
		fmt.Printf("bolt %s <command> [args]\n", cmd.Name())
	} else {
		dimC.Print("    $ ")
		fmt.Printf("bolt %s\n", cmd.UseLine())
	}
	fmt.Println()

	// Flags
	if cmd.HasAvailableLocalFlags() {
		bold.Println("  FLAGS")
		fmt.Println()
		fmt.Print("  " + strings.ReplaceAll(strings.TrimRight(cmd.LocalFlags().FlagUsages(), "\n"), "\n", "\n  "))
		fmt.Println()
		fmt.Println()
	}

	// Examples from subcommands
	if cmd.HasAvailableSubCommands() {
		bold.Println("  EXAMPLES")
		fmt.Println()
		for _, sub := range cmd.Commands() {
			if sub.IsAvailableCommand() {
				dimC.Print("    $ ")
				fmt.Printf("bolt %s %s\n", cmd.Name(), sub.Use)
			}
		}
		fmt.Println()
	}
}

func printExample(command, description string) {
	dimC.Print("    $ ")
	fmt.Printf("%-42s", command)
	dimC.Printf(" # %s\n", description)
}
