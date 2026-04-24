package secret

import (
	"strings"
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"C": "3"}
	r, err := Merge(base, incoming, DefaultMergeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Merged["A"] != "1" || r.Merged["B"] != "2" || r.Merged["C"] != "3" {
		t.Errorf("unexpected merged map: %v", r.Merged)
	}
	if len(r.Added) != 1 || r.Added[0] != "C" {
		t.Errorf("expected Added=[C], got %v", r.Added)
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", r.Conflicts)
	}
}

func TestMerge_ConflictStrategyTheirs(t *testing.T) {
	base := map[string]string{"X": "old"}
	incoming := map[string]string{"X": "new"}
	opts := MergeOptions{Strategy: MergeStrategyTheirs}
	r, err := Merge(base, incoming, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Merged["X"] != "new" {
		t.Errorf("expected 'new', got %q", r.Merged["X"])
	}
	if len(r.Conflicts) != 1 || r.Conflicts[0] != "X" {
		t.Errorf("expected conflict on X, got %v", r.Conflicts)
	}
}

func TestMerge_ConflictStrategyOurs(t *testing.T) {
	base := map[string]string{"X": "old"}
	incoming := map[string]string{"X": "new"}
	opts := MergeOptions{Strategy: MergeStrategyOurs}
	r, err := Merge(base, incoming, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Merged["X"] != "old" {
		t.Errorf("expected 'old', got %q", r.Merged["X"])
	}
}

func TestMerge_ConflictStrategyError(t *testing.T) {
	base := map[string]string{"KEY": "v1"}
	incoming := map[string]string{"KEY": "v2"}
	opts := MergeOptions{Strategy: MergeStrategyError}
	_, err := Merge(base, incoming, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "KEY") {
		t.Errorf("error should mention key name, got: %v", err)
	}
}

func TestMerge_IgnoreKeys(t *testing.T) {
	base := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "99", "B": "2"}
	opts := MergeOptions{Strategy: MergeStrategyTheirs, IgnoreKeys: []string{"B"}}
	r, err := Merge(base, incoming, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Merged["B"]; ok {
		t.Error("B should have been ignored")
	}
	if len(r.Skipped) != 1 || r.Skipped[0] != "B" {
		t.Errorf("expected Skipped=[B], got %v", r.Skipped)
	}
}

func TestFormatMergeReport_ContainsSummary(t *testing.T) {
	r := MergeResult{
		Merged:    map[string]string{"A": "1", "B": "2", "C": "3"},
		Added:     []string{"C"},
		Conflicts: []string{"B"},
		Skipped:   []string{},
	}
	out := FormatMergeReport(r)
	if !strings.Contains(out, "3 keys total") {
		t.Errorf("expected total count in report, got: %s", out)
	}
	if !strings.Contains(out, "added") {
		t.Errorf("expected 'added' section in report, got: %s", out)
	}
	if !strings.Contains(out, "conflicts") {
		t.Errorf("expected 'conflicts' section in report, got: %s", out)
	}
}
