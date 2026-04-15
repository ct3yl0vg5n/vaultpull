# vaultpull

> CLI tool to sync HashiCorp Vault secrets to local `.env` files with diff preview

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Authenticate with Vault and pull secrets into a local `.env` file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

vaultpull --path secret/data/myapp --out .env
```

Before writing, `vaultpull` shows a diff of what will change:

```
~ DB_PASSWORD  [changed]
+ NEW_API_KEY  [added]
- OLD_SECRET   [removed]

Apply changes? [y/N]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | *(required)* |
| `--out` | Output `.env` file path | `.env` |
| `--yes` | Skip confirmation prompt | `false` |
| `--addr` | Vault server address | `$VAULT_ADDR` |

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance
- A valid Vault token or auth method configured

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername