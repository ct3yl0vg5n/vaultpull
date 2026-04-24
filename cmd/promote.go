package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/env"
	"github.com/yourorg/vaultpull/internal/secret"
)

var (
	promoteFrom      string
	promoteTo        string
	promoteOverwrite bool
	promoteDryRun    bool
	promoteIgnore    []string
)

var promoteCmd = &cobra.Command{
	Use:   "promote",
	Short: "Promote secrets from one environment file to another",
	RunE:  runPromote,
}

func init() {
	promoteCmd.Flags().StringVar(&promoteFrom, "from", ".env.staging", "source .env file")
	promoteCmd.Flags().StringVar(&promoteTo, "to", ".env.production", "destination .env file")
	promoteCmd.Flags().BoolVar(&promoteOverwrite, "overwrite", false, "overwrite existing keys in destination")
	promoteCmd.Flags().BoolVar(&promoteDryRun, "dry-run", false, "preview changes without writing")
	promoteCmd.Flags().StringSliceVar(&promoteIgnore, "ignore", nil, "comma-separated keys to skip")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, _ []string) error {
	src, err := env.ParseFile(promoteFrom)
	if err != nil {
		return fmt.Errorf("reading source file %q: %w", promoteFrom, err)
	}

	dst, err := env.ParseFile(promoteTo)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading destination file %q: %w", promoteTo, err)
	}
	if dst == nil {
		dst = map[string]string{}
	}

	opts := secret.DefaultPromoteOptions()
	opts.FromEnv = promoteFrom
	opts.ToEnv = promoteTo
	opts.DryRun = promoteDryRun
	opts.OverwriteExisting = promoteOverwrite
	opts.IgnoreKeys = promoteIgnore

	merged, results := secret.Promote(src, dst, opts)

	report := secret.FormatPromoteReport(promoteFrom, promoteTo, results)
	fmt.Fprint(cmd.OutOrStdout(), report)

	if promoteDryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "Dry run — no files written.")
		return nil
	}

	if err := env.WriteFile(promoteTo, merged); err != nil {
		return fmt.Errorf("writing destination file %q: %w", promoteTo, err)
	}
	return nil
}
