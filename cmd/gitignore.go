package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/AkshayS96/bolt/internal/utils"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

func init() {
	gitignoreCmd := &cobra.Command{
		Use:   "gitignore <language>",
		Short: "Generate a .gitignore file from GitHub templates",
		Long: `Generate a .gitignore from GitHub's official templates.

  Examples:
    bolt gitignore go
    bolt gitignore node
    bolt gitignore python
    bolt gitignore java
    bolt gitignore list          # Show available templates`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			lang := strings.ToLower(args[0])

			if lang == "list" || lang == "ls" {
				showGitignoreList()
				return
			}

			save, _ := cmd.Flags().GetBool("save")
			fetchGitignore(lang, save)
		},
	}

	gitignoreCmd.Flags().BoolP("save", "s", false, "Save as .gitignore file in current directory")
	rootCmd.AddCommand(gitignoreCmd)
}

func fetchGitignore(lang string, save bool) {
	// GitHub's gitignore API uses title case
	name := strings.Title(lang)

	// Common mappings
	aliases := map[string]string{
		"js":         "Node",
		"javascript": "Node",
		"node":       "Node",
		"ts":         "Node",
		"typescript": "Node",
		"go":         "Go",
		"golang":     "Go",
		"py":         "Python",
		"python":     "Python",
		"rb":         "Ruby",
		"ruby":       "Ruby",
		"rs":         "Rust",
		"rust":       "Rust",
		"java":       "Java",
		"kotlin":     "Java",
		"swift":      "Swift",
		"c":          "C",
		"cpp":        "C++",
		"c++":        "C++",
		"csharp":     "VisualStudio",
		"cs":         "VisualStudio",
		"dotnet":     "VisualStudio",
		"dart":       "Dart",
		"flutter":    "Dart",
		"unity":      "Unity",
		"react":      "Node",
		"vue":        "Node",
		"nextjs":     "Node",
		"rails":      "Rails",
		"django":     "Python",
		"laravel":    "Laravel",
		"android":    "Android",
		"ios":        "Swift",
		"macos":      "macOS",
		"linux":      "Linux",
		"windows":    "Windows",
		"terraform":  "Terraform",
		"elixir":     "Elixir",
		"haskell":    "Haskell",
		"scala":      "Scala",
		"r":          "R",
		"tex":        "TeX",
		"latex":      "TeX",
	}

	if mapped, ok := aliases[lang]; ok {
		name = mapped
	}

	url := fmt.Sprintf("https://raw.githubusercontent.com/github/gitignore/main/%s.gitignore", name)
	resp, err := http.Get(url)
	if err != nil {
		utils.PrintError("Failed to fetch: " + err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		utils.PrintError(fmt.Sprintf("No gitignore template found for '%s'", lang))
		dimC := color.New(color.Faint)
		dimC.Println("  Run: bolt gitignore list")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.PrintError("Failed to read response: " + err.Error())
		return
	}

	content := string(body)

	if save {
		if err := os.WriteFile(".gitignore", body, 0644); err != nil {
			utils.PrintError("Failed to write .gitignore: " + err.Error())
			return
		}
		utils.PrintSuccess(fmt.Sprintf("Saved .gitignore for %s (%d bytes)", name, len(body)))
	} else {
		accent := color.New(color.FgCyan, color.Bold)
		accent.Printf("# .gitignore for %s\n", name)
		fmt.Println(content)
	}
}

func showGitignoreList() {
	accent := color.New(color.FgCyan, color.Bold)
	dimC := color.New(color.Faint)

	accent.Println("  Available gitignore templates")
	fmt.Println()

	categories := []struct {
		name  string
		langs []string
	}{
		{"Languages", []string{"go", "python", "node/js", "rust", "java", "swift", "c", "c++", "ruby", "dart", "elixir", "haskell", "scala", "r", "tex"}},
		{"Frameworks", []string{"react/vue/nextjs → node", "rails", "django → python", "laravel", "flutter → dart"}},
		{"Platforms", []string{"android", "ios → swift", "unity", "macos", "linux", "windows"}},
		{"Tools", []string{"terraform"}},
	}

	for _, cat := range categories {
		accent.Printf("  %s\n", cat.name)
		for _, l := range cat.langs {
			dimC.Printf("    • %s\n", l)
		}
		fmt.Println()
	}

	dimC.Println("  Usage: bolt gitignore <name>")
	dimC.Println("  Save:  bolt gitignore <name> --save")
}
