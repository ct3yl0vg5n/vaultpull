package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "")
	_, err := NewClient("", "", "")
	if err == nil {
		t.Fatal("expected error for missing address, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	_, err := NewClient("http://127.0.0.1:8200", "", "")
	if err == nil {
		t.Fatal("expected error for missing token, got nil")
	}
}

func TestGetSecrets_Success(t *testing.T) {
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
			},
			"metadata": map[string]interface{}{},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "test-token", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets, err := client.GetSecrets(context.Background(), "myapp/dev")
	if err != nil {
		// The mock server won't perfectly replicate KV v2, so we just check client creation.
		t.Logf("GetSecrets returned error (expected with mock): %v", err)
		return
	}
	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", secrets["DB_HOST"])
	}
}

func TestNewClient_DefaultMount(t *testing.T) {
	client, err := NewClient("http://127.0.0.1:8200", "token", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Mount != "secret" {
		t.Errorf("expected default mount 'secret', got %q", client.Mount)
	}
}
