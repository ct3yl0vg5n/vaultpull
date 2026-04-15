package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds the resolved configuration for a sync operation.
type Config struct {
	VaultAddr  string
	VaultToken string
	Mount      string
	SecretPath string
	EnvFile    string
	DryRun     bool
}

// FromEnv populates configuration fields from environment variables,
// returning an error if any required field is missing.
func FromEnv() (*Config, error) {
	cfg := &Config{
		VaultAddr:  strings.TrimRight(os.Getenv("VAULT_ADDR"), "/"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
		Mount:      os.Getenv("VAULT_MOUNT"),
		SecretPath: os.Getenv("VAULT_SECRET_PATH"),
		EnvFile:    os.Getenv("VAULTPULL_ENV_FILE"),
	}

	if cfg.Mount == "" {
		cfg.Mount = "secret"
	}

	if cfg.EnvFile == "" {
		cfg.EnvFile = ".env"
	}

	return cfg, cfg.validate()
}

// Merge overlays non-zero flag values onto the config, allowing CLI flags
// to override environment-derived values.
func (c *Config) Merge(addr, token, mount, secretPath, envFile string, dryRun bool) {
	if addr != "" {
		c.VaultAddr = strings.TrimRight(addr, "/")
	}
	if token != "" {
		c.VaultToken = token
	}
	if mount != "" {
		c.Mount = mount
	}
	if secretPath != "" {
		c.SecretPath = secretPath
	}
	if envFile != "" {
		c.EnvFile = envFile
	}
	c.DryRun = dryRun
}

func (c *Config) validate() error {
	if c.VaultAddr == "" {
		return errors.New("config: VAULT_ADDR is required")
	}
	if c.VaultToken == "" {
		return errors.New("config: VAULT_TOKEN is required")
	}
	if c.SecretPath == "" {
		return errors.New("config: VAULT_SECRET_PATH is required")
	}
	return nil
}
