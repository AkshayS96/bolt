package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	interactiveCmd := &cobra.Command{
		Use:     "interactive",
		Aliases: []string{"i"},
		Short:   "Start interactive REPL mode",
		Run: func(cmd *cobra.Command, args []string) {
			runInteractive()
		},
	}

	rootCmd.AddCommand(interactiveCmd)
}

func getCommandNames() []string {
	var names []string
	for _, cmd := range rootCmd.Commands() {
		if cmd.IsAvailableCommand() && cmd.Name() != "interactive" {
			names = append(names, cmd.Name())
			// Add subcommand completions
			for _, sub := range cmd.Commands() {
				if sub.IsAvailableCommand() {
					names = append(names, cmd.Name()+" "+sub.Name())
				}
			}
		}
	}
	sort.Strings(names)
	return names
}

func buildCompleter() *readline.PrefixCompleter {
	var items []readline.PrefixCompleterInterface

	for _, cmd := range rootCmd.Commands() {
		if !cmd.IsAvailableCommand() || cmd.Name() == "interactive" {
			continue
		}
		if cmd.HasAvailableSubCommands() {
			var subItems []readline.PrefixCompleterInterface
			for _, sub := range cmd.Commands() {
				if sub.IsAvailableCommand() {
					subItems = append(subItems, readline.PcItem(sub.Name()))
				}
			}
			items = append(items, readline.PcItem(cmd.Name(), subItems...))
		} else {
			items = append(items, readline.PcItem(cmd.Name()))
		}
	}

	// Add built-in commands
	items = append(items, readline.PcItem("help"))
	items = append(items, readline.PcItem("clear"))
	items = append(items, readline.PcItem("exit"))

	return readline.NewPrefixCompleter(items...)
}

func runInteractive() {
	accent := color.New(color.FgCyan, color.Bold)
	boldC := color.New(color.Bold)
	dimC := color.New(color.Faint)

	completer := buildCompleter()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[36m  bolt› \033[0m",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		HistoryFile:     os.TempDir() + "/bolt_history",
	})
	if err != nil {
		fmt.Println("Error starting interactive mode:", err)
		return
	}
	defer rl.Close()

	// Welcome banner
	fmt.Println()
	accent.Println("  ╔══════════════════════════════════════════╗")
	accent.Print("  ║  ⚡ ")
	boldC.Print("Bolt Interactive Mode")
	accent.Println("               ║")
	accent.Println("  ╚══════════════════════════════════════════╝")
	fmt.Println()
	dimC.Println("  Features: ↑↓ history • Tab completion • 'help' for commands")
	dimC.Println("  Type 'exit' or Ctrl+D to quit")
	fmt.Println()

	cmdCount := 0
	startTime := time.Now()

	for {
		line, err := rl.Readline()
		if err != nil { // EOF or interrupt
			printExitMessage(dimC, cmdCount, startTime)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Exit commands
		if line == "exit" || line == "quit" || line == "q" {
			printExitMessage(dimC, cmdCount, startTime)
			return
		}

		// Clear screen
		if line == "clear" || line == "cls" {
			fmt.Print("\033[2J\033[H")
			continue
		}

		// Help command
		if line == "help" {
			fmt.Println()
			rootHelp(rootCmd)
			continue
		}

		// Parse the input into args
		args := parseArgs(line)
		if len(args) == 0 {
			continue
		}

		// Handle subcommand help
		if len(args) >= 2 && (args[len(args)-1] == "--help" || args[len(args)-1] == "-h") {
			sub, _, err := rootCmd.Find(args[:len(args)-1])
			if err == nil && sub != rootCmd {
				fmt.Println()
				subcommandHelp(sub)
				continue
			}
		}

		// Execute the command
		newRoot := *rootCmd
		newRoot.Use = ""
		newRoot.SetArgs(args)
		newRoot.SilenceErrors = true
		newRoot.SilenceUsage = true

		fmt.Println()
		if err := newRoot.Execute(); err != nil {
			errC := color.New(color.FgRed)
			errC.Printf("  ✗ %s\n", err.Error())
		}
		fmt.Println()

		cmdCount++
	}
}

func printExitMessage(dimC *color.Color, cmdCount int, startTime time.Time) {
	elapsed := time.Since(startTime).Round(time.Second)
	fmt.Println()
	dimC.Printf("  👋 Bye! (%d commands in %s)\n", cmdCount, elapsed)
	fmt.Println()
}

// parseArgs handles basic shell-like argument parsing with quote support
func parseArgs(input string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := byte(0)

	for i := 0; i < len(input); i++ {
		c := input[i]
		switch {
		case inQuote:
			if c == quoteChar {
				inQuote = false
			} else {
				current.WriteByte(c)
			}
		case c == '"' || c == '\'':
			inQuote = true
			quoteChar = c
		case c == ' ' || c == '\t':
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
