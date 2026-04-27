package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/env"
	"github.com/yourorg/vaultpull/internal/secret"
)

var (
	chainFile    string
	chainSteps   []string
	chainDryRun  bool
	chainOutFile string
)

var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Apply a sequence of transformations to a .env file",
	RunE:  runChain,
}

func init() {
	chainCmd.Flags().StringVarP(&chainFile, "file", "f", ".env", "source .env file")
	chainCmd.Flags().StringArrayVarP(&chainSteps, "step", "s", nil, "transformation steps: uppercase, prefix=<val>")
	chainCmd.Flags().BoolVar(&chainDryRun, "dry-run", false, "preview without writing")
	chainCmd.Flags().StringVarP(&chainOutFile, "out", "o", "", "output file (defaults to source file)")
	rootCmd.AddCommand(chainCmd)
}

func runChain(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ParseFile(chainFile)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	steps, err := buildChainSteps(chainSteps)
	if err != nil {
		return err
	}

	opts := secret.DefaultChainOptions()
	opts.DryRun = chainDryRun

	out, results, err := secret.RunChain(secrets, steps, opts)
	if err != nil {
		return fmt.Errorf("chain failed: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), secret.FormatChainReport(results))

	if chainDryRun {
		return nil
	}

	dest := chainFile
	if chainOutFile != "" {
		dest = chainOutFile
	}
	return env.WriteFile(dest, out)
}

func buildChainSteps(names []string) ([]secret.ChainStep, error) {
	var steps []secret.ChainStep
	for _, name := range names {
		switch {
		case name == "uppercase":
			steps = append(steps, secret.ChainStep{
				Name: "uppercase",
				Action: func(m map[string]string) (map[string]string, error) {
					out := make(map[string]string, len(m))
					for k, v := range m {
						out[k] = strings.ToUpper(v)
					}
					return out, nil
				},
			})
		case strings.HasPrefix(name, "prefix="):
			pfx := strings.TrimPrefix(name, "prefix=")
			steps = append(steps, secret.ChainStep{
				Name: "prefix=" + pfx,
				Action: func(m map[string]string) (map[string]string, error) {
					out := make(map[string]string, len(m))
					for k, v := range m {
						out[k] = pfx + v
					}
					return out, nil
				},
			})
		default:
			fmt.Fprintf(os.Stderr, "warning: unknown step %q, skipping\n", name)
		}
	}
	return steps, nil
}
