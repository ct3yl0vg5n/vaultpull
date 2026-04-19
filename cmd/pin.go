package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/secret"
)

var pinFile string

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Manage pinned secret values",
}

var pinAddCmd = &cobra.Command{
	Use:   "add <key> <value>",
	Short: "Pin a secret key to a specific value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := secret.LoadPins(pinFile)
		if err != nil {
			return err
		}
		secret.Pin(store, args[0], args[1], os.Getenv("USER"))
		if err := secret.SavePins(pinFile, store); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "pinned %s\n", args[0])
		return nil
	},
}

var pinRemoveCmd = &cobra.Command{
	Use:   "remove <key>",
	Short: "Unpin a secret key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := secret.LoadPins(pinFile)
		if err != nil {
			return err
		}
		if !secret.Unpin(store, args[0]) {
			return fmt.Errorf("key %q is not pinned", args[0])
		}
		if err := secret.SavePins(pinFile, store); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "unpinned %s\n", args[0])
		return nil
	},
}

var pinListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pinned secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := secret.LoadPins(pinFile)
		if err != nil {
			return err
		}
		if len(store.Entries) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "no pinned secrets")
			return nil
		}
		for _, e := range store.Entries {
			fmt.Fprintf(cmd.OutOrStdout(), "%s = %s (pinned %s)\n", e.Key, e.Value, e.PinnedAt.Format("2006-01-02"))
		}
		return nil
	},
}

func init() {
	pinCmd.PersistentFlags().StringVar(&pinFile, "pin-file", ".vault-pins.json", "path to pin store file")
	pinCmd.AddCommand(pinAddCmd, pinRemoveCmd, pinListCmd)
	rootCmd.AddCommand(pinCmd)
}
