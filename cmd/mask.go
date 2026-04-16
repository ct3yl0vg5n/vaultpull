package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
	"vaultpull/internal/secret"
)

var (
	maskRevealLast int
	maskChar       string
	maskFile       string
)

var maskCmd = &cobra.Command{
	Use:   "mask",
	Short: "Print .env file with secret values masked",
	RunE: func(cmd *cobra.Command, args []string) error {
		secrets, err := env.ParseFile(maskFile)
		if err != nil {
			return fmt.Errorf("parse file: %w", err)
		}

		opts := secret.MaskOptions{
			Enabled:    true,
			MaskChar:   maskChar,
			RevealLast: maskRevealLast,
		}

		masked := secret.MaskMap(secrets, opts)

		for k, v := range masked {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	},
}

func init() {
	maskCmd.Flags().StringVarP(&maskFile, "file", "f", ".env", "path to .env file")
	maskCmd.Flags().IntVar(&maskRevealLast, "reveal-last", 4, "number of characters to reveal at end of value")
	maskCmd.Flags().StringVar(&maskChar, "mask-char", "*", "character used for masking")
	rootCmd.AddCommand(maskCmd)
}
