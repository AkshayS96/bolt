package cmd

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea", "commodo",
	"consequat", "duis", "aute", "irure", "in", "reprehenderit", "voluptate",
	"velit", "esse", "cillum", "fugiat", "nulla", "pariatur", "excepteur", "sint",
	"occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia",
	"deserunt", "mollit", "anim", "id", "est", "laborum", "at", "vero", "eos",
	"accusamus", "iusto", "odio", "dignissimos", "ducimus", "blanditiis",
	"praesentium", "voluptatum", "deleniti", "atque", "corrupti", "quos",
	"dolores", "quas", "molestias", "recusandae", "itaque", "earum", "rerum",
	"hic", "tenetur", "sapiente", "delectus", "aut", "reiciendis", "voluptatibus",
	"maiores", "alias", "perferendis", "doloribus", "asperiores", "repellat",
}

func init() {
	loremCmd := &cobra.Command{
		Use:   "lorem [paragraphs]",
		Short: "Generate lorem ipsum text",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			count := 1
			if len(args) > 0 {
				n, err := strconv.Atoi(args[0])
				if err != nil || n <= 0 {
					count = 1
				} else {
					count = n
				}
			}

			wordsPerSentence, _ := cmd.Flags().GetBool("words")
			if wordsPerSentence {
				n := count
				if n > 500 {
					n = 500
				}
				words := make([]string, n)
				for i := 0; i < n; i++ {
					words[i] = loremWords[rand.Intn(len(loremWords))]
				}
				fmt.Println(strings.Join(words, " "))
				return
			}

			for i := 0; i < count; i++ {
				fmt.Println(generateParagraph())
				if i < count-1 {
					fmt.Println()
				}
			}
		},
	}
	loremCmd.Flags().BoolP("words", "w", false, "Generate N words instead of paragraphs")

	rootCmd.AddCommand(loremCmd)
}

func generateParagraph() string {
	sentences := 4 + rand.Intn(4) // 4-7 sentences
	var parts []string
	for i := 0; i < sentences; i++ {
		parts = append(parts, generateSentence())
	}
	return strings.Join(parts, " ")
}

func generateSentence() string {
	wordCount := 6 + rand.Intn(10) // 6-15 words
	words := make([]string, wordCount)
	for i := 0; i < wordCount; i++ {
		words[i] = loremWords[rand.Intn(len(loremWords))]
	}
	words[0] = strings.ToUpper(words[0][:1]) + words[0][1:]
	return strings.Join(words, " ") + "."
}
