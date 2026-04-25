package secret

import (
	"os"
	"path/filepath"
	"testing"
)

func tempSnapshotPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	path := tempSnapshotPath(t)
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := SaveSnapshot(path, "test-label", secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if s.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", s.Label)
	}
	if len(s.Secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(s.Secrets))
	}
	if s.Secrets["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", s.Secrets["FOO"])
	}
}

func TestLoadSnapshot_NotExist(t *testing.T) {
	s, err := LoadSnapshot("/nonexistent/snapshot.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != nil {
		t.Error("expected nil snapshot for missing file")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	path := tempSnapshotPath(t)
	os.WriteFile(path, []byte("not-json"), 0600)
	_, err := LoadSnapshot(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestFormatSnapshotSummary_WithLabel(t *testing.T) {
	path := tempSnapshotPath(t)
	SaveSnapshot(path, "prod", map[string]string{"A": "1"})
	s, _ := LoadSnapshot(path)
	out := FormatSnapshotSummary(s)
	if out == "" {
		t.Error("expected non-empty summary")
	}
}

func TestFormatSnapshotSummary_Nil(t *testing.T) {
	out := FormatSnapshotSummary(nil)
	if out != "no snapshot found" {
		t.Errorf("unexpected: %q", out)
	}
}

func TestFormatSnapshotSummary_NoLabel(t *testing.T) {
	path := tempSnapshotPath(t)
	SaveSnapshot(path, "", map[string]string{"X": "y"})
	s, _ := LoadSnapshot(path)
	out := FormatSnapshotSummary(s)
	if out == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSaveSnapshot_CreatesFileWithRestrictedPermissions(t *testing.T) {
	path := tempSnapshotPath(t)
	if err := SaveSnapshot(path, "perm-test", map[string]string{"K": "v"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("could not stat snapshot file: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected file permissions 0600, got %04o", perm)
	}
}
