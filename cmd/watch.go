package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/secret"
)

var (
	watchFile     string
	watchInterval int
	watchMax      int
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Poll a .env file for changes and report diffs",
	RunE:  runWatch,
}

func init() {
	watchCmd.Flags().StringVarP(&watchFile, "file", "f", ".env", "path to .env file to watch")
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 30, "poll interval in seconds")
	watchCmd.Flags().IntVarP(&watchMax, "max", "n", 0, "maximum number of checks (0 = unlimited)")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, _ []string) error {
	opts := secret.DefaultWatchOptions()
	opts.Interval = time.Duration(watchInterval) * time.Second
	opts.MaxChecks = watchMax

	prev, err := env.ParseFile(watchFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading %s: %w", watchFile, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "watching %s every %s\n", watchFile, opts.Interval)

	checks := 0
	for {
		time.Sleep(opts.Interval)
		checks++

		curr, err := env.ParseFile(watchFile)
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(cmd.ErrOrStderr(), "error reading file: %v\n", err)
			continue
		}

		events := secret.DetectChanges(prev, curr, opts)
		report := secret.FormatWatchReport(events)
		fmt.Fprintln(cmd.OutOrStdout(), report)

		prev = curr

		if opts.MaxChecks > 0 && checks >= opts.MaxChecks {
			break
		}
	}
	return nil
}
