package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/secret"
)

var truncateCmd = &cobra.Command{
	Use:   "truncate",
	Short: "Truncate long secret values in a .env file",
	RunE:  runTruncate,
}

func init() {
	truncateCmd.Flags().String("file", ".env", "Path to the .env file")
	truncateCmd.Flags().Int("max-length", 64, "Maximum value length before truncation")
	truncateCmd.Flags().String("suffix", "...", "Suffix appended to truncated values")
	truncateCmd.Flags().StringSlice("skip", nil, "Keys to skip during truncation")
	truncateCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	truncateCmd.Flags().Bool("write", false, "Write truncated values back to file")
	rootCmd.AddCommand(truncateCmd)
}

func runTruncate(cmd *cobra.Command, _ []string) error {
	filePath, _ := cmd.Flags().GetString("file")
	maxLen, _ := cmd.Flags().GetInt("max-length")
	suffix, _ := cmd.Flags().GetString("suffix")
	skipKeys, _ := cmd.Flags().GetStringSlice("skip")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	write, _ := cmd.Flags().GetBool("write")

	secrets, err := env.ParseFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(cmd.OutOrStdout(), "file not found: %s\n", filePath)
			return nil
		}
		return fmt.Errorf("parse file: %w", err)
	}

	opts := secret.DefaultTruncateOptions()
	opts.MaxLength = maxLen
	opts.Suffix = suffix
	opts.SkipKeys = skipKeys

	truncated, results := secret.TruncateMap(secrets, opts)
	report := secret.FormatTruncateReport(results)
	fmt.Fprint(cmd.OutOrStdout(), report)

	if dryRun || !write {
		return nil
	}

	if err := env.WriteFile(filePath, truncated); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Written to %s\n", filePath)
	return nil
}
