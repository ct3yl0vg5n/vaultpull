package secret

import (
	"testing"
)

func TestImport_BasicEnv(t *testing.T) {
	dst := map[string]string{}
	lines := []string{"FOO=bar", "BAZ=qux"}
	r, err := Import(dst, lines, DefaultImportOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Imported != 2 || r.Skipped != 0 {
		t.Errorf("expected 2 imported, got %+v", r)
	}
	if dst["FOO"] != "bar" || dst["BAZ"] != "qux" {
		t.Errorf("unexpected dst: %v", dst)
	}
}

func TestImport_StripExportPrefix(t *testing.T) {
	dst := map[string]string{}
	lines := []string{`export DB_URL="postgres://localhost"`}
	opts := DefaultImportOptions()
	opts.StripExport = true
	r, _ := Import(dst, lines, opts)
	if r.Imported != 1 {
		t.Errorf("expected 1 imported, got %d", r.Imported)
	}
	if dst["DB_URL"] != "postgres://localhost" {
		t.Errorf("unexpected value: %q", dst["DB_URL"])
	}
}

func TestImport_SkipsExistingWithoutOverwrite(t *testing.T) {
	dst := map[string]string{"FOO": "original"}
	lines := []string{"FOO=new"}
	opts := DefaultImportOptions()
	opts.Overwrite = false
	r, _ := Import(dst, lines, opts)
	if r.Skipped != 1 || r.Imported != 0 {
		t.Errorf("expected 1 skipped, got %+v", r)
	}
	if dst["FOO"] != "original" {
		t.Errorf("value should not have changed")
	}
}

func TestImport_OverwriteExisting(t *testing.T) {
	dst := map[string]string{"FOO": "original"}
	lines := []string{"FOO=new"}
	opts := DefaultImportOptions()
	opts.Overwrite = true
	r, _ := Import(dst, lines, opts)
	if r.Imported != 1 {
		t.Errorf("expected 1 imported, got %+v", r)
	}
	if dst["FOO"] != "new" {
		t.Errorf("expected overwritten value, got %q", dst["FOO"])
	}
}

func TestImport_IgnoreKeys(t *testing.T) {
	dst := map[string]string{}
	lines := []string{"FOO=bar", "SECRET=hidden"}
	opts := DefaultImportOptions()
	opts.IgnoreKeys = []string{"SECRET"}
	r, _ := Import(dst, lines, opts)
	if r.Imported != 1 || r.Skipped != 1 {
		t.Errorf("expected 1 imported 1 skipped, got %+v", r)
	}
	if _, ok := dst["SECRET"]; ok {
		t.Error("SECRET should have been ignored")
	}
}

func TestImport_InvalidLineReturnsError(t *testing.T) {
	dst := map[string]string{}
	lines := []string{"NOTVALID", "GOOD=value"}
	r, _ := Import(dst, lines, DefaultImportOptions())
	if len(r.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(r.Errors))
	}
	if r.Imported != 1 {
		t.Errorf("expected 1 imported despite error, got %d", r.Imported)
	}
}

func TestImport_CommentsAndBlankLinesSkipped(t *testing.T) {
	dst := map[string]string{}
	lines := []string{"", "# comment", "KEY=val"}
	r, _ := Import(dst, lines, DefaultImportOptions())
	if r.Imported != 1 {
		t.Errorf("expected 1 imported, got %d", r.Imported)
	}
}

func TestFormatImportReport_WithErrors(t *testing.T) {
	r := ImportResult{Imported: 3, Skipped: 1, Errors: []string{"line 2: invalid format"}}
	out := FormatImportReport(r)
	if out == "" {
		t.Error("expected non-empty report")
	}
	if len(out) < 10 {
		t.Errorf("report too short: %q", out)
	}
}
