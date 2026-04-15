package sync_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpull/internal/sync"
	"github.com/user/vaultpull/internal/vault"
)

func vaultServer(t *testing.T, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

func tempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if content != "" {
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatalf("failed to write temp env file: %v", err)
		}
	}
	return path
}

func TestRun_DryRun_NoWrite(t *testing.T) {
	server := vaultServer(t, `{"data":{"data":{"KEY":"newval"}}}`)
	defer server.Close()

	client, _ := vault.NewClient(server.URL, "test-token", "")
	envPath := tempEnvFile(t, "KEY=oldval\n")

	s := sync.New(client, sync.Options{
		DryRun:    true,
		EnvFile:   envPath,
		VaultPath: "secret/myapp",
	})

	result, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Applied {
		t.Error("expected Applied=false in dry-run mode")
	}
	if len(result.Changes) == 0 {
		t.Error("expected at least one change")
	}
}

func TestRun_AppliesChanges(t *testing.T) {
	server := vaultServer(t, `{"data":{"data":{"KEY":"newval"}}}`)
	defer server.Close()

	client, _ := vault.NewClient(server.URL, "test-token", "")
	envPath := tempEnvFile(t, "KEY=oldval\n")

	s := sync.New(client, sync.Options{
		DryRun:    false,
		EnvFile:   envPath,
		VaultPath: "secret/myapp",
	})

	result, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Applied {
		t.Error("expected Applied=true when changes exist and dry-run is off")
	}
}

func TestRun_NoChanges_NotApplied(t *testing.T) {
	server := vaultServer(t, `{"data":{"data":{"KEY":"same"}}}`)
	defer server.Close()

	client, _ := vault.NewClient(server.URL, "test-token", "")
	envPath := tempEnvFile(t, "KEY=same\n")

	s := sync.New(client, sync.Options{
		DryRun:    false,
		EnvFile:   envPath,
		VaultPath: "secret/myapp",
	})

	result, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Applied {
		t.Error("expected Applied=false when there are no changes")
	}
}
