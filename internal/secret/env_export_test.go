package secret

import (
	"strings"
	"testing"
)

func TestExport_BasicOutput(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := DefaultExportOptions()
	results := Export(src, opts)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// sorted: BAZ first
	if results[0].Line != "BAZ=qux" {
		t.Errorf("unexpected line: %s", results[0].Line)
	}
	if results[1].Line != "FOO=bar" {
		t.Errorf("unexpected line: %s", results[1].Line)
	}
}

func TestExport_WithPrefix(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	opts := DefaultExportOptions()
	opts.Prefix = "APP_"
	results := Export(src, opts)
	if results[0].Key != "APP_KEY" {
		t.Errorf("expected APP_KEY, got %s", results[0].Key)
	}
	if results[0].Line != "APP_KEY=val" {
		t.Errorf("unexpected line: %s", results[0].Line)
	}
}

func TestExport_ExportDecl(t *testing.T) {
	src := map[string]string{"TOKEN": "abc"}
	opts := DefaultExportOptions()
	opts.ExportDecl = true
	results := Export(src, opts)
	if results[0].Line != "export TOKEN=abc" {
		t.Errorf("unexpected line: %s", results[0].Line)
	}
}

func TestExport_QuoteValues(t *testing.T) {
	src := map[string]string{"MSG": `hello "world"`}
	opts := DefaultExportOptions()
	opts.QuoteValues = true
	results := Export(src, opts)
	expected := `MSG="hello \"world\""`
	if results[0].Line != expected {
		t.Errorf("expected %s, got %s", expected, results[0].Line)
	}
}

func TestExport_SkipEmpty(t *testing.T) {
	src := map[string]string{"A": "val", "B": ""}
	opts := DefaultExportOptions()
	opts.SkipEmpty = true
	results := Export(src, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "A" {
		t.Errorf("expected key A, got %s", results[0].Key)
	}
}

func TestFormatExport_JoinsLines(t *testing.T) {
	src := map[string]string{"X": "1", "Y": "2"}
	opts := DefaultExportOptions()
	results := Export(src, opts)
	out := FormatExport(results)
	if !strings.Contains(out, "\n") {
		t.Errorf("expected newline-joined output, got: %s", out)
	}
	if !strings.Contains(out, "X=1") || !strings.Contains(out, "Y=2") {
		t.Errorf("missing expected lines in: %s", out)
	}
}

func TestExport_EmptyMap(t *testing.T) {
	results := Export(map[string]string{}, DefaultExportOptions())
	if len(results) != 0 {
		t.Errorf("expected empty results")
	}
}
