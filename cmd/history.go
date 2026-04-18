package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/internal/secret"
)

var (
	historyFile string
	historyKey  string
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View secret change history",
	RunE:  runHistory,
}

func init() {
	historyCmd.Flags().StringVar(&historyFile, "history-file", ".vaultpull_history.json", "Path to history log file")
	historyCmd.Flags().StringVar(&historyKey, "key", "", "Filter history by secret key")
	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, _ []string) error {
	h, err := secret.LoadHistory(historyFile)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}

	entries := h.Entries
	if historyKey != "" {
		entries = h.FilterByKey(historyKey)
	}

	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No history entries found.")
		return nil
	}

	w := cmd.OutOrStdout()
	for _, e := range entries {
		fmt.Fprintf(w, "[%s] %s = %s\n", e.Timestamp.Format("2006-01-02 15:04:05"), e.Key, e.Value)
	}
	return nil
}

var _ = os.Stderr // ensure os import used
