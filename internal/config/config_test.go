package config

import (
	"testing"
)

func setEnv(t *testing.T, pairs map[string]string) {
	t.Helper()
	for k, v := range pairs {
		t.Setenv(k, v)
	}
}

func TestFromEnv_Valid(t *testing.T) {
	setEnv(t, map[string]string{
		"VAULT_ADDR":        "http://localhost:8200",
		"VAULT_TOKEN":       "root",
		"VAULT_SECRET_PATH": "myapp/prod",
	})

	cfg, err := FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://localhost:8200" {
		t.Errorf("expected addr http://localhost:8200, got %s", cfg.VaultAddr)
	}
	if cfg.Mount != "secret" {
		t.Errorf("expected default mount 'secret', got %s", cfg.Mount)
	}
	if cfg.EnvFile != ".env" {
		t.Errorf("expected default env file '.env', got %s", cfg.EnvFile)
	}
}

func TestFromEnv_MissingAddr(t *testing.T) {
	setEnv(t, map[string]string{
		"VAULT_TOKEN":       "root",
		"VAULT_SECRET_PATH": "myapp/prod",
	})
	_, err := FromEnv()
	if err == nil {
		t.Fatal("expected error for missing VAULT_ADDR")
	}
}

func TestFromEnv_MissingToken(t *testing.T) {
	setEnv(t, map[string]string{
		"VAULT_ADDR":        "http://localhost:8200",
		"VAULT_SECRET_PATH": "myapp/prod",
	})
	_, err := FromEnv()
	if err == nil {
		t.Fatal("expected error for missing VAULT_TOKEN")
	}
}

func TestFromEnv_TrailingSlashStripped(t *testing.T) {
	setEnv(t, map[string]string{
		"VAULT_ADDR":        "http://localhost:8200/",
		"VAULT_TOKEN":       "root",
		"VAULT_SECRET_PATH": "myapp/prod",
	})
	cfg, err := FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://localhost:8200" {
		t.Errorf("trailing slash not stripped, got %s", cfg.VaultAddr)
	}
}

func TestMerge_OverridesValues(t *testing.T) {
	cfg := &Config{
		VaultAddr:  "http://old:8200",
		VaultToken: "old-token",
		Mount:      "secret",
		SecretPath: "old/path",
		EnvFile:    ".env",
	}
	cfg.Merge("http://new:8200", "new-token", "kv", "new/path", ".env.local", true)

	if cfg.VaultAddr != "http://new:8200" {
		t.Errorf("addr not overridden")
	}
	if cfg.Mount != "kv" {
		t.Errorf("mount not overridden")
	}
	if !cfg.DryRun {
		t.Errorf("dry-run not set")
	}
}
