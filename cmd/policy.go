package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/env"
	"github.com/yourorg/vaultpull/internal/secret"
)

var (
	policyFile        string
	policyRequireKeys []string
	policyDenyValues  []string
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Enforce secret policy rules against a .env file",
	RunE:  runPolicy,
}

func init() {
	policyCmd.Flags().StringVar(&policyFile, "file", ".env", "path to .env file")
	policyCmd.Flags().StringArrayVar(&policyRequireKeys, "require", nil, "keys that must be present")
	policyCmd.Flags().StringArrayVar(&policyDenyValues, "deny-value", nil, "regex patterns values must not match (applied to all keys)")
	rootCmd.AddCommand(policyCmd)
}

func runPolicy(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ParseFile(policyFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	opts := secret.DefaultPolicyOptions()

	for _, k := range policyRequireKeys {
		opts.Rules = append(opts.Rules, secret.PolicyRule{
			Name:     k,
			Required: true,
		})
	}

	for i, pattern := range policyDenyValues {
		opts.Rules = append(opts.Rules, secret.PolicyRule{
			Name:        fmt.Sprintf("deny-value-%d", i),
			DenyPattern: pattern,
		})
	}

	violations := secret.EnforcePolicy(secrets, opts)
	fmt.Fprint(cmd.OutOrStdout(), secret.FormatPolicyReport(violations))

	if len(violations) > 0 {
		os.Exit(1)
	}
	return nil
}
