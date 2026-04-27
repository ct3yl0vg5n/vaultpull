package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/internal/env"
	"github.com/vaultpull/internal/secret"
)

var envExportCmd = &cobra.Command{
	Use:   "env-export",
	Short: "Export secrets from an env file in shell-compatible format",
	RunE:  runEnvExport,
}

func init() {
	envExportCmd.Flags().String("file", ".env", "Path to the env file")
	envExportCmd.Flags().String("prefix", "", "Key prefix to prepend")
	envExportCmd.Flags().Bool("export", false, "Add 'export' declaration to each line")
	envExportCmd.Flags().Bool("quote", false, "Quote all values")
	envExportCmd.Flags().Bool("skip-empty", false, "Skip keys with empty values")
	rootCmd.AddCommand(envExportCmd)
}

func runEnvExport(cmd *cobra.Command, _ []string) error {
	filePath, _ := cmd.Flags().GetString("file")
	prefix, _ := cmd.Flags().GetString("prefix")
	exportDecl, _ := cmd.Flags().GetBool("export")
	quote, _ := cmd.Flags().GetBool("quote")
	skipEmpty, _ := cmd.Flags().GetBool("skip-empty")

	secrets, err := env.ParseFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("env file not found: %s", filePath)
		}
		return fmt.Errorf("failed to parse env file: %w", err)
	}

	opts := secret.DefaultExportOptions()
	opts.Prefix = prefix
	opts.ExportDecl = exportDecl
	opts.QuoteValues = quote
	opts.SkipEmpty = skipEmpty

	results := secret.Export(secrets, opts)
	fmt.Println(secret.FormatExport(results))
	return nil
}
