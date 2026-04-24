package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
	"vaultpull/internal/secret"
)

var (
	redactReveal  int
	redactMask    string
	redactEnvFile string
)

var redactCmd = &cobra.Command{
	Use:   "redact",
	Short: "Print env file with secret values redacted",
	RunE: func(cmd *cobra.Command, args []string) error {
		secrets, err := env.ParseFile(redactEnvFile)
		if err != nil {
			return fmt.Errorf("parse env file: %w", err)
		}

		opts := secret.RedactOptions{
			Enabled:     true,
			MaskChar:    redactMask,
			RevealChars: redactReveal,
		}

		redacted := secret.RedactMap(secrets, opts)
		printRedacted(redacted)
		return nil
	},
}

// printRedacted writes redacted key=value pairs to stdout in sorted order
// so that output is deterministic regardless of map iteration order.
func printRedacted(redacted map[string]string) {
	keys := make([]string, 0, len(redacted))
	for k := range redacted {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, redacted[k])
	}
}

func init() {
	redactCmd.Flags().StringVar(&redactEnvFile, "file", ".env", "path to .env file")
	redactCmd.Flags().IntVar(&redactReveal, "reveal", 4, "number of trailing characters to reveal")
	redactCmd.Flags().StringVar(&redactMask, "mask", "*", "mask character")
	rootCmd.AddCommand(redactCmd)
}
