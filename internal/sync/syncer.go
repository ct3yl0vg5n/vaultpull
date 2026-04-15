package sync

import (
	"fmt"

	"github.com/user/vaultpull/internal/diff"
	"github.com/user/vaultpull/internal/env"
	"github.com/user/vaultpull/internal/vault"
)

// Options configures the sync operation.
type Options struct {
	DryRun    bool
	Quiet     bool
	EnvFile   string
	VaultPath string
}

// Result holds the outcome of a sync operation.
type Result struct {
	Changes  []diff.Change
	Applied  bool
	EnvFile  string
}

// Syncer orchestrates pulling secrets from Vault and writing them to a .env file.
type Syncer struct {
	client *vault.Client
	opts   Options
}

// New creates a new Syncer with the given Vault client and options.
func New(client *vault.Client, opts Options) *Syncer {
	return &Syncer{client: client, opts: opts}
}

// Run executes the sync: fetches remote secrets, compares with local, and optionally writes.
func (s *Syncer) Run() (*Result, error) {
	remote, err := s.client.GetSecrets(s.opts.VaultPath)
	if err != nil {
		return nil, fmt.Errorf("fetching secrets from vault: %w", err)
	}

	local, err := env.ParseFile(s.opts.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("parsing local env file: %w", err)
	}

	changes := diff.Compare(local, remote)

	result := &Result{
		Changes: changes,
		EnvFile: s.opts.EnvFile,
	}

	if !diff.HasChanges(changes) || s.opts.DryRun {
		return result, nil
	}

	if err := env.WriteFile(s.opts.EnvFile, remote); err != nil {
		return nil, fmt.Errorf("writing env file: %w", err)
	}

	result.Applied = true
	return result, nil
}
