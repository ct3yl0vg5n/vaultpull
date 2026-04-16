package secret

import (
	"strings"
	"testing"
)

func TestRedact_DefaultOptions(t *testing.T) {
	opts := DefaultRedactOptions()
	result := Redact("supersecret", opts)
	if !strings.HasSuffix(result, "cret") {
		t.Errorf("expected suffix 'cret', got %q", result)
	}
	if strings.Contains(result, "super") {
		t.Errorf("expected prefix to be masked, got %q", result)
	}
}

func TestRedact_FullMask(t *testing.T) {
	opts := RedactOptions{Enabled: true, MaskChar: "*", RevealChars: 0}
	result := Redact("hello", opts)
	if result != "*****" {
		t.Errorf("expected '****', got %q", result)
	}
}

func TestRedact_Disabled(t *testing.T) {
	opts := RedactOptions{Enabled: false}
	result := Redact("plaintext", opts)
	if result != "plaintext" {
		t.Errorf("expected original value, got %q", result)
	}
}

func TestRedact_EmptyValue(t *testing.T) {
	opts := DefaultRedactOptions()
	result := Redact("", opts)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestRedact_RevealMoreThanLength(t *testing.T) {
	opts := RedactOptions{Enabled: true, MaskChar: "#", RevealChars: 20}
	result := Redact("short", opts)
	if result != "#####" {
		t.Errorf("expected full mask when reveal >= len, got %q", result)
	}
}

func TestRedactMap(t *testing.T) {
	opts := RedactOptions{Enabled: true, MaskChar: "*", RevealChars: 0}
	input := map[string]string{"KEY1": "value1", "KEY2": "value2"}
	out := RedactMap(input, opts)
	for k, v := range out {
		if strings.Contains(v, "value") {
			t.Errorf("key %s: expected redacted value, got %q", k, v)
		}
	}
	if len(out) != len(input) {
		t.Errorf("expected same number of keys")
	}
}

func TestRedactMap_Disabled(t *testing.T) {
	opts := RedactOptions{Enabled: false}
	input := map[string]string{"A": "secret"}
	out := RedactMap(input, opts)
	if out["A"] != "secret" {
		t.Errorf("expected original value when disabled")
	}
}
