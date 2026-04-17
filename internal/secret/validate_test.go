package secret

import (
	"testing"
)

func TestValidateKey_Valid(t *testing.T) {
	opts := DefaultValidateOptions()
	if err := ValidateKey("MY_KEY", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateKey_Empty(t *testing.T) {
	opts := DefaultValidateOptions()
	if err := ValidateKey("", opts); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestValidateKey_InvalidPattern(t *testing.T) {
	opts := DefaultValidateOptions()
	if err := ValidateKey("lower_case", opts); err == nil {
		t.Fatal("expected error for lowercase key")
	}
}

func TestValidateValue_TooShort(t *testing.T) {
	opts := DefaultValidateOptions()
	opts.MinLength = 5
	if err := ValidateValue("abc", opts); err == nil {
		t.Fatal("expected error for short value")
	}
}

func TestValidateValue_TooLong(t *testing.T) {
	opts := DefaultValidateOptions()
	opts.MaxLength = 5
	if err := ValidateValue("toolongvalue", opts); err == nil {
		t.Fatal("expected error for long value")
	}
}

func TestValidateValue_Valid(t *testing.T) {
	opts := DefaultValidateOptions()
	if err := ValidateValue("somevalue", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateMap_AllValid(t *testing.T) {
	opts := DefaultValidateOptions()
	secrets := map[string]string{"MY_KEY": "value", "OTHER_KEY": "data"}
	if errs := ValidateMap(secrets, opts); errs != nil {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateMap_SomeInvalid(t *testing.T) {
	opts := DefaultValidateOptions()
	secrets := map[string]string{"valid_key": "v", "GOOD": "ok"}
	errs := ValidateMap(secrets, opts)
	if errs == nil {
		t.Fatal("expected errors")
	}
	if _, ok := errs["valid_key"]; !ok {
		t.Error("expected error for valid_key")
	}
}
