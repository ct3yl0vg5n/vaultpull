package secret

import (
	"strings"
	"testing"
)

func TestClone_BasicCopy(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	dst := map[string]string{}
	opts := DefaultCloneOptions()

	r := Clone(src, dst, opts)

	if len(r.Copied) != 2 {
		t.Fatalf("expected 2 copied, got %d", len(r.Copied))
	}
	if r.Dest["FOO"] != "bar" || r.Dest["BAZ"] != "qux" {
		t.Error("expected cloned values to appear in Dest")
	}
}

func TestClone_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	dst := map[string]string{"FOO": "old"}
	opts := DefaultCloneOptions()

	r := Clone(src, dst, opts)

	if len(r.Skipped) != 1 || r.Skipped[0] != "FOO" {
		t.Fatalf("expected FOO to be skipped, got %+v", r.Skipped)
	}
	if r.Dest["FOO"] != "old" {
		t.Error("expected original value to be preserved")
	}
}

func TestClone_OverwriteExisting(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	dst := map[string]string{"FOO": "old"}
	opts := DefaultCloneOptions()
	opts.Overwrite = true

	r := Clone(src, dst, opts)

	if len(r.Copied) != 1 {
		t.Fatalf("expected 1 copied, got %d", len(r.Copied))
	}
	if r.Dest["FOO"] != "new" {
		t.Error("expected overwritten value")
	}
}

func TestClone_WithPrefix(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	dst := map[string]string{}
	opts := DefaultCloneOptions()
	opts.KeyPrefix = "PROD_"

	r := Clone(src, dst, opts)

	if _, ok := r.Dest["PROD_KEY"]; !ok {
		t.Error("expected prefixed key PROD_KEY in Dest")
	}
	if len(r.Copied) != 1 || r.Copied[0] != "PROD_KEY" {
		t.Errorf("unexpected Copied list: %v", r.Copied)
	}
}

func TestClone_DryRun_DoesNotMutate(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	dst := map[string]string{}
	opts := DefaultCloneOptions()
	opts.DryRun = true

	r := Clone(src, dst, opts)

	if len(r.Copied) != 1 {
		t.Fatal("expected 1 in Copied for dry run")
	}
	if _, ok := r.Dest["FOO"]; ok {
		t.Error("dry run should not write to Dest")
	}
}

func TestClone_IgnoreKeys(t *testing.T) {
	src := map[string]string{"FOO": "1", "BAR": "2"}
	dst := map[string]string{}
	opts := DefaultCloneOptions()
	opts.IgnoreKeys = []string{"FOO"}

	r := Clone(src, dst, opts)

	if _, ok := r.Dest["FOO"]; ok {
		t.Error("ignored key FOO should not appear in Dest")
	}
	if r.Dest["BAR"] != "2" {
		t.Error("expected BAR to be cloned")
	}
}

func TestFormatCloneReport_ContainsSummary(t *testing.T) {
	r := CloneResult{
		Copied:  []string{"FOO"},
		Skipped: []string{"BAR"},
		Dest:    map[string]string{"FOO": "1"},
	}
	out := FormatCloneReport(r)
	if !strings.Contains(out, "Cloned: 1") || !strings.Contains(out, "Skipped: 1") {
		t.Errorf("unexpected report output: %s", out)
	}
}
