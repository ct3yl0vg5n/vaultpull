package secret

import (
	"strings"
	"testing"
)

func TestRollback_NoChanges(t *testing.T) {
	dst := map[string]string{"KEY": "val"}
	snap := map[string]string{"KEY": "val"}
	_, result := Rollback(dst, snap, DefaultRollbackOptions())
	if len(result.Applied) != 0 {
		t.Fatalf("expected no applied entries, got %d", len(result.Applied))
	}
}

func TestRollback_RestoresModifiedKey(t *testing.T) {
	dst := map[string]string{"KEY": "new"}
	snap := map[string]string{"KEY": "old"}
	out, result := Rollback(dst, snap, DefaultRollbackOptions())
	if len(result.Applied) != 1 {
		t.Fatalf("expected 1 applied entry, got %d", len(result.Applied))
	}
	if out["KEY"] != "old" {
		t.Errorf("expected out[KEY]=old, got %q", out["KEY"])
	}
}

func TestRollback_RestoresMissingKey(t *testing.T) {
	dst := map[string]string{}
	snap := map[string]string{"MISSING": "restored"}
	out, result := Rollback(dst, snap, DefaultRollbackOptions())
	if len(result.Applied) != 1 {
		t.Fatalf("expected 1 applied entry, got %d", len(result.Applied))
	}
	if out["MISSING"] != "restored" {
		t.Errorf("expected restored value, got %q", out["MISSING"])
	}
}

func TestRollback_DryRun_DoesNotMutate(t *testing.T) {
	dst := map[string]string{"KEY": "current"}
	snap := map[string]string{"KEY": "original"}
	opts := DefaultRollbackOptions()
	opts.DryRun = true
	out, result := Rollback(dst, snap, opts)
	if !result.DryRun {
		t.Error("expected DryRun=true in result")
	}
	if out["KEY"] != "current" {
		t.Errorf("dry-run should not mutate dst, got %q", out["KEY"])
	}
	if len(result.Applied) != 1 {
		t.Fatalf("expected 1 applied entry in dry-run, got %d", len(result.Applied))
	}
}

func TestRollback_IgnoreKeys(t *testing.T) {
	dst := map[string]string{"KEY": "current", "SKIP": "current"}
	snap := map[string]string{"KEY": "old", "SKIP": "old"}
	opts := DefaultRollbackOptions()
	opts.IgnoreKeys = []string{"SKIP"}
	_, result := Rollback(dst, snap, opts)
	for _, e := range result.Applied {
		if e.Key == "SKIP" {
			t.Error("SKIP should not appear in applied")
		}
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "SKIP" {
		t.Errorf("expected SKIP in skipped, got %v", result.Skipped)
	}
}

func TestFormatRollbackReport_DryRun(t *testing.T) {
	r := RollbackResult{
		DryRun:  true,
		Applied: []RollbackEntry{{Key: "FOO", OldValue: "a", NewValue: "b"}},
	}
	out := FormatRollbackReport(r)
	if !strings.Contains(out, "dry-run") {
		t.Error("expected dry-run label in report")
	}
	if !strings.Contains(out, "FOO") {
		t.Error("expected key FOO in report")
	}
}

func TestFormatRollbackReport_NoChanges(t *testing.T) {
	r := RollbackResult{}
	out := FormatRollbackReport(r)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes' message, got %q", out)
	}
}
