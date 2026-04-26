package secret

import (
	"strings"
	"testing"
)

func TestApplyNamespace_BasicPrefix(t *testing.T) {
	src := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out, results, err := ApplyNamespace(src, NamespaceOptions{Namespace: "prod", Separator: "__", Uppercase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if _, ok := out["PROD__DB_HOST"]; !ok {
		t.Errorf("expected PROD__DB_HOST in output")
	}
	if _, ok := out["PROD__DB_PORT"]; !ok {
		t.Errorf("expected PROD__DB_PORT in output")
	}
}

func TestApplyNamespace_LowercaseNamespace(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	out, _, err := ApplyNamespace(src, NamespaceOptions{Namespace: "staging", Separator: "_", Uppercase: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["staging_KEY"]; !ok {
		t.Errorf("expected staging_KEY, got %v", out)
	}
}

func TestApplyNamespace_EmptyNamespace_Error(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	_, _, err := ApplyNamespace(src, NamespaceOptions{Namespace: ""})
	if err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestApplyNamespace_StripOnRead(t *testing.T) {
	src := map[string]string{
		"PROD__DB_HOST": "localhost",
		"PROD__DB_PORT": "5432",
		"OTHER_KEY":     "ignored",
	}
	out, results, err := ApplyNamespace(src, NamespaceOptions{
		Namespace:   "prod",
		Separator:   "__",
		Uppercase:   true,
		StripOnRead: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Errorf("expected DB_HOST after strip")
	}
	if _, ok := out["OTHER_KEY"]; ok {
		t.Errorf("OTHER_KEY should have been excluded")
	}
}

func TestApplyNamespace_EmptySrc(t *testing.T) {
	out, results, err := ApplyNamespace(map[string]string{}, NamespaceOptions{Namespace: "dev", Separator: "__", Uppercase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 || len(results) != 0 {
		t.Errorf("expected empty output for empty src")
	}
}

func TestFormatNamespaceReport_WithResults(t *testing.T) {
	results := []NamespaceResult{
		{Key: "DB_HOST", NewKey: "PROD__DB_HOST", Value: "localhost"},
	}
	report := FormatNamespaceReport(results, false)
	if !strings.Contains(report, "1 key(s) namespaced") {
		t.Errorf("unexpected report: %s", report)
	}
	if !strings.Contains(report, "DB_HOST -> PROD__DB_HOST") {
		t.Errorf("expected key mapping in report")
	}
}

func TestFormatNamespaceReport_NoResults(t *testing.T) {
	report := FormatNamespaceReport(nil, false)
	if !strings.Contains(report, "no keys affected") {
		t.Errorf("expected no-op message, got: %s", report)
	}
}
