package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
	"vaultpull/internal/secret"
)

var (
	snapshotFile  string
	snapshotLabel string
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Save or load a snapshot of the current .env secrets",
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save current .env to a snapshot file",
	RunE: func(cmd *cobra.Command, args []string) error {
		secrets, err := env.ParseFile(envFile)
		if err != nil {
			return fmt.Errorf("parse env: %w", err)
		}
		if err := secret.SaveSnapshot(snapshotFile, snapshotLabel, secrets); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "snapshot saved to %s\n", snapshotFile)
		return nil
	},
}

var snapshotShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show summary of an existing snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := secret.LoadSnapshot(snapshotFile)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), secret.FormatSnapshotSummary(s))
		return nil
	},
}

func init() {
	snapshotSaveCmd.Flags().StringVar(&envFile, "env-file", ".env", "path to .env file")
	snapshotSaveCmd.Flags().StringVar(&snapshotFile, "out", "snapshot.json", "path to snapshot output file")
	snapshotSaveCmd.Flags().StringVar(&snapshotLabel, "label", "", "optional label for the snapshot")

	snapshotShowCmd.Flags().StringVar(&snapshotFile, "file", "snapshot.json", "path to snapshot file")

	snapshotCmd.AddCommand(snapshotSaveCmd)
	snapshotCmd.AddCommand(snapshotShowCmd)

	if err := rootCmd.GenBashCompletion(os.Discard); err == nil {
		rootCmd.AddCommand(snapshotCmd)
	}
	rootCmd.AddCommand(snapshotCmd)
}
