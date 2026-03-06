package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/AkshayS96/bolt/internal/utils"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

func init() {
	jsonCmd := &cobra.Command{
		Use:   "json",
		Short: "JSON tools (format, minify, query, validate, convert)",
	}

	jsonFormatCmd := &cobra.Command{
		Use:   "format [file]",
		Short: "Pretty print JSON",
		Long:  "Pretty print JSON from a file or stdin pipe",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := utils.ReadFileOrStdin(args, 0)
			if err != nil {
				utils.PrintError("No input provided. Use: bolt json format <file> or pipe input")
				return
			}
			var out bytes.Buffer
			if err := json.Indent(&out, []byte(input), "", "  "); err != nil {
				utils.PrintError("Invalid JSON: " + err.Error())
				return
			}
			fmt.Println(out.String())
		},
	}

	jsonMinifyCmd := &cobra.Command{
		Use:   "minify [file]",
		Short: "Minify JSON",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := utils.ReadFileOrStdin(args, 0)
			if err != nil {
				utils.PrintError("No input provided. Use: bolt json minify <file> or pipe input")
				return
			}
			var out bytes.Buffer
			if err := json.Compact(&out, []byte(input)); err != nil {
				utils.PrintError("Invalid JSON: " + err.Error())
				return
			}
			fmt.Println(out.String())
		},
	}

	jsonQueryCmd := &cobra.Command{
		Use:   "query <path> [file]",
		Short: "Query JSON using dot-path syntax",
		Long:  `Query JSON using GJSON path syntax (e.g., "user.name", "users.#.email")`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			query := args[0]
			input, err := utils.ReadFileOrStdin(args, 1)
			if err != nil {
				utils.PrintError("No JSON input provided")
				return
			}
			result := gjson.Get(input, query)
			if !result.Exists() {
				utils.PrintWarning("No match found for path: " + query)
				return
			}
			fmt.Println(result.String())
		},
	}

	jsonValidateCmd := &cobra.Command{
		Use:   "validate [file]",
		Short: "Validate JSON syntax",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := utils.ReadFileOrStdin(args, 0)
			if err != nil {
				utils.PrintError("No input provided")
				return
			}
			if json.Valid([]byte(input)) {
				utils.PrintSuccess("Valid JSON")
			} else {
				utils.PrintError("Invalid JSON")
			}
		},
	}

	jsonToYamlCmd := &cobra.Command{
		Use:   "to-yaml [file]",
		Short: "Convert JSON to YAML",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := utils.ReadFileOrStdin(args, 0)
			if err != nil {
				utils.PrintError("No input provided")
				return
			}
			var data interface{}
			if err := json.Unmarshal([]byte(input), &data); err != nil {
				utils.PrintError("Invalid JSON: " + err.Error())
				return
			}
			yamlBytes, err := yaml.Marshal(data)
			if err != nil {
				utils.PrintError("Failed to convert to YAML: " + err.Error())
				return
			}
			fmt.Print(string(yamlBytes))
		},
	}

	jsonFromYamlCmd := &cobra.Command{
		Use:   "from-yaml [file]",
		Short: "Convert YAML to JSON",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := utils.ReadFileOrStdin(args, 0)
			if err != nil {
				utils.PrintError("No input provided")
				return
			}
			var data interface{}
			if err := yaml.Unmarshal([]byte(input), &data); err != nil {
				utils.PrintError("Invalid YAML: " + err.Error())
				return
			}
			data = convertYAMLToJSON(data)
			jsonBytes, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				utils.PrintError("Failed to convert to JSON: " + err.Error())
				return
			}
			fmt.Println(string(jsonBytes))
		},
	}

	jsonCmd.AddCommand(jsonFormatCmd, jsonMinifyCmd, jsonQueryCmd, jsonValidateCmd, jsonToYamlCmd, jsonFromYamlCmd)
	rootCmd.AddCommand(jsonCmd)
}

// convertYAMLToJSON recursively converts map[interface{}]interface{} to map[string]interface{}
// because YAML unmarshals maps with interface{} keys, but JSON needs string keys.
func convertYAMLToJSON(v interface{}) interface{} {
	switch val := v.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			m[fmt.Sprintf("%v", k)] = convertYAMLToJSON(v)
		}
		return m
	case map[string]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			m[k] = convertYAMLToJSON(v)
		}
		return m
	case []interface{}:
		for i, v := range val {
			val[i] = convertYAMLToJSON(v)
		}
		return val
	default:
		return v
	}
}
