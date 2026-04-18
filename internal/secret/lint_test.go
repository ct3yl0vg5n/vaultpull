package secret

import (
	"testing"
)

func TestLintMap_NoIssues(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc123"}
	results := LintMap(secrets, DefaultLintOptions())
	for _, r := range results {
		t.Errorf("unexpected lint result: %s: %s", r.Key, r.Message)
	}
}

func TestLintMap_DisallowedPrefix(t *testing.T) {
	secrets := map[string]string{"TMP_TOKEN": "xyz"}
	results := LintMap(secrets, DefaultLintOptions())
	if !hasResult(results, "TMP_TOKEN", "prefix") {
		t.Error("expected disallowed prefix warning")
	}
}

func TestLintMap_DisallowedSuffix(t *testing.T) {
	secrets := map[string]string{"SECRET_OLD": "val"}
	results := LintMap(secrets, DefaultLintOptions())
	if !hasResult(results, "SECRET_OLD", "suffix") {
		t.Error("expected disallowed suffix warning")
	}
}

func TestLintMap_LowercaseKey(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost"}
	results := LintMap(secrets, DefaultLintOptions())
	if !hasResult(results, "db_host", "lowercase") {
		t.Error("expected lowercase key warning")
	}
}

func TestLintMap_WhitespaceValue(t *testing.T) {
	secrets := map[string]string{"API_KEY": "hello world"}
	results := LintMap(secrets, DefaultLintOptions())
	if !hasResult(results, "API_KEY", "whitespace") {
		t.Error("expected whitespace value warning")
	}
}

func TestLintMap_WarnOnLowercase_Disabled(t *testing.T) {
	opts := DefaultLintOptions()
	opts.WarnOnLowercase = false
	secrets := map[string]string{"db_host": "val"}
	results := LintMap(secrets, opts)
	if hasResult(results, "db_host", "lowercase") {
		t.Error("did not expect lowercase warning when disabled")
	}
}

func hasResult(results []LintResult, key, msgSubstr string) bool {
	for _, r := range results {
		if r.Key == key {
			for i := range r.Message {
				if i+len(msgSubstr) <= len(r.Message) && r.Message[i:i+len(msgSubstr)] == msgSubstr {
					return true
				}
			}
		}
	}
	return false
}
