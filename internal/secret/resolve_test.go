package secret

import (
	"strings"
	"testing"
)

func TestResolve_NoRefs(t *testing.T) {
	src := map[string]string{"A": "hello", "B": "world"}
	out, results, err := Resolve(src, DefaultResolveOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "hello" || out["B"] != "world" {
		t.Errorf("unexpected output: %v", out)
	}
	for _, r := range results {
		if r.WasRef {
			t.Errorf("expected no refs, got WasRef=true for %s", r.Key)
		}
	}
}

func TestResolve_BasicRef(t *testing.T) {
	src := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "ref:BASE_URL",
	}
	out, _, err := Resolve(src, DefaultResolveOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_URL"] != "https://example.com" {
		t.Errorf("expected resolved value, got %q", out["API_URL"])
	}
}

func TestResolve_MissingRef_Strict(t *testing.T) {
	src := map[string]string{"A": "ref:MISSING"}
	opts := DefaultResolveOptions()
	opts.Strict = true
	_, _, err := Resolve(src, opts)
	if err == nil {
		t.Fatal("expected error for missing ref in strict mode")
	}
}

func TestResolve_MissingRef_NonStrict(t *testing.T) {
	src := map[string]string{"A": "ref:MISSING"}
	opts := DefaultResolveOptions()
	opts.Strict = false
	out, results, err := Resolve(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "ref:MISSING" {
		t.Errorf("expected original value preserved, got %q", out["A"])
	}
	for _, r := range results {
		if r.Key == "A" && !r.Missing {
			t.Error("expected Missing=true")
		}
	}
}

func TestResolve_CustomPrefix(t *testing.T) {
	src := map[string]string{
		"TOKEN": "secret123",
		"AUTH":  "$ref:TOKEN",
	}
	opts := DefaultResolveOptions()
	opts.RefPrefix = "$ref:"
	out, _, err := Resolve(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["AUTH"] != "secret123" {
		t.Errorf("expected resolved value, got %q", out["AUTH"])
	}
}

func TestFormatResolveReport_ContainsSummary(t *testing.T) {
	results := []ResolveResult{
		{Key: "A", Original: "ref:B", Resolved: "val", WasRef: true},
		{Key: "C", Original: "ref:X", Resolved: "ref:X", WasRef: true, Missing: true},
	}
	report := FormatResolveReport(results)
	if !strings.Contains(report, "2 reference(s)") {
		t.Errorf("expected ref count in report, got: %s", report)
	}
	if !strings.Contains(report, "1 missing") {
		t.Errorf("expected missing count in report, got: %s", report)
	}
}
