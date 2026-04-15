package filter_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/filter"
)

var sampleSecrets = map[string]string{
	"APP_DB_HOST":  "localhost",
	"APP_DB_PASS":  "secret",
	"APP_API_KEY":  "key123",
	"OTHER_SECRET": "ignore",
	"PLAIN":        "value",
}

func TestApply_NoRule(t *testing.T) {
	r := &filter.Rule{}
	out := r.Apply(sampleSecrets)
	if len(out) != len(sampleSecrets) {
		t.Errorf("expected %d keys, got %d", len(sampleSecrets), len(out))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	r := &filter.Rule{Prefix: "APP_"}
	out := r.Apply(sampleSecrets)
	if len(out) != 3 {
		t.Errorf("expected 3 keys with APP_ prefix, got %d", len(out))
	}
	if _, ok := out["OTHER_SECRET"]; ok {
		t.Error("OTHER_SECRET should have been filtered out")
	}
}

func TestApply_StripPrefix(t *testing.T) {
	r := &filter.Rule{Prefix: "APP_", StripPrefix: true}
	out := r.Apply(sampleSecrets)
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected key DB_HOST after stripping APP_ prefix")
	}
	if _, ok := out["APP_DB_HOST"]; ok {
		t.Error("original key APP_DB_HOST should not be present")
	}
}

func TestApply_Exclude(t *testing.T) {
	r := &filter.Rule{Exclude: []string{"PLAIN", "OTHER_SECRET"}}
	out := r.Apply(sampleSecrets)
	if _, ok := out["PLAIN"]; ok {
		t.Error("PLAIN should have been excluded")
	}
	if _, ok := out["OTHER_SECRET"]; ok {
		t.Error("OTHER_SECRET should have been excluded")
	}
}

func TestApply_PrefixAndExclude(t *testing.T) {
	r := &filter.Rule{Prefix: "APP_", Exclude: []string{"APP_API_KEY"}}
	out := r.Apply(sampleSecrets)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_API_KEY"]; ok {
		t.Error("APP_API_KEY should have been excluded")
	}
}

func TestIsEmpty(t *testing.T) {
	if !(&filter.Rule{}).IsEmpty() {
		t.Error("empty rule should return true for IsEmpty")
	}
	if (&filter.Rule{Prefix: "X_"}).IsEmpty() {
		t.Error("rule with prefix should not be empty")
	}
	if (&filter.Rule{Exclude: []string{"FOO"}}).IsEmpty() {
		t.Error("rule with exclusions should not be empty")
	}
}
