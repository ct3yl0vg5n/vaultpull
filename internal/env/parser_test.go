package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_Basic(t *testing.T) {
	content := `# comment
DB_HOST=localhost
DB_PORT=5432
API_KEY="supersecret"
`
	path := writeTempFile(t, content)

	got, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "supersecret",
	}
	for k, v := range expect {
		if got[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, got[k])
		}
	}
}

func TestParseFile_NotExist(t *testing.T) {
	got, err := ParseFile("/nonexistent/.env")
	if err != nil {
		t.Fatalf("expected empty map for missing file, got error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestParseFile_InvalidSyntax(t *testing.T) {
	content := "BADLINE\n"
	path := writeTempFile(t, content)
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for invalid syntax, got nil")
	}
}

func TestWriteFile_RoundTrip(t *testing.T) {
	data := map[string]string{
		"Z_KEY": "zval",
		"A_KEY": "aval",
		"M_KEY": "mval",
	}
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := WriteFile(path, data); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	parsed, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	for k, v := range data {
		if parsed[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, parsed[k])
		}
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}
