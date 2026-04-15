package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultpull/internal/sync"
	"github.com/example/vaultpull/internal/vault"
)

var (
	vaultAddr  string
	vaultToken string
	mountPath  string
	secretPath string
	envFile    string
	dryRun     bool
)

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync HashiCorp Vault secrets to local .env files",
	Long: `vaultpull fetches secrets from HashiCorp Vault and syncs them
to a local .env file, showing a diff preview before applying changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Config{
			Address:    vaultAddr,
			Token:      vaultToken,
			MountPath:  mountPath,
			SecretPath: secretPath,
		})
		if err != nil {
			return fmt.Errorf("failed to create vault client: %w", err)
		}

		syncer := sync.New(client, envFile)
		return syncer.Run(cmd.Context(), sync.Options{
			DryRun: dryRun,
		})
	},
}

func init() {
	rootCmd.Flags().StringVar(&vaultAddr, "vault-addr", os.Getenv("VAULT_ADDR"), "Vault server address (env: VAULT_ADDR)")
	rootCmd.Flags().StringVar(&vaultToken, "vault-token", os.Getenv("VAULT_TOKEN"), "Vault token (env: VAULT_TOKEN)")
	rootCmd.Flags().StringVar(&mountPath, "mount", "secret", "Vault KV mount path")
	rootCmd.Flags().StringVar(&secretPath, "path", "", "Secret path within the mount")
	rootCmd.Flags().StringVar(&envFile, "env-file", ".env", "Path to the target .env file")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing to file")

	_ = rootCmd.MarkFlagRequired("path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
