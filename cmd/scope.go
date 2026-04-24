package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vaultpull/internal/env"
	"github.com/vaultpull/internal/secret"
)

var (
	scopeFile        string
	scopeDefinitions []string
	scopeDefault     string
)

var scopeCmd = &cobra.Command{
	Use:   "scope",
	Short: "Partition secrets from a .env file into named scopes by key prefix",
	RunE:  runScope,
}

func init() {
	scopeCmd.Flags().StringVarP(&scopeFile, "file", "f", ".env", "source .env file")
	scopeCmd.Flags().StringArrayVarP(&scopeDefinitions, "scope", "s", nil,
		"scope definition as name=PREFIX1,PREFIX2 (repeatable)")
	scopeCmd.Flags().StringVar(&scopeDefault, "default", "default", "name for the default scope")
	rootCmd.AddCommand(scopeCmd)
}

func runScope(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ParseFile(scopeFile)
	if err != nil {
		return fmt.Errorf("parse %s: %w", scopeFile, err)
	}

	opts := secret.ScopeOptions{
		Scopes:       map[string][]string{},
		DefaultScope: scopeDefault,
	}

	for _, def := range scopeDefinitions {
		parts := strings.SplitN(def, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid scope definition %q: expected name=PREFIX,...", def)
		}
		name := strings.TrimSpace(parts[0])
		prefixes := strings.Split(parts[1], ",")
		for i, p := range prefixes {
			prefixes[i] = strings.TrimSpace(p)
		}
		opts.Scopes[name] = prefixes
	}

	results := secret.Partition(secrets, opts)
	fmt.Fprint(os.Stdout, secret.FormatScopeReport(results))
	return nil
}
