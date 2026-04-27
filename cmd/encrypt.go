package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/secret"
)

var (
	encryptFile    string
	encryptKeyHex  string
	encryptDecrypt bool
	encryptOutput  string
)

func init() {
	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt or decrypt values in a .env file using AES-256-GCM",
		RunE:  runEncrypt,
	}

	encryptCmd.Flags().StringVarP(&encryptFile, "file", "f", ".env", "path to .env file")
	encryptCmd.Flags().StringVar(&encryptKeyHex, "key", "", "32-byte AES key (required)")
	encryptCmd.Flags().BoolVar(&encryptDecrypt, "decrypt", false, "decrypt values instead of encrypting")
	encryptCmd.Flags().StringVarP(&encryptOutput, "output", "o", "", "output file path (defaults to stdout)")
	_ = encryptCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(encryptCmd)
}

func runEncrypt(cmd *cobra.Command, _ []string) error {
	if len(encryptKeyHex) != 32 {
		return fmt.Errorf("--key must be exactly 32 ASCII characters (got %d)", len(encryptKeyHex))
	}

	secrets, err := env.ParseFile(encryptFile)
	if err != nil {
		return fmt.Errorf("parse %s: %w", encryptFile, err)
	}

	opts := secret.DefaultEncryptOptions()
	opts.Key = []byte(encryptKeyHex)

	var result map[string]string
	if encryptDecrypt {
		result, err = secret.DecryptMap(secrets, opts)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
	} else {
		result, err = secret.EncryptMap(secrets, opts)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
	}

	if encryptOutput == "" {
		for k, v := range result {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
		}
		return nil
	}

	if err := env.WriteFile(encryptOutput, result); err != nil {
		return fmt.Errorf("write %s: %w", encryptOutput, err)
	}

	action := "Encrypted"
	if encryptDecrypt {
		action = "Decrypted"
	}
	fmt.Fprintf(os.Stderr, "%s %d key(s) → %s\n", action, len(result), encryptOutput)
	return nil
}
