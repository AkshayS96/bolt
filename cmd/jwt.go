package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bolt/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

func init() {
	jwtCmd := &cobra.Command{
		Use:   "jwt",
		Short: "JWT tools (decode, inspect, verify)",
	}

	jwtDecodeCmd := &cobra.Command{
		Use:   "decode <token>",
		Short: "Decode and display full JWT (header + payload)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			token := args[0]
			header, payload, err := parseJWTParts(token)
			if err != nil {
				utils.PrintError(err.Error())
				return
			}
			utils.Key.Println("── Header ──")
			printJSON(header)
			fmt.Println()
			utils.Key.Println("── Payload ──")
			printJSON(payload)
		},
	}

	jwtHeaderCmd := &cobra.Command{
		Use:   "header <token>",
		Short: "Show JWT header",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			header, _, err := parseJWTParts(args[0])
			if err != nil {
				utils.PrintError(err.Error())
				return
			}
			printJSON(header)
		},
	}

	jwtPayloadCmd := &cobra.Command{
		Use:   "payload <token>",
		Short: "Show JWT payload",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, payload, err := parseJWTParts(args[0])
			if err != nil {
				utils.PrintError(err.Error())
				return
			}
			printJSON(payload)
		},
	}

	jwtExpCmd := &cobra.Command{
		Use:   "exp <token>",
		Short: "Show JWT expiration time",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, payload, err := parseJWTParts(args[0])
			if err != nil {
				utils.PrintError(err.Error())
				return
			}
			expVal, ok := payload["exp"]
			if !ok {
				utils.PrintWarning("No 'exp' claim found in token")
				return
			}
			expFloat, ok := expVal.(float64)
			if !ok {
				utils.PrintError("Invalid 'exp' claim format")
				return
			}
			expTime := time.Unix(int64(expFloat), 0)
			now := time.Now()
			utils.PrintKeyValue("Expires:", expTime.Format(time.RFC3339))
			if now.After(expTime) {
				utils.PrintError(fmt.Sprintf("Token EXPIRED %s ago", now.Sub(expTime).Round(time.Second)))
			} else {
				utils.PrintSuccess(fmt.Sprintf("Token valid for %s", expTime.Sub(now).Round(time.Second)))
			}
		},
	}

	var secret string
	jwtVerifyCmd := &cobra.Command{
		Use:   "verify <token>",
		Short: "Verify JWT signature",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if secret == "" {
				utils.PrintError("--secret flag is required")
				return
			}
			token, err := jwt.Parse(args[0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})
			if err != nil {
				utils.PrintError("Verification failed: " + err.Error())
				return
			}
			if token.Valid {
				utils.PrintSuccess("Signature is valid ✓")
			} else {
				utils.PrintError("Signature is invalid ✗")
			}
		},
	}
	jwtVerifyCmd.Flags().StringVar(&secret, "secret", "", "HMAC secret key")

	jwtCmd.AddCommand(jwtDecodeCmd, jwtHeaderCmd, jwtPayloadCmd, jwtExpCmd, jwtVerifyCmd)
	rootCmd.AddCommand(jwtCmd)
}

func parseJWTParts(token string) (header map[string]interface{}, payload map[string]interface{}, err error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, nil, fmt.Errorf("invalid JWT: expected at least 2 parts, got %d", len(parts))
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode header: %w", err)
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, nil, fmt.Errorf("failed to parse header JSON: %w", err)
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode payload: %w", err)
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, nil, fmt.Errorf("failed to parse payload JSON: %w", err)
	}

	return header, payload, nil
}

func printJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(data)
		return
	}
	fmt.Println(string(b))
}
