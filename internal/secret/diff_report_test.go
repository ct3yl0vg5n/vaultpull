package secret

import (
	"strings"
	"testing"
)

func TestBuildDiffReport_Added(t *testing.T) {
	old := map[string]string{}
	next := map[string]string{"NEW_KEY": "value"}
	entries := BuildDiffReport(old, next, DiffReportOptions{RedactValues: false})
	if len(entries) != 1 || entries[0].Status != "added" {
		t.Fatalf("expected 1 added entry, got %+v", entries)
	}
}

func TestBuildDiffReport_Removed(t *testing.T) {
	old := map[string]string{"OLD_KEY": "val"}
	next := map[string]string{}
	entries := BuildDiffReport(old, next, DiffReportOptions{RedactValues: false})
	if len(entries) != 1 || entries[0].Status != "removed" {
		t.Fatalf("expected 1 removed entry, got %+v", entries)
	}
}

func TestBuildDiffReport_Modified(t *testing.T) {
	old := map[string]string{"KEY": "old"}
	next := map[string]string{"KEY": "new"}
	entries := BuildDiffReport(old, next, DiffReportOptions{RedactValues: false})
	if len(entries) != 1 || entries[0].Status != "modified" {
		t.Fatalf("expected 1 modified entry, got %+v", entries)
	}
	if entries[0].Old != "old" || entries[0].New != "new" {
		t.Errorf("unexpected values: %+v", entries[0])
	}
}

func TestBuildDiffReport_Unchanged(t *testing.T) {
	old := map[string]string{"KEY": "same"}
	next := map[string]string{"KEY": "same"}
	entries := BuildDiffReport(old, next, DiffReportOptions{RedactValues: false})
	if len(entries) != 1 || entries[0].Status != "unchanged" {
		t.Fatalf("expected 1 unchanged entry, got %+v", entries)
	}
}

func TestBuildDiffReport_Redacted(t *testing.T) {
	old := map[string]string{}
	next := map[string]string{"SECRET": "supersecret"}
	entries := BuildDiffReport(old, next, DefaultDiffReportOptions())
	if entries[0].New == "supersecret" {
		t.Error("expected value to be redacted")
	}
}

func TestFormatDiffReport_Empty(t *testing.T) {
	out := FormatDiffReport(nil)
	if out != "no differences found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatDiffReport_Symbols(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", New: "1", Status: "added"},
		{Key: "B", Old: "2", Status: "removed"},
		{Key: "C", Old: "x", New: "y", Status: "modified"},
	}
	out := FormatDiffReport(entries)
	if !strings.Contains(out, "+ A") {
		t.Error("expected added line")
	}
	if !strings.Contains(out, "- B") {
		t.Error("expected removed line")
	}
	if !strings.Contains(out, "~ C") {
		t.Error("expected modified line")
	}
}
