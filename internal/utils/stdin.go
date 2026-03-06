package utils

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ReadStdin reads all data from stdin if it's being piped.
// Returns the trimmed string and whether data was available.
func ReadStdin() (string, bool) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", false
	}
	if info.Mode()&os.ModeCharDevice != 0 {
		return "", false // not piped
	}
	reader := bufio.NewReader(os.Stdin)
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(data)), true
}

// GetInputFromArgsOrStdin gets input from args first, then falls back to stdin.
// Returns the input string and whether input was found.
func GetInputFromArgsOrStdin(args []string, index int) (string, bool) {
	if len(args) > index {
		return args[index], true
	}
	return ReadStdin()
}

// ReadFileOrStdin reads from a file path if provided, otherwise reads from stdin.
func ReadFileOrStdin(args []string, index int) (string, error) {
	if len(args) > index {
		data, err := os.ReadFile(args[index])
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	input, ok := ReadStdin()
	if !ok {
		return "", io.ErrUnexpectedEOF
	}
	return input, nil
}
