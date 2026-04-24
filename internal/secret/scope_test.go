package secret

import (
	"strings"
	"testing"
)

func TestPartition_NoScopes_AllDefault(t *testing.T) {
	src := map[string]string{"FOO": "1", "BAR": "2"}
	opts := DefaultScopeOptions()
	results := Partition(src, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 scope, got %d", len(results))
	}
	if results[0].Name != "default" {
		t.Errorf("expected scope 'default', got %q", results[0].Name)
	}
	if len(results[0].Secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(results[0].Secrets))
	}
}

func TestPartition_PrefixRouting(t *testing.T) {
	src := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_KEY":  "secret",
		"OTHER":   "val",
	}
	opts := ScopeOptions{
		Scopes: map[string][]string{
			"database": {"DB_"},
			"app":      {"APP_"},
		},
		DefaultScope: "misc",
	}
	results := Partition(src, opts)
	scopes := map[string]ScopeResult{}
	for _, r := range results {
		scopes[r.Name] = r
	}
	if len(scopes["database"].Secrets) != 2 {
		t.Errorf("expected 2 db secrets, got %d", len(scopes["database"].Secrets))
	}
	if len(scopes["app"].Secrets) != 1 {
		t.Errorf("expected 1 app secret, got %d", len(scopes["app"].Secrets))
	}
	if len(scopes["misc"].Secrets) != 1 {
		t.Errorf("expected 1 misc secret, got %d", len(scopes["misc"].Secrets))
	}
}

func TestPartition_EmptySrc(t *testing.T) {
	opts := DefaultScopeOptions()
	results := Partition(map[string]string{}, opts)
	if len(results) != 0 {
		t.Errorf("expected 0 scopes for empty src, got %d", len(results))
	}
}

func TestPartition_FirstMatchWins(t *testing.T) {
	src := map[string]string{"APP_DB_HOST": "localhost"}
	opts := ScopeOptions{
		Scopes: map[string][]string{
			"app": {"APP_"},
			"db":  {"APP_DB_"},
		},
		DefaultScope: "default",
	}
	results := Partition(src, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 scope, got %d", len(results))
	}
	// alphabetically 'app' < 'db', so APP_DB_HOST goes to 'app'
	if results[0].Name != "app" {
		t.Errorf("expected scope 'app', got %q", results[0].Name)
	}
}

func TestFormatScopeReport_NoResults(t *testing.T) {
	out := FormatScopeReport(nil)
	if !strings.Contains(out, "no secrets") {
		t.Errorf("expected 'no secrets' message, got %q", out)
	}
}

func TestFormatScopeReport_WithData(t *testing.T) {
	results := []ScopeResult{
		{Name: "db", Secrets: map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}},
	}
	out := FormatScopeReport(results)
	if !strings.Contains(out, "[db] 2 key(s)") {
		t.Errorf("missing scope header in output: %q", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("missing key DB_HOST in output: %q", out)
	}
}
