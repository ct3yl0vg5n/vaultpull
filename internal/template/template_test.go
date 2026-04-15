package template

import (
	"testing"
)

func TestApply_NoRule(t *testing.T) {
	key, err := Apply("DATABASE_URL", Rule{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %q", key)
	}
}

func TestApply_Prefix(t *testing.T) {
	key, err := Apply("HOST", Rule{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %q", key)
	}
}

func TestApply_Upper(t *testing.T) {
	key, err := Apply("db_password", Rule{Upper: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD, got %q", key)
	}
}

func TestApply_Replace(t *testing.T) {
	key, err := Apply("db-host", Rule{
		Replace: [][2]string{{"\u002d", "_"}},
		Upper:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", key)
	}
}

func TestApply_InvalidResult(t *testing.T) {
	_, err := Apply("123bad", Rule{})
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

func TestApplyMap_TransformsKeys(t *testing.T) {
	secrets := map[string]string{
		"host":     "localhost",
		"password": "secret",
	}
	out := ApplyMap(secrets, Rule{Prefix: "DB_", Upper: true})

	if v, ok := out["DB_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", out)
	}
	if v, ok := out["DB_PASSWORD"]; !ok || v != "secret" {
		t.Errorf("expected DB_PASSWORD=secret, got %v", out)
	}
}

func TestApplyMap_SkipsInvalidKeys(t *testing.T) {
	secrets := map[string]string{
		"valid_key": "value",
		"123bad":    "skip",
	}
	out := ApplyMap(secrets, Rule{})

	if _, ok := out["123bad"]; ok {
		t.Error("expected invalid key to be skipped")
	}
	if _, ok := out["valid_key"]; !ok {
		t.Error("expected valid_key to be present")
	}
}
