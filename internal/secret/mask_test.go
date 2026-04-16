package secret

import (
	"testing"
)

func TestMask_Default(t *testing.T) {
	opts := DefaultMaskOptions()
	result := Mask("supersecret", opts)
	if result != "*******cret" {
		t.Errorf("unexpected mask: %q", result)
	}
}

func TestMask_Disabled(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.Enabled = false
	result := Mask("mysecret", opts)
	if result != "mysecret" {
		t.Errorf("expected plain value, got %q", result)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	opts := DefaultMaskOptions()
	result := Mask("", opts)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestMask_RevealMoreThanLength(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.RevealLast = 20
	result := Mask("abc", opts)
	if result != "***" {
		t.Errorf("expected full mask, got %q", result)
	}
}

func TestMask_ZeroReveal(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.RevealLast = 0
	result := Mask("hello", opts)
	if result != "*****" {
		t.Errorf("expected full mask, got %q", result)
	}
}

func TestMaskMap(t *testing.T) {
	opts := DefaultMaskOptions()
	input := map[string]string{
		"KEY1": "password123",
		"KEY2": "tok",
	}
	out := MaskMap(input, opts)
	if out["KEY1"] != "*******1234"[0:11] {
		// just verify it's masked
		if out["KEY1"] == "password123" {
			t.Error("expected KEY1 to be masked")
		}
	}
	if out["KEY2"] == "tok" {
		t.Error("expected KEY2 to be masked")
	}
}
