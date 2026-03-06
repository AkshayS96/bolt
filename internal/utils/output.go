package utils

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Bold    = color.New(color.Bold)
	Success = color.New(color.FgGreen, color.Bold)
	Error   = color.New(color.FgRed, color.Bold)
	Warning = color.New(color.FgYellow, color.Bold)
	Info    = color.New(color.FgCyan)
	Dim     = color.New(color.Faint)
	Key     = color.New(color.FgMagenta, color.Bold)
)

// PrintKeyValue prints a key-value pair with colored key
func PrintKeyValue(key, value string) {
	Key.Printf("%-15s", key)
	fmt.Println(value)
}

// PrintSuccess prints a success message
func PrintSuccess(msg string) {
	Success.Println("✓ " + msg)
}

// PrintError prints an error message
func PrintError(msg string) {
	Error.Println("✗ " + msg)
}

// PrintWarning prints a warning message
func PrintWarning(msg string) {
	Warning.Println("⚠ " + msg)
}
