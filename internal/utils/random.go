package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	AlphaNum   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AlphaLower = "abcdefghijklmnopqrstuvwxyz"
	AlphaUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits     = "0123456789"
	Symbols    = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	AllChars   = AlphaNum + Symbols
)

// RandomString generates a cryptographically secure random string of the given length using the provided charset.
func RandomString(length int, charset string) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))
	for i := range result {
		idx, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[idx.Int64()]
	}
	return string(result), nil
}
