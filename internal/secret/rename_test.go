package secret

import (
	"strings"
	"testing"
)

func TestRename_BasicRule(t *testing.T) {
	src := map[string]string{"OLD_KEY": "value1", "OTHER": "value2"}
	opts := DefaultRenameOptions()
	opts.Rules = map[string]string{"OLD_KEY": "NEW_KEY"}

	dst, results, err := Rename(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := dst["OLD_KEY"]; ok {
		t.Error("old key should have been removed")
	}
	if dst["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", dst["NEW_KEY"])
	}
	if len(results) != 1 || !results[0].Applied {
		t.Error("expected one applied result")
	}
}

func TestRename_KeyNotFound_NoError(t *testing.T) {
	src := map[string]string{"EXISTING": "v"}
	opts := DefaultRenameOptions()
	opts.Rules = map[string]string{"MISSING": "NEW"}

	_, results, err := Rename(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Error("expected one skipped result")
	}
}

func TestRename_KeyNotFound_ErrorOnMissing(t *testing.T) {
	src := map[string]string{"EXISTING": "v"}
	opts := DefaultRenameOptions()
	opts.ErrorOnMissing = true
	opts.Rules = map[string]string{"MISSING": "NEW"}

	_, _, err := Rename(src, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "MISSING") {
		t.Errorf("error should mention key name, got: %v", err)
	}
}

func TestRename_DryRun_DoesNotMutate(t *testing.T) {
	src := map[string]string{"OLD": "val"}
	opts := DefaultRenameOptions()
	opts.DryRun = true
	opts.Rules = map[string]string{"OLD": "NEW"}

	dst, results, err := Rename(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := dst["OLD"]; !ok {
		t.Error("dry-run should not remove old key")
	}
	if _, ok := dst["NEW"]; ok {
		t.Error("dry-run should not add new key")
	}
	if len(results) != 1 || results[0].Applied {
		t.Error("dry-run result should not be marked applied")
	}
}

func TestRename_OriginalUnchanged(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	opts := DefaultRenameOptions()
	opts.Rules = map[string]string{"A": "Z"}

	dst, _, _ := Rename(src, opts)
	if src["A"] != "1" {
		t.Error("src map should not be mutated")
	}
	if dst["B"] != "2" {
		t.Error("unrelated keys should be preserved in dst")
	}
}

func TestFormatRenameReport_NoResults(t *testing.T) {
	out := FormatRenameReport(nil)
	if out != "no rename rules applied" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatRenameReport_Summary(t *testing.T) {
	results := []RenameResult{
		{OldKey: "A", NewKey: "B", Applied: true},
		{OldKey: "C", NewKey: "D", Skipped: true},
	}
	out := FormatRenameReport(results)
	if !strings.Contains(out, "1 renamed") {
		t.Errorf("expected renamed count in output: %q", out)
	}
	if !strings.Contains(out, "1 skipped") {
		t.Errorf("expected skipped count in output: %q", out)
	}
}
