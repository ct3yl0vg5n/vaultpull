package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/vaultpull/internal/audit"
)

func TestWrite_CreatesFileAndAppendsEntries(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "audit-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	logger := audit.NewLogger(tmp.Name())

	entry1 := audit.Entry{
		Timestamp: time.Now().UTC(),
		Operation: "sync",
		Mount:     "secret",
		Path:      "myapp/prod",
		EnvFile:   ".env",
		DryRun:    false,
		Added:     2,
		Modified:  1,
	}
	entry2 := audit.Entry{
		Operation: "sync",
		DryRun:    true,
		Error:     "vault unreachable",
	}

	if err := logger.Write(entry1); err != nil {
		t.Fatalf("Write entry1: %v", err)
	}
	if err := logger.Write(entry2); err != nil {
		t.Fatalf("Write entry2: %v", err)
	}

	f, err := os.Open(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var entries []audit.Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var e audit.Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		entries = append(entries, e)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Added != 2 {
		t.Errorf("expected Added=2, got %d", entries[0].Added)
	}
	if entries[1].Error != "vault unreachable" {
		t.Errorf("unexpected error field: %q", entries[1].Error)
	}
	if entries[1].Timestamp.IsZero() {
		t.Error("expected timestamp to be populated automatically")
	}
}

func TestWrite_NoPath_IsNoop(t *testing.T) {
	logger := audit.NewLogger("")
	if err := logger.Write(audit.Entry{Operation: "sync"}); err != nil {
		t.Errorf("expected no error for no-op logger, got %v", err)
	}
}

func TestWrite_InvalidPath_ReturnsError(t *testing.T) {
	logger := audit.NewLogger("/nonexistent-dir/audit.jsonl")
	err := logger.Write(audit.Entry{Operation: "sync"})
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
