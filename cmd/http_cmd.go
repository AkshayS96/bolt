package cmd

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"bytes"
	"encoding/json"

	"github.com/AkshayS96/bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	httpCmd := &cobra.Command{
		Use:   "http",
		Short: "HTTP/API request tools",
	}

	httpGetCmd := &cobra.Command{
		Use:   "get <url>",
		Short: "Send a GET request",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			doHTTP("GET", args[0], "")
		},
	}

	httpPostCmd := &cobra.Command{
		Use:   "post <url>",
		Short: "Send a POST request (reads body from stdin)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			body, _ := utils.ReadStdin()
			doHTTP("POST", args[0], body)
		},
	}

	httpHeadersCmd := &cobra.Command{
		Use:   "headers <url>",
		Short: "Show response headers",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Head(args[0])
			if err != nil {
				utils.PrintError("Request failed: " + err.Error())
				return
			}
			defer resp.Body.Close()

			utils.Key.Printf("HTTP %s\n", resp.Status)
			fmt.Println()
			for key, values := range resp.Header {
				utils.PrintKeyValue(key+":", strings.Join(values, ", "))
			}
		},
	}

	httpJSONCmd := &cobra.Command{
		Use:   "json <url>",
		Short: "GET request with pretty-printed JSON response",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Get(args[0])
			if err != nil {
				utils.PrintError("Request failed: " + err.Error())
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				utils.PrintError("Failed to read response: " + err.Error())
				return
			}

			var out bytes.Buffer
			if err := json.Indent(&out, body, "", "  "); err != nil {
				// Not JSON, print raw
				fmt.Println(string(body))
				return
			}
			fmt.Println(out.String())
		},
	}

	httpCmd.AddCommand(httpGetCmd, httpPostCmd, httpHeadersCmd, httpJSONCmd)
	rootCmd.AddCommand(httpCmd)
}

func doHTTP(method, url, body string) {
	client := &http.Client{Timeout: 10 * time.Second}

	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		utils.PrintError("Invalid request: " + err.Error())
		return
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		utils.PrintError("Request failed: " + err.Error())
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.PrintError("Failed to read response: " + err.Error())
		return
	}

	utils.Dim.Printf("HTTP %s %s (%s)\n", method, resp.Status, elapsed.Round(time.Millisecond))
	fmt.Println()
	fmt.Println(string(respBody))
}
