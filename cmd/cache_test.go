package cmd_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourorg/vaultpull/cmd"
	"github.com/yourorg/vaultpull/internal/cache"
)

func runCacheCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	root := cmd.RootCmd
	root.SetOut(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func writeCacheFile(t *testing.T, path string, expired bool) {
	t.Helper()
	ttl := time.Hour
	if expired {
		ttl = time.Millisecond
		defer time.Sleep(5 * time.Millisecond)
	}
	if err := cache.Store(path, map[string]string{"K": "v"}, ttl); err != nil {
		t.Fatal(err)
	}
	if expired {
		time.Sleep(5 * time.Millisecond)
	}
}

func TestCacheStatus_NoCache(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	out, err := runCacheCmd(t, "cache", "status", "--cache-path", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" || !bytes.Contains([]byte(out), []byte("No cache found")) {
		t.Errorf("expected 'No cache found', got: %s", out)
	}
}

func TestCacheStatus_ValidCache(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	writeCacheFile(t, path, false)
	out, err := runCacheCmd(t, "cache", "status", "--cache-path", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains([]byte(out), []byte("valid")) {
		t.Errorf("expected 'valid' in output, got: %s", out)
	}
}

func TestCacheInvalidate_RemovesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	writeCacheFile(t, path, false)
	_, err := runCacheCmd(t, "cache", "invalidate", "--cache-path", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, serr := os.Stat(path); !os.IsNotExist(serr) {
		t.Error("expected cache file to be removed")
	}
}

// Ensure cache entry JSON is well-formed.
func TestStore_JSONWellFormed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	_ = cache.Store(path, map[string]string{"A": "b"}, time.Minute)
	data, _ := os.ReadFile(path)
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Errorf("cache file is not valid JSON: %v", err)
	}
}
