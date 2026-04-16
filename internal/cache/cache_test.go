package cache_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/vaultpull/internal/cache"
)

func tempCachePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "cache", "secrets.json")
}

func TestStore_And_Load(t *testing.T) {
	path := tempCachePath(t)
	secrets := map[string]string{"KEY": "value", "OTHER": "123"}

	if err := cache.Store(path, secrets, 10*time.Minute); err != nil {
		t.Fatalf("Store: %v", err)
	}

	entry, err := cache.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if entry.Secrets["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", entry.Secrets["KEY"])
	}
}

func TestLoad_NotExist(t *testing.T) {
	entry, err := cache.Load("/tmp/vaultpull_no_such_file_xyz.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Fatal("expected nil entry for missing file")
	}
}

func TestIsExpired_WithShortTTL(t *testing.T) {
	path := tempCachePath(t)
	if err := cache.Store(path, map[string]string{"A": "1"}, 1*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Millisecond)
	entry, _ := cache.Load(path)
	if !entry.IsExpired() {
		t.Error("expected entry to be expired")
	}
}

func TestIsExpired_WithLongTTL(t *testing.T) {
	path := tempCachePath(t)
	if err := cache.Store(path, map[string]string{"A": "1"}, 1*time.Hour); err != nil {
		t.Fatal(err)
	}
	entry, _ := cache.Load(path)
	if entry.IsExpired() {
		t.Error("expected entry to not be expired")
	}
}

func TestInvalidate(t *testing.T) {
	path := tempCachePath(t)
	_ = cache.Store(path, map[string]string{"X": "y"}, time.Minute)
	if err := cache.Invalidate(path); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestInvalidate_NoOp_WhenMissing(t *testing.T) {
	if err := cache.Invalidate("/tmp/vaultpull_ghost_cache.json"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
