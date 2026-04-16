package secret

import (
	"strings"
	"testing"
	"time"
)

func TestCheckRotation_NoEntries(t *testing.T) {
	results := CheckRotation(nil, DefaultRotateOptions())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestCheckRotation_FreshSecret_NotRotated(t *testing.T) {
	entries := []RotateEntry{
		{Key: "DB_PASS", CreatedAt: time.Now().AddDate(0, 0, -10)},
	}
	results := CheckRotation(entries, DefaultRotateOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].Rotated {
		t.Error("expected fresh secret not to be rotated")
	}
}

func TestCheckRotation_OldSecret_Rotated(t *testing.T) {
	entries := []RotateEntry{
		{Key: "API_KEY", CreatedAt: time.Now().AddDate(0, 0, -100)},
	}
	opts := DefaultRotateOptions()
	results := CheckRotation(entries, opts)
	if !results[0].Rotated {
		t.Error("expected old secret to be rotated")
	}
	if results[0].Reason == "" {
		t.Error("expected reason to be set")
	}
}

func TestCheckRotation_DryRun_NotRotated(t *testing.T) {
	entries := []RotateEntry{
		{Key: "TOKEN", CreatedAt: time.Now().AddDate(0, 0, -200)},
	}
	opts := DefaultRotateOptions()
	opts.DryRun = true
	results := CheckRotation(entries, opts)
	if results[0].Rotated {
		t.Error("dry-run should not mark as rotated")
	}
	if results[0].Reason == "" {
		t.Error("expected reason even in dry-run")
	}
}

func TestCheckRotation_EmptyKey_Skipped(t *testing.T) {
	entries := []RotateEntry{
		{Key: "", CreatedAt: time.Now().AddDate(0, 0, -200)},
	}
	results := CheckRotation(entries, DefaultRotateOptions())
	if len(results) != 0 {
		t.Error("empty key entries should be skipped")
	}
}

func TestFormatRotateReport_ContainsSummary(t *testing.T) {
	results := []RotateResult{
		{Key: "A", Rotated: true, Reason: "age 100 days exceeds limit 90"},
		{Key: "B", Rotated: false, Reason: ""},
	}
	out := FormatRotateReport(results)
	if !strings.Contains(out, "Rotated: 1") {
		t.Errorf("expected rotated count in report, got: %s", out)
	}
	if !strings.Contains(out, "Skipped: 1") {
		t.Errorf("expected skipped count in report, got: %s", out)
	}
}
