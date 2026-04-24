package secret

import (
	"strings"
	"testing"
)

func TestPromote_AddNew(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	dst := map[string]string{}
	opts := DefaultPromoteOptions()

	out, results := Promote(src, dst, opts)

	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected keys to be promoted, got %v", out)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Action != "added" {
			t.Errorf("expected action 'added', got %q for key %s", r.Action, r.Key)
		}
	}
}

func TestPromote_SkipsExistingByDefault(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	dst := map[string]string{"FOO": "old"}
	opts := DefaultPromoteOptions()

	out, results := Promote(src, dst, opts)

	if out["FOO"] != "old" {
		t.Errorf("expected original value preserved, got %q", out["FOO"])
	}
	if len(results) != 1 || results[0].Action != "skipped" {
		t.Errorf("expected skipped result, got %+v", results)
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	dst := map[string]string{"FOO": "old"}
	opts := DefaultPromoteOptions()
	opts.OverwriteExisting = true

	out, results := Promote(src, dst, opts)

	if out["FOO"] != "new" {
		t.Errorf("expected overwritten value, got %q", out["FOO"])
	}
	if results[0].Action != "overwritten" {
		t.Errorf("expected 'overwritten', got %q", results[0].Action)
	}
}

func TestPromote_DryRun_DoesNotMutate(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	dst := map[string]string{}
	opts := DefaultPromoteOptions()
	opts.DryRun = true

	out, results := Promote(src, dst, opts)

	if _, ok := out["FOO"]; ok {
		t.Error("dry run should not write to destination")
	}
	if results[0].Action != "added" {
		t.Errorf("expected action 'added' in dry run, got %q", results[0].Action)
	}
}

func TestPromote_IgnoreKeys(t *testing.T) {
	src := map[string]string{"FOO": "bar", "SECRET": "hidden"}
	dst := map[string]string{}
	opts := DefaultPromoteOptions()
	opts.IgnoreKeys = []string{"SECRET"}

	out, results := Promote(src, dst, opts)

	if _, ok := out["SECRET"]; ok {
		t.Error("ignored key should not appear in output")
	}
	actions := map[string]string{}
	for _, r := range results {
		actions[r.Key] = r.Action
	}
	if actions["SECRET"] != "ignored" {
		t.Errorf("expected SECRET to be ignored, got %q", actions["SECRET"])
	}
	if actions["FOO"] != "added" {
		t.Errorf("expected FOO to be added, got %q", actions["FOO"])
	}
}

func TestFormatPromoteReport_ContainsSummary(t *testing.T) {
	results := []PromoteResult{
		{Key: "A", Action: "added"},
		{Key: "B", Action: "skipped"},
		{Key: "C", Action: "ignored"},
	}
	report := FormatPromoteReport("staging", "production", results)
	if !strings.Contains(report, "staging → production") {
		t.Error("report missing environment labels")
	}
	if !strings.Contains(report, "1 added") {
		t.Error("report missing added count")
	}
	if !strings.Contains(report, "1 skipped") {
		t.Error("report missing skipped count")
	}
}
