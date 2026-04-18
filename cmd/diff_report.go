package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
	"vaultpull/internal/secret"
)

var (
	diffReportFile    string
	diffReportCompare string
	diffReportRedact  bool
)

var diffReportCmd = &cobra.Command{
	Use:   "diff-report",
	Short: "Show a diff report between two .env files",
	RunE:  runDiffReport,
}

func init() {
	diffReportCmd.Flags().StringVar(&diffReportFile, "file", ".env", "base .env file")
	diffReportCmd.Flags().StringVar(&diffReportCompare, "compare", "", "file to compare against (required)")
	diffReportCmd.Flags().BoolVar(&diffReportRedact, "redact", true, "redact secret values in output")
	_ = diffReportCmd.MarkFlagRequired("compare")
	rootCmd.AddCommand(diffReportCmd)
}

func runDiffReport(cmd *cobra.Command, _ []string) error {
	base, err := env.ParseFile(diffReportFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading base file: %w", err)
	}

	next, err := env.ParseFile(diffReportCompare)
	if err != nil {
		return fmt.Errorf("reading compare file: %w", err)
	}

	opts := secret.DefaultDiffReportOptions()
	opts.RedactValues = diffReportRedact

	entries := secret.BuildDiffReport(base, next, opts)
	fmt.Fprint(cmd.OutOrStdout(), secret.FormatDiffReport(entries))
	return nil
}
