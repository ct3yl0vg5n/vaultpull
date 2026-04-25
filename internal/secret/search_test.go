package secret

import (
	"strings"
	"testing"
)

func TestSearch_NoResults(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	results, err := Search(secrets, "REDIS", DefaultSearchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearch_KeyMatch_CaseInsensitive(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"}
	results, err := Search(secrets, "db", DefaultSearchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %v", results)
	}
	if !results[0].MatchedKey {
		t.Error("expected MatchedKey to be true")
	}
}

func TestSearch_CaseSensitive_NoMatch(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	opts := DefaultSearchOptions()
	opts.CaseSensitive = true
	results, err := Search(secrets, "db", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results with case-sensitive search, got %d", len(results))
	}
}

func TestSearch_ValueMatch(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "prod-server", "APP_NAME": "myapp"}
	opts := DefaultSearchOptions()
	opts.SearchValues = true
	results, err := Search(secrets, "prod", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST matched by value, got %v", results)
	}
	if !results[0].MatchedVal {
		t.Error("expected MatchedVal to be true")
	}
}

func TestSearch_RegexMatch(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "API_KEY": "abc"}
	opts := DefaultSearchOptions()
	opts.Regex = true
	results, err := Search(secrets, "^DB_", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_InvalidRegex_ReturnsError(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	opts := DefaultSearchOptions()
	opts.Regex = true
	_, err := Search(secrets, "[", opts)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestSearch_ResultsSortedByKey(t *testing.T) {
	secrets := map[string]string{"Z_KEY": "v", "A_KEY": "v", "M_KEY": "v"}
	results, err := Search(secrets, "KEY", DefaultSearchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results")
	}
	if results[0].Key != "A_KEY" || results[1].Key != "M_KEY" || results[2].Key != "Z_KEY" {
		t.Errorf("results not sorted: %v", results)
	}
}

func TestFormatSearchReport_NoResults(t *testing.T) {
	out := FormatSearchReport(nil, "MISSING")
	if !strings.Contains(out, "no matches") {
		t.Errorf("expected no-match message, got: %s", out)
	}
}

func TestFormatSearchReport_WithResults(t *testing.T) {
	results := []SearchResult{
		{Key: "DB_HOST", Value: "localhost", MatchedKey: true},
	}
	out := FormatSearchReport(results, "DB")
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "[key]") {
		t.Errorf("expected [key] tag in output, got: %s", out)
	}
}
