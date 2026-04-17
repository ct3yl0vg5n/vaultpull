package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourorg/vaultpull/internal/cache"
)

var cachePath string

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage the local secrets cache",
}

var cacheInvalidateCmd = &cobra.Command{
	Use:   "invalidate",
	Short: "Remove the local cache file",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cachePath
		if path == "" {
			path = cache.DefaultOptions().Path
		}
		if err := cache.Invalidate(path); err != nil {
			return fmt.Errorf("invalidate cache: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Cache invalidated: %s\n", path)
		return nil
	},
}

var cacheStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cache status (exists / expired)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cachePath
		if path == "" {
			path = cache.DefaultOptions().Path
		}
		entry, err := cache.Load(path)
		if err != nil {
			return fmt.Errorf("load cache: %w", err)
		}
		if entry == nil {
			fmt.Fprintln(cmd.OutOrStdout(), "No cache found.")
			return nil
		}
		status := "valid"
		if entry.IsExpired() {
			status = "expired"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Cache status: %s (fetched at %s, TTL %s)\n",
			status, entry.FetchedAt.Format("2006-01-02 15:04:05"), entry.TTL)
		return nil
	},
}

// resolveCachePath returns the provided path if non-empty, otherwise falls back
// to the default cache path defined by cache.DefaultOptions.
func resolveCachePath(p string) string {
	if p != "" {
		return p
	}
	return cache.DefaultOptions().Path
}

func init() {
	cacheCmd.PersistentFlags().StringVar(&cachePath, "cache-path", "", "path to cache file")
	cacheCmd.AddCommand(cacheInvalidateCmd)
	cacheCmd.AddCommand(cacheStatusCmd)
	RootCmd.AddCommand(cacheCmd)
}
