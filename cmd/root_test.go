package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd_MissingRequiredFlag(t *testing.T) {
	// Reset flags to defaults before test
	vaultAddr = ""
	vaultToken = ""
	mountPath = "secret"
	secretPath = ""
	envFile = ".env"
	dryRun = false

	rootCmd.SetArgs([]string{})

	buf := &bytes.Buffer{}
	rootCmd.SetErr(buf)

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --path flag is missing, got nil")
	}
}

func TestRootCmd_FlagDefaults(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "test-token")

	// Re-run init logic by checking env fallback values directly
	addr := rootCmd.Flags().Lookup("vault-addr")
	if addr == nil {
		t.Fatal("expected vault-addr flag to be registered")
	}

	token := rootCmd.Flags().Lookup("vault-token")
	if token == nil {
		t.Fatal("expected vault-token flag to be registered")
	}

	mount := rootCmd.Flags().Lookup("mount")
	if mount == nil {
		t.Fatal("expected mount flag to be registered")
	}
	if mount.DefValue != "secret" {
		t.Errorf("expected mount default to be 'secret', got %q", mount.DefValue)
	}

	ef := rootCmd.Flags().Lookup("env-file")
	if ef == nil {
		t.Fatal("expected env-file flag to be registered")
	}
	if ef.DefValue != ".env" {
		t.Errorf("expected env-file default to be '.env', got %q", ef.DefValue)
	}
}

func TestRootCmd_DryRunFlag(t *testing.T) {
	flag := rootCmd.Flags().Lookup("dry-run")
	if flag == nil {
		t.Fatal("expected dry-run flag to be registered")
	}
	if flag.DefValue != "false" {
		t.Errorf("expected dry-run default to be 'false', got %q", flag.DefValue)
	}
}
