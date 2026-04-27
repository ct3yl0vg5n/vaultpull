package secret

import (
	"strings"
	"testing"
)

func makeKey(t *testing.T) []byte {
	t.Helper()
	// 32 bytes => AES-256
	return []byte("12345678901234567890123456789012")
}

func TestEncrypt_RoundTrip(t *testing.T) {
	opts := DefaultEncryptOptions()
	opts.Key = makeKey(t)

	plain := "super-secret-value"
	enc, err := Encrypt(plain, opts)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if enc == plain {
		t.Fatal("expected ciphertext to differ from plaintext")
	}

	dec, err := Decrypt(enc, opts)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec != plain {
		t.Errorf("got %q, want %q", dec, plain)
	}
}

func TestEncrypt_Disabled_Passthrough(t *testing.T) {
	opts := EncryptOptions{Enabled: false}
	out, err := Encrypt("hello", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello" {
		t.Errorf("expected passthrough, got %q", out)
	}
}

func TestEncrypt_EmptyKey_Error(t *testing.T) {
	opts := EncryptOptions{Enabled: true, Key: nil}
	_, err := Encrypt("value", opts)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestDecrypt_InvalidBase64_Error(t *testing.T) {
	opts := DefaultEncryptOptions()
	opts.Key = makeKey(t)
	_, err := Decrypt("not-valid-base64!!!", opts)
	if err == nil {
		t.Fatal("expected base64 decode error")
	}
}

func TestDecrypt_TruncatedCiphertext_Error(t *testing.T) {
	opts := DefaultEncryptOptions()
	opts.Key = makeKey(t)
	// base64 of a very short byte slice
	_, err := Decrypt("YQ==", opts)
	if err == nil || !strings.Contains(err.Error(), "too short") {
		t.Fatalf("expected 'too short' error, got %v", err)
	}
}

func TestEncryptMap_And_DecryptMap(t *testing.T) {
	opts := DefaultEncryptOptions()
	opts.Key = makeKey(t)

	src := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	}

	enc, err := EncryptMap(src, opts)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	for k, v := range enc {
		if v == src[k] {
			t.Errorf("key %q: expected encrypted value, got plaintext", k)
		}
	}

	dec, err := DecryptMap(enc, opts)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}
	for k, want := range src {
		if got := dec[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestEncrypt_NonceDiffers(t *testing.T) {
	opts := DefaultEncryptOptions()
	opts.Key = makeKey(t)

	enc1, _ := Encrypt("same", opts)
	enc2, _ := Encrypt("same", opts)
	if enc1 == enc2 {
		t.Error("expected different ciphertexts due to random nonce")
	}
}
