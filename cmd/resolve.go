package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/secret"
)

var (
	resolveFile      string
	resolveRefPrefix string
	resolveStrict    bool
)

var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve secret references (ref:<KEY>) within an env file",
	RunE:  runResolve,
}

func init() {
	resolveCmd.Flags().StringVarP(&resolveFile, "file", "f", ".env", "env file to resolve")
	resolveCmd.Flags().StringVar(&resolveRefPrefix, "prefix", "ref:", "prefix that marks a value as a reference")
	resolveCmd.Flags().BoolVar(&resolveStrict, "strict", true, "fail if a referenced key is missing")
	rootCmd.AddCommand(resolveCmd)
}

func runResolve(cmd *cobra.Command, _ []string) error {
	src, err := env.ParseFile(resolveFile)
	if err != nil {
		return fmt.Errorf("resolve: parse %q: %w", resolveFile, err)
	}

	opts := secret.DefaultResolveOptions()
	opts.RefPrefix = resolveRefPrefix
	opts.Strict = resolveStrict

	out, results, err := secret.Resolve(src, opts)
	if err != nil {
		return err
	}

	fmt.Fprint(cmd.OutOrStdout(), secret.FormatResolveReport(results))

	if err := env.WriteFile(resolveFile, out); err != nil {
		return fmt.Errorf("resolve: write %q: %w", resolveFile, err)
	}

	fmt.Fprintf(os.Stdout, "Written to %s\n", resolveFile)
	return nil
}
