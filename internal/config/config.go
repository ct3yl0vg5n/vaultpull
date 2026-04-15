package config

import (
	"errors"
	"os"
	"strings"

	"github.com/your-org/vaultpull/internal/filter"
)

// Config holds all runtime configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	MountPath  string
	SecretPath string
	EnvFile    string
	DryRun     bool
	Filter     filter.Rule
}

// FromEnv builds a Config from environment variables.
func FromEnv() (*Config, error) {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		return nil, errors.New("VAULT_ADDR is required")
	}
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, errors.New("VAULT_TOKEN is required")
	}

	mount := os.Getenv("VAULT_MOUNT")
	if mount == "" {
		mount = "secret"
	}

	envFile := os.Getenv("VAULTPULL_ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	excludeRaw := os.Getenv("VAULTPULL_EXCLUDE")
	var excludeKeys []string
	if excludeRaw != "" {
		for _, k := range strings.Split(excludeRaw, ",") {
			if trimmed := strings.TrimSpace(k); trimmed != "" {
				excludeKeys = append(excludeKeys, trimmed)
			}
		}
	}

	return &Config{
		VaultAddr:  strings.TrimRight(addr, "/"),
		VaultToken: token,
		MountPath:  mount,
		SecretPath: os.Getenv("VAULT_SECRET_PATH"),
		EnvFile:    envFile,
		Filter: filter.Rule{
			Prefix:      os.Getenv("VAULTPULL_PREFIX"),
			Exclude:     excludeKeys,
			StripPrefix: os.Getenv("VAULTPULL_STRIP_PREFIX") == "true",
		},
	}, nil
}
