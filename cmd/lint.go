package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/internal/env"
	"github.com/vaultpull/internal/secret"
)

var lintFile string

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint secret keys and values in a .env file",
	RunE:  runLint,
}

func init() {
	lintCmd.Flags().StringVarP(&lintFile, "file", "f", ".env", "Path to .env file")
	RootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	secrets, err := env.ParseFile(lintFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	opts := secret.DefaultLintOptions()
	results := secret.LintMap(secrets, opts)

	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✔ No lint issues found.")
		return nil
	}

	for _, r := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "⚠  [%s] %s\n", r.Key, r.Message)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "\n%d issue(s) found.\n", len(results))
	os.Exit(1)
	return nil
}
