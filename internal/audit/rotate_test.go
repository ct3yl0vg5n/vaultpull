package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpull/internal/audit"
)

func writeSizedFile(t *testing.T, path string, size int) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	data := make([]byte, size)
	for i := range data {
		data[i] = 'x'
	}
	f.Write(data)
}

func TestRotate_BelowThreshold_NoOp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	writeSizedFile(t, path, 100)

	backup, err := audit.Rotate(path, audit.RotateOptions{MaxBytes: 1024})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backup != "" {
		t.Errorf("expected no rotation, got backup %q", backup)
	}
}

func TestRotate_AboveThreshold_CreatesBackup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	writeSizedFile(t, path, 2048)

	backup, err := audit.Rotate(path, audit.RotateOptions{MaxBytes: 512})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backup == "" {
		t.Fatal("expected a backup path to be returned")
	}

	if _, err := os.Stat(backup); err != nil {
		t.Errorf("backup file not found: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("original path should exist after rotation: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("new log file should be empty, got size %d", info.Size())
	}
}

func TestRotate_FileNotExist_NoOp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.jsonl")

	backup, err := audit.Rotate(path, audit.RotateOptions{MaxBytes: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backup != "" {
		t.Errorf("expected no backup for missing file, got %q", backup)
	}
}

func TestRotate_ZeroMaxBytes_AlwaysRotates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	writeSizedFile(t, path, 1)

	backup, err := audit.Rotate(path, audit.RotateOptions{MaxBytes: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backup == "" {
		t.Error("expected rotation when MaxBytes is 0")
	}
}
