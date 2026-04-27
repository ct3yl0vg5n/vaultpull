package secret_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/secret"
)

func tempAuditPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit.jsonl")
}

func TestAppendAuditEntry_CreatesFile(t *testing.T) {
	path := tempAuditPath(t)
	opts := secret.DefaultAuditTrailOptions()
	opts.Path = path

	err := secret.AppendAuditEntry(opts, "sync", "DB_PASSWORD", "applied")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected audit file to be created")
	}
}

func TestAppendAuditEntry_AppendsMultiple(t *testing.T) {
	path := tempAuditPath(t)
	opts := secret.DefaultAuditTrailOptions()
	opts.Path = path

	actions := []struct {
		action, key, status string
	}{
		{"sync", "API_KEY", "applied"},
		{"rotate", "DB_PASS", "skipped"},
		{"promote", "SECRET_TOKEN", "applied"},
	}

	for _, a := range actions {
		if err := secret.AppendAuditEntry(opts, a.action, a.key, a.status); err != nil {
			t.Fatalf("AppendAuditEntry(%q): %v", a.key, err)
		}
	}

	entries, err := secret.LoadAuditTrail(opts)
	if err != nil {
		t.Fatalf("LoadAuditTrail: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestLoadAuditTrail_NotExist(t *testing.T) {
	opts := secret.DefaultAuditTrailOptions()
	opts.Path = "/nonexistent/audit.jsonl"

	entries, err := secret.LoadAuditTrail(opts)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestLoadAuditTrail_RespectsMaxEntries(t *testing.T) {
	path := tempAuditPath(t)
	opts := secret.DefaultAuditTrailOptions()
	opts.Path = path
	opts.MaxEntries = 2

	for i := 0; i < 5; i++ {
		_ = secret.AppendAuditEntry(opts, "sync", "KEY", "applied")
	}

	entries, err := secret.LoadAuditTrail(opts)
	if err != nil {
		t.Fatalf("LoadAuditTrail: %v", err)
	}
	if len(entries) > 2 {
		t.Fatalf("expected at most 2 entries, got %d", len(entries))
	}
}

func TestFormatAuditTrail_ContainsKeyAndAction(t *testing.T) {
	path := tempAuditPath(t)
	opts := secret.DefaultAuditTrailOptions()
	opts.Path = path

	_ = secret.AppendAuditEntry(opts, "promote", "MY_SECRET", "applied")

	entries, _ := secret.LoadAuditTrail(opts)
	out := secret.FormatAuditTrail(entries)

	if !strings.Contains(out, "MY_SECRET") {
		t.Errorf("expected output to contain key name, got:\n%s", out)
	}
	if !strings.Contains(out, "promote") {
		t.Errorf("expected output to contain action, got:\n%s", out)
	}
}

func TestFormatAuditTrail_EmptyEntries(t *testing.T) {
	out := secret.FormatAuditTrail(nil)
	if !strings.Contains(out, "no audit") && !strings.Contains(out, "No audit") && out == "" {
		// acceptable: empty string or a "no entries" message
	}
	// Should not panic
}

func TestAppendAuditEntry_TimestampIsRecent(t *testing.T) {
	path := tempAuditPath(t)
	opts := secret.DefaultAuditTrailOptions()
	opts.Path = path

	before := time.Now().Add(-time.Second)
	_ = secret.AppendAuditEntry(opts, "sync", "TOKEN", "applied")
	after := time.Now().Add(time.Second)

	entries, _ := secret.LoadAuditTrail(opts)
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}

	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v outside expected range [%v, %v]", ts, before, after)
	}
}
