package secret_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/secret"
)

func TestApplyKey_Uppercase(t *testing.T) {
	out, err := secret.ApplyKey("db_host", secret.Rule{UppercaseKeys: true})
	if err != nil {
		t.Fatal(err)
	}
	if out != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", out)
	}
}

func TestApplyKey_Prefix(t *testing.T) {
	out, err := secret.ApplyKey("HOST", secret.Rule{KeyPrefix: "APP_"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %s", out)
	}
}

func TestApplyKey_Replace(t *testing.T) {
	out, err := secret.ApplyKey("DB-HOST", secret.Rule{
		UppercaseKeys: true,
		KeyReplace:    &secret.KeyReplace{From: "-", To: "_"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", out)
	}
}

func TestApplyKey_InvalidChars(t *testing.T) {
	_, err := secret.ApplyKey("db-host", secret.Rule{})
	if err == nil {
		t.Error("expected error for invalid key characters")
	}
}

func TestApplyValue_Redact(t *testing.T) {
	out := secret.ApplyValue("supersecret", secret.Rule{RedactValues: true})
	if out != "***REDACTED***" {
		t.Errorf("expected redacted value, got %s", out)
	}
}

func TestApplyValue_NoRedact(t *testing.T) {
	out := secret.ApplyValue("supersecret", secret.Rule{})
	if out != "supersecret" {
		t.Errorf("expected plain value, got %s", out)
	}
}

func TestApplyMap_Success(t *testing.T) {
	input := map[string]string{"db_host": "localhost", "db_port": "5432"}
	out, err := secret.ApplyMap(input, secret.Rule{UppercaseKeys: true})
	if err != nil {
		t.Fatal(err)
	}
	if out["DB_HOST"] != "localhost" || out["DB_PORT"] != "5432" {
		t.Errorf("unexpected map result: %v", out)
	}
}

func TestApplyMap_Error(t *testing.T) {
	input := map[string]string{"bad-key": "value"}
	_, err := secret.ApplyMap(input, secret.Rule{})
	if err == nil {
		t.Error("expected error for invalid key in map")
	}
}
