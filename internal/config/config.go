package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/user/vaultpull/internal/template"
)

// Config holds all runtime configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	Mount      string
	SecretPath string
	EnvFile    string
	AuditLog   string
	Template   template.Rule
}

// FromEnv reads configuration from environment variables.
func FromEnv() (*Config, error) {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		return nil, errors.New("config: VAULT_ADDR is required")
	}
	addr = strings.TrimRight(addr, "/")

	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, errors.New("config: VAULT_TOKEN is required")
	}

	mount := os.Getenv("VAULT_MOUNT")
	if mount == "" {
		mount = "secret"
	}

	envFile := os.Getenv("VAULTPULL_ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	prefix := os.Getenv("VAULTPULL_KEY_PREFIX")
	upper := strings.EqualFold(os.Getenv("VAULTPULL_KEY_UPPER"), "true")

	path := os.Getenv("VAULT_SECRET_PATH")
	if path == "" {
		return nil, fmt.Errorf("config: VAULT_SECRET_PATH is required")
	}

	return &Config{
		VaultAddr:  addr,
		VaultToken: token,
		Mount:      mount,
		SecretPath: path,
		EnvFile:    envFile,
		AuditLog:   os.Getenv("VAULTPULL_AUDIT_LOG"),
		Template: template.Rule{
			Prefix: prefix,
			Upper:  upper,
		},
	}, nil
}
