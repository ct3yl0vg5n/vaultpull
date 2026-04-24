package secret

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tempArchivePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "archive.json")
}

func TestLoadArchive_NotExist(t *testing.T) {
	entries, err := LoadArchive("/nonexistent/archive.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestAppendArchive_CreatesAndReads(t *testing.T) {
	path := tempArchivePath(t)
	opts := DefaultArchiveOptions()
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}

	if err := AppendArchive(path, "v1", secrets, opts); err != nil {
		t.Fatalf("AppendArchive: %v", err)
	}

	entries, err := LoadArchive(path)
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Label != "v1" {
		t.Errorf("expected label 'v1', got %q", entries[0].Label)
	}
	if entries[0].Secrets["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", entries[0].Secrets["FOO"])
	}
}

func TestAppendArchive_RespectsMaxEntries(t *testing.T) {
	path := tempArchivePath(t)
	opts := ArchiveOptions{MaxEntries: 3}
	secrets := map[string]string{"KEY": "val"}

	for i := 0; i < 5; i++ {
		if err := AppendArchive(path, "label", secrets, opts); err != nil {
			t.Fatalf("AppendArchive iteration %d: %v", i, err)
		}
	}

	entries, err := LoadArchive(path)
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries (max), got %d", len(entries))
	}
}

func TestLoadArchive_InvalidJSON(t *testing.T) {
	path := tempArchivePath(t)
	if err := os.WriteFile(path, []byte("not-json"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadArchive(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFormatArchiveSummary_Empty(t *testing.T) {
	out := FormatArchiveSummary([]ArchiveEntry{})
	if !strings.Contains(out, "no archive entries") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatArchiveSummary_WithEntries(t *testing.T) {
	path := tempArchivePath(t)
	opts := DefaultArchiveOptions()
	_ = AppendArchive(path, "release-1", map[string]string{"A": "1", "B": "2"}, opts)
	entries, _ := LoadArchive(path)
	out := FormatArchiveSummary(entries)
	if !strings.Contains(out, "release-1") {
		t.Errorf("expected label in output, got: %q", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected key count in output, got: %q", out)
	}
}
