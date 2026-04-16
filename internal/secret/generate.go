package secret

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

const (
	Alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AlphaNumericSpecial = AlphanumericSpecial
	AlphanumericSpecial = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
)

// GenerateOptions controls secret generation.
type GenerateOptions struct {
	Length  int
	Charset string
	Base64  bool
}

// DefaultGenerateOptions returns sensible defaults.
func DefaultGenerateOptions() GenerateOptions {
	return GenerateOptions{
		Length:  32,
		Charset: Alphanumeric,
		Base64:  false,
	}
}

// Generate creates a random secret string.
func Generate(opts GenerateOptions) (string, error) {
	if opts.Length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}
	if opts.Base64 {
		buf := make([]byte, opts.Length)
		if _, err := rand.Read(buf); err != nil {
			return "", fmt.Errorf("rand read: %w", err)
		}
		return base64.StdEncoding.EncodeToString(buf)[:opts.Length], nil
	}
	charset := opts.Charset
	if charset == "" {
		charset = Alphanumeric
	}
	out := make([]byte, opts.Length)
	for i := range out {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("rand int: %w", err)
		}
		out[i] = charset[n.Int64()]
	}
	return string(out), nil
}
