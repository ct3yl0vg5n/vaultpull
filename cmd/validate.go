package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
	"vaultpull/internal/secret"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate keys and values in a local .env file",
	RunE:  runValidate,
}

func init() {
	validateCmd.Flags().String("file", ".env", "path to .env file")
	validateCmd.Flags().Int("min-length", 1, "minimum value length")
	validateCmd.Flags().Int("max-length", 4096, "maximum value length")
	validateCmd.Flags().Bool("strict-keys", true, "enforce UPPER_SNAKE_CASE key pattern")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, _ []string) error {
	filePath, _ := cmd.Flags().GetString("file")
	minLen, _ := cmd.Flags().GetInt("min-length")
	maxLen, _ := cmd.Flags().GetInt("max-length")
	strictKeys, _ := cmd.Flags().GetBool("strict-keys")

	secrets, err := env.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parsing file: %w", err)
	}

	opts := secret.DefaultValidateOptions()
	opts.MinLength = minLen
	opts.MaxLength = maxLen
	if !strictKeys {
		opts.KeyPattern = ""
	}

	errs := secret.ValidateMap(secrets, opts)
	if errs == nil {
		fmt.Fprintln(cmd.OutOrStdout(), "all secrets valid")
		return nil
	}

	for k, e := range errs {
		fmt.Fprintf(os.Stderr, "  [invalid] %s: %v\n", k, e)
	}
	return fmt.Errorf("%d validation error(s) found", len(errs))
}
