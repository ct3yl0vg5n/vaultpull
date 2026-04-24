package secret

import (
	"strings"
	"testing"
)

func TestEnforcePolicy_NoViolations(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "xK9#mP2$qR",
	}
	opts := DefaultPolicyOptions()
	violations := EnforcePolicy(secrets, opts)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestEnforcePolicy_DenyPatternTriggered(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "password",
	}
	opts := DefaultPolicyOptions()
	violations := EnforcePolicy(secrets, opts)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DB_PASSWORD" {
		t.Errorf("expected key DB_PASSWORD, got %s", violations[0].Key)
	}
}

func TestEnforcePolicy_RequiredKeyMissing(t *testing.T) {
	secrets := map[string]string{}
	opts := PolicyOptions{
		Rules: []PolicyRule{
			{Name: "API_KEY", Required: true},
		},
	}
	violations := EnforcePolicy(secrets, opts)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Msg, "missing") {
		t.Errorf("expected missing message, got: %s", violations[0].Msg)
	}
}

func TestEnforcePolicy_KeyPatternNoMatch(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "password", // key doesn't match PASSWORD pattern
	}
	opts := DefaultPolicyOptions()
	violations := EnforcePolicy(secrets, opts)
	if len(violations) != 0 {
		t.Errorf("expected no violations for non-matching key, got %d", len(violations))
	}
}

func TestEnforcePolicy_EmptySecrets(t *testing.T) {
	violations := EnforcePolicy(map[string]string{}, DefaultPolicyOptions())
	if len(violations) != 0 {
		t.Errorf("expected no violations for empty secrets, got %d", len(violations))
	}
}

func TestFormatPolicyReport_NoViolations(t *testing.T) {
	out := FormatPolicyReport(nil)
	if !strings.Contains(out, "passed") {
		t.Errorf("expected 'passed' in output, got: %s", out)
	}
}

func TestFormatPolicyReport_WithViolations(t *testing.T) {
	violations := []PolicyViolation{
		{Key: "DB_PASSWORD", Rule: "no-plaintext-password", Msg: "value matches denied pattern"},
	}
	out := FormatPolicyReport(violations)
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected key in report, got: %s", out)
	}
	if !strings.Contains(out, "violations (1)") {
		t.Errorf("expected count in report, got: %s", out)
	}
}
