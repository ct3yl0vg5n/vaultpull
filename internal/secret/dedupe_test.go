package secret

import (
	"strings"
	"testing"
)

func TestDedupe_NoDuplicates(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	r := Dedupe(src, nil, DefaultDedupeOptions())
	if len(r.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %v", r.Duplicates)
	}
	if len(r.Out) != 3 {
		t.Fatalf("expected 3 keys in output, got %d", len(r.Out))
	}
}

func TestDedupe_CaseSensitive_KeepsFirst(t *testing.T) {
	src := map[string]string{"KEY": "first", "key": "second"}
	ordered := []string{"KEY", "key"}
	opts := DefaultDedupeOptions()
	opts.CaseSensitive = true
	r := Dedupe(src, ordered, opts)
	// case-sensitive: KEY != key, so no duplicates
	if len(r.Duplicates) != 0 {
		t.Fatalf("expected no duplicates in case-sensitive mode, got %v", r.Duplicates)
	}
	if len(r.Out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Out))
	}
}

func TestDedupe_CaseInsensitive_RemovesDuplicate(t *testing.T) {
	src := map[string]string{"KEY": "first", "key": "second"}
	ordered := []string{"KEY", "key"}
	opts := DefaultDedupeOptions()
	opts.CaseSensitive = false
	r := Dedupe(src, ordered, opts)
	if len(r.Duplicates) != 1 || r.Duplicates[0] != "key" {
		t.Fatalf("expected [key] as duplicate, got %v", r.Duplicates)
	}
	if _, ok := r.Out["KEY"]; !ok {
		t.Fatal("expected KEY to be kept")
	}
	if _, ok := r.Out["key"]; ok {
		t.Fatal("expected key to be removed")
	}
}

func TestDedupe_PreferLast(t *testing.T) {
	src := map[string]string{"KEY": "first", "key": "second"}
	ordered := []string{"KEY", "key"}
	opts := DefaultDedupeOptions()
	opts.CaseSensitive = false
	opts.PreferLast = true
	r := Dedupe(src, ordered, opts)
	if v := r.Out["key"]; v != "second" {
		t.Fatalf("expected last value 'second', got %q", v)
	}
	if _, ok := r.Out["KEY"]; ok {
		t.Fatal("expected KEY to be replaced by key")
	}
}

func TestDedupe_ReportOnly_DoesNotMutate(t *testing.T) {
	src := map[string]string{"KEY": "first", "key": "second"}
	ordered := []string{"KEY", "key"}
	opts := DefaultDedupeOptions()
	opts.CaseSensitive = false
	opts.ReportOnly = true
	r := Dedupe(src, ordered, opts)
	if len(r.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate, got %d", len(r.Duplicates))
	}
	if len(r.Out) != 2 {
		t.Fatalf("expected original map untouched (2 keys), got %d", len(r.Out))
	}
}

func TestFormatDedupeReport_NoDuplicates(t *testing.T) {
	r := DedupeResult{}
	out := FormatDedupeReport(r)
	if out != "no duplicate keys found" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestFormatDedupeReport_WithDuplicates(t *testing.T) {
	r := DedupeResult{Duplicates: []string{"foo", "bar"}}
	out := FormatDedupeReport(r)
	if !strings.Contains(out, "2 duplicate") {
		t.Fatalf("expected count in output, got: %q", out)
	}
	if !strings.Contains(out, "foo") || !strings.Contains(out, "bar") {
		t.Fatalf("expected key names in output, got: %q", out)
	}
}
