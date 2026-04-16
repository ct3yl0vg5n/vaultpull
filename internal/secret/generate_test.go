package secret

import (
	"strings"
	"testing"
)

func TestGenerate_DefaultLength(t *testing.T) {
	opts := DefaultGenerateOptions()
	v, err := Generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v) != 32 {
		t.Errorf("expected length 32, got %d", len(v))
	}
}

func TestGenerate_CustomLength(t *testing.T) {
	opts := DefaultGenerateOptions()
	opts.Length = 16
	v, err := Generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v) != 16 {
		t.Errorf("expected length 16, got %d", len(v))
	}
}

func TestGenerate_ZeroLength_Error(t *testing.T) {
	opts := DefaultGenerateOptions()
	opts.Length = 0
	_, err := Generate(opts)
	if err == nil {
		t.Fatal("expected error for zero length")
	}
}

func TestGenerate_Base64(t *testing.T) {
	opts := DefaultGenerateOptions()
	opts.Base64 = true
	opts.Length = 24
	v, err := Generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v) != 24 {
		t.Errorf("expected length 24, got %d", len(v))
	}
}

func TestGenerate_OnlyAlphanumeric(t *testing.T) {
	opts := DefaultGenerateOptions()
	opts.Length = 64
	v, err := Generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, c := range v {
		if !strings.ContainsRune(Alphanumeric, c) {
			t.Errorf("unexpected char %q in output", c)
		}
	}
}

func TestGenerate_Uniqueness(t *testing.T) {
	opts := DefaultGenerateOptions()
	a, _ := Generate(opts)
	b, _ := Generate(opts)
	if a == b {
		t.Error("two generated secrets should not be equal")
	}
}
