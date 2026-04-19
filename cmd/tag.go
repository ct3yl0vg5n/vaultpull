package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/internal/env"
	"github.com/vaultpull/internal/secret"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag secrets in a .env file and filter or report by tag",
	RunE:  runTag,
}

var (
	tagFile       string
	tagFilterTag  string
	tagShowReport bool
)

func init() {
	tagCmd.Flags().StringVar(&tagFile, "file", ".env", "path to .env file")
	tagCmd.Flags().StringVar(&tagFilterTag, "filter", "", "filter keys by tag")
	tagCmd.Flags().BoolVar(&tagShowReport, "report", false, "show full tag report")
	RootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) error {
	secrets, err := env.ParseFile(tagFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", tagFile)
		}
		return err
	}

	opts := secret.DefaultTagOptions()
	for key := range secrets {
		if len(key) > 4 && key[len(key)-4:] == "PASS" || containsSuffix(key, "_KEY", "_SECRET", "_TOKEN") {
			opts = secret.AddTag(opts, key, "sensitive")
		} else {
			opts = secret.AddTag(opts, key, "config")
		}
	}

	if tagFilterTag != "" {
		keys := secret.FilterByTag(opts, tagFilterTag)
		if len(keys) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "no keys found with tag %q\n", tagFilterTag)
			return nil
		}
		for _, k := range keys {
			fmt.Fprintln(cmd.OutOrStdout(), k)
		}
		return nil
	}

	fmt.Fprint(cmd.OutOrStdout(), secret.FormatTagReport(opts))
	return nil
}

func containsSuffix(key string, suffixes ...string) bool {
	for _, s := range suffixes {
		if len(key) >= len(s) && key[len(key)-len(s):] == s {
			return true
		}
	}
	return false
}
