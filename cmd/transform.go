package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/secret"
)

var (
	transformFile    string
	transformPrefix  string
	transformReplace string
	transformUpper   bool
	transformRedact  bool
)

var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "Apply key/value transformations to a local .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		secrets, err := env.ParseFile(transformFile)
		if err != nil {
			return fmt.Errorf("parse file: %w", err)
		}

		rule := secret.Rule{
			UppercaseKeys: transformUpper,
			KeyPrefix:     transformPrefix,
			RedactValues:  transformRedact,
		}

		if transformReplace != "" {
			var from, to string
			_, err := fmt.Sscanf(transformReplace, "%s %s", &from, &to)
			if err != nil {
				return fmt.Errorf("--replace must be 'FROM TO': %w", err)
			}
			rule.KeyReplace = &secret.KeyReplace{From: from, To: to}
		}

		transformed, err := secret.ApplyMap(secrets, rule)
		if err != nil {
			return fmt.Errorf("transform: %w", err)
		}

		if err := env.WriteFile(transformFile, transformed); err != nil {
			return fmt.Errorf("write file: %w", err)
		}

		fmt.Fprintf(os.Stdout, "transformed %d keys in %s\n", len(transformed), transformFile)
		return nil
	},
}

func init() {
	transformCmd.Flags().StringVarP(&transformFile, "file", "f", ".env", "path to .env file")
	transformCmd.Flags().StringVar(&transformPrefix, "prefix", "", "prefix to prepend to all keys")
	transformCmd.Flags().StringVar(&transformReplace, "replace", "", "find and replace in key names, e.g. '- _'")
	transformCmd.Flags().BoolVar(&transformUpper, "uppercase", false, "uppercase all keys")
	transformCmd.Flags().BoolVar(&transformRedact, "redact", false, "redact all values (for display/audit use)")
	rootCmd.AddCommand(transformCmd)
}
