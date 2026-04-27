package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/secret"
)

var dedupeCmd = &cobra.Command{
	Use:   "dedupe",
	Short: "Remove duplicate keys from a .env file",
	RunE:  runDedupe,
}

var (
	dedupeCaseSensitive bool
	dedupePreferLast    bool
	dedupeReportOnly    bool
	dedupeFile          string
)

func init() {
	dedupeCmd.Flags().StringVarP(&dedupeFile, "file", "f", ".env", "path to .env file")
	dedupeCmd.Flags().BoolVar(&dedupeCaseSensitive, "case-sensitive", true, "treat keys as case-sensitive")
	dedupeCmd.Flags().BoolVar(&dedupePreferLast, "prefer-last", false, "keep last occurrence of duplicate keys")
	dedupeCmd.Flags().BoolVar(&dedupeReportOnly, "report-only", false, "report duplicates without modifying the file")
	rootCmd.AddCommand(dedupeCmd)
}

func runDedupe(cmd *cobra.Command, _ []string) error {
	parsed, err := env.ParseFile(dedupeFile)
	if err != nil {
		return fmt.Errorf("parse %s: %w", dedupeFile, err)
	}

	opts := secret.DefaultDedupeOptions()
	opts.CaseSensitive = dedupeCaseSensitive
	opts.PreferLast = dedupePreferLast
	opts.ReportOnly = dedupeReportOnly

	result := secret.Dedupe(parsed, nil, opts)

	fmt.Fprintln(cmd.OutOrStdout(), secret.FormatDedupeReport(result))

	if dedupeReportOnly || len(result.Duplicates) == 0 {
		return nil
	}

	if err := env.WriteFile(dedupeFile, result.Out); err != nil {
		return fmt.Errorf("write %s: %w", dedupeFile, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "wrote deduplicated secrets to %s\n", dedupeFile)
	return nil
}

// exitIfErr is a small helper used across cmd files.
func exitIfErrDedupe(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
