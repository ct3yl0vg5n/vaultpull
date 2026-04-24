package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultpull/internal/env"
	"github.com/example/vaultpull/internal/secret"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Restore secrets in an env file to values from a snapshot",
	RunE:  runRollback,
}

var (
	rollbackFile     string
	rollbackSnapshot string
	rollbackDryRun   bool
	rollbackIgnore   []string
)

func init() {
	rollbackCmd.Flags().StringVar(&rollbackFile, "file", ".env", "env file to roll back")
	rollbackCmd.Flags().StringVar(&rollbackSnapshot, "snapshot", "", "snapshot file to restore from (required)")
	rollbackCmd.Flags().BoolVar(&rollbackDryRun, "dry-run", false, "preview changes without writing")
	rollbackCmd.Flags().StringSliceVar(&rollbackIgnore, "ignore", nil, "keys to skip during rollback")
	_ = rollbackCmd.MarkFlagRequired("snapshot")
	rootCmd.AddCommand(rollbackCmd)
}

func runRollback(cmd *cobra.Command, _ []string) error {
	current, err := env.ParseFile(rollbackFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading env file: %w", err)
	}
	if current == nil {
		current = map[string]string{}
	}

	snap, err := secret.LoadSnapshot(rollbackSnapshot)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	opts := secret.DefaultRollbackOptions()
	opts.DryRun = rollbackDryRun
	opts.IgnoreKeys = rollbackIgnore

	updated, result := secret.Rollback(current, snap, opts)

	fmt.Fprint(cmd.OutOrStdout(), secret.FormatRollbackReport(result))

	if !rollbackDryRun && len(result.Applied) > 0 {
		if err := env.WriteFile(rollbackFile, updated); err != nil {
			return fmt.Errorf("writing env file: %w", err)
		}
	}
	return nil
}
