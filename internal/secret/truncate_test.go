package secret

import (
	"strings"
	"testing"
)

func TestTruncateMap_NoChanges(t *testing.T) {
	opts := DefaultTruncateOptions()
	src := map[string]string{"KEY": "short"}
	out, results := TruncateMap(src, opts)
	if out["KEY"] != "short" {
		t.Errorf("expected 'short', got %q", out["KEY"])
	}
	for _, r := range results {
		if r.Changed {
			t.Error("expected no changes")
		}
	}
}

func TestTruncateMap_TruncatesLongValue(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.MaxLength = 10
	long := "this_is_a_very_long_value"
	src := map[string]string{"KEY": long}
	out, results := TruncateMap(src, opts)
	if len(out["KEY"]) > 10 {
		t.Errorf("expected truncated value, got %q", out["KEY"])
	}
	if !strings.HasSuffix(out["KEY"], "...") {
		t.Errorf("expected suffix '...', got %q", out["KEY"])
	}
	if len(results) == 0 || !results[0].Changed {
		t.Error("expected Changed=true in result")
	}
}

func TestTruncateMap_SkipKeys(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.MaxLength = 5
	opts.SkipKeys = []string{"SKIP_ME"}
	src := map[string]string{
		"SKIP_ME": "this_should_not_be_truncated",
		"OTHER":   "this_should_be_truncated",
	}
	out, _ := TruncateMap(src, opts)
	if out["SKIP_ME"] != "this_should_not_be_truncated" {
		t.Errorf("SKIP_ME should not be truncated, got %q", out["SKIP_ME"])
	}
	if len(out["OTHER"]) > 5 {
		t.Errorf("OTHER should be truncated, got %q", out["OTHER"])
	}
}

func TestTruncateMap_CustomSuffix(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.MaxLength = 8
	opts.Suffix = "~~"
	src := map[string]string{"K": "abcdefghij"}
	out, _ := TruncateMap(src, opts)
	if !strings.HasSuffix(out["K"], "~~") {
		t.Errorf("expected suffix '~~', got %q", out["K"])
	}
}

func TestTruncateMap_ZeroMaxLength_TruncatesAll(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.MaxLength = 0
	src := map[string]string{"K": "somevalue"}
	out, _ := TruncateMap(src, opts)
	if out["K"] != "somevalue" {
		t.Errorf("zero MaxLength should be no-op, got %q", out["K"])
	}
}

func TestFormatTruncateReport_NoChanges(t *testing.T) {
	report := FormatTruncateReport([]TruncateResult{
		{Key: "A", Original: "x", Truncated: "x", Changed: false},
	})
	if !strings.Contains(report, "No values truncated") {
		t.Errorf("expected no-change message, got %q", report)
	}
}

func TestFormatTruncateReport_WithChanges(t *testing.T) {
	report := FormatTruncateReport([]TruncateResult{
		{Key: "FOO", Original: "longvalue", Truncated: "lon...", Changed: true},
	})
	if !strings.Contains(report, "FOO") {
		t.Errorf("expected key in report, got %q", report)
	}
	if !strings.Contains(report, "Truncated 1") {
		t.Errorf("expected count in report, got %q", report)
	}
}
