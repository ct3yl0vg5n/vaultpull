package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/env"
	"github.com/vaultpull/vaultpull/internal/secret"
)

var (
	groupFile      string
	groupDelimiter string
	groupMaxDepth  int
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Group secrets by key prefix",
	Long:  `Reads a .env file and groups secrets by their key prefix, displaying a structured summary.`,
	RunE:  runGroup,
}

func init() {
	groupCmd.Flags().StringVarP(&groupFile, "file", "f", ".env", "path to the .env file")
	groupCmd.Flags().StringVar(&groupDelimiter, "delimiter", "_", "delimiter used to split key segments")
	groupCmd.Flags().IntVar(&groupMaxDepth, "depth", 1, "number of prefix segments that form the group name")
	rootCmd.AddCommand(groupCmd)
}

func runGroup(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ParseFile(groupFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(cmd.OutOrStdout(), "file not found: %s\n", groupFile)
			return nil
		}
		return fmt.Errorf("reading env file: %w", err)
	}

	opts := secret.GroupOptions{
		Delimiter: groupDelimiter,
		MaxDepth:  groupMaxDepth,
	}

	result := secret.Group(secrets, opts)
	report := secret.FormatGroupReport(result)
	fmt.Fprint(cmd.OutOrStdout(), report)
	return nil
}
