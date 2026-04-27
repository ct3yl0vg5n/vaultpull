package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// DefaultEncryptOptions returns sensible defaults for encryption.
func DefaultEncryptOptions() EncryptOptions {
	return EncryptOptions{
		Enabled: true,
	}
}

// EncryptOptions controls how values are encrypted or decrypted.
type EncryptOptions struct {
	Enabled bool
	Key     []byte // must be 16, 24, or 32 bytes for AES-128/192/256
}

// Encrypt encrypts a plaintext string using AES-GCM and returns a base64-encoded ciphertext.
func Encrypt(plaintext string, opts EncryptOptions) (string, error) {
	if !opts.Enabled {
		return plaintext, nil
	}
	if len(opts.Key) == 0 {
		return "", errors.New("encrypt: key must not be empty")
	}
	block, err := aes.NewCipher(opts.Key)
	if err != nil {
		return "", fmt.Errorf("encrypt: new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: new gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encrypt: nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded AES-GCM ciphertext and returns the plaintext.
func Decrypt(encoded string, opts EncryptOptions) (string, error) {
	if !opts.Enabled {
		return encoded, nil
	}
	if len(opts.Key) == 0 {
		return "", errors.New("decrypt: key must not be empty")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decrypt: base64 decode: %w", err)
	}
	block, err := aes.NewCipher(opts.Key)
	if err != nil {
		return "", fmt.Errorf("decrypt: new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("decrypt: new gcm: %w", err)
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", errors.New("decrypt: ciphertext too short")
	}
	plaintext, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: open: %w", err)
	}
	return string(plaintext), nil
}

// EncryptMap encrypts all values in the map, returning a new map.
func EncryptMap(src map[string]string, opts EncryptOptions) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		enc, err := Encrypt(v, opts)
		if err != nil {
			return nil, fmt.Errorf("encrypt map key %q: %w", k, err)
		}
		out[k] = enc
	}
	return out, nil
}

// DecryptMap decrypts all values in the map, returning a new map.
func DecryptMap(src map[string]string, opts EncryptOptions) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		dec, err := Decrypt(v, opts)
		if err != nil {
			return nil, fmt.Errorf("decrypt map key %q: %w", k, err)
		}
		out[k] = dec
	}
	return out, nil
}
