package vault

import (
	"context"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the HashiCorp Vault API client.
type Client struct {
	vc   *vaultapi.Client
	Mount string
}

// NewClient creates a new Vault client using environment variables or explicit config.
func NewClient(address, token, mount string) (*Client, error) {
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if address == "" {
		return nil, fmt.Errorf("vault address is required (set VAULT_ADDR or use --vault-addr flag)")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required (set VAULT_TOKEN or use --vault-token flag)")
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	vc, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}
	vc.SetToken(token)

	if mount == "" {
		mount = "secret"
	}

	return &Client{vc: vc, Mount: mount}, nil
}

// GetSecrets fetches key-value pairs from a KV v2 secret path.
func (c *Client) GetSecrets(ctx context.Context, path string) (map[string]string, error) {
	secret, err := c.vc.KVv2(c.Mount).Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at path %q", path)
	}

	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		result[k] = str
	}
	return result, nil
}
