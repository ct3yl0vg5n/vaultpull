package secret

import (
	"errors"
	"fmt"
	"regexp"
)

type ValidateOptions struct {
	MinLength  int
	MaxLength  int
	RequireKey bool
	KeyPattern string
}

func DefaultValidateOptions() ValidateOptions {
	return ValidateOptions{
		MinLength:  1,
		MaxLength:  4096,
		RequireKey: true,
		KeyPattern: `^[A-Z][A-Z0-9_]*$`,
	}
}

func ValidateKey(key string, opts ValidateOptions) error {
	if opts.RequireKey && key == "" {
		return errors.New("key must not be empty")
	}
	if opts.KeyPattern != "" && key != "" {
		re, err := regexp.Compile(opts.KeyPattern)
		if err != nil {
			return fmt.Errorf("invalid key pattern: %w", err)
		}
		if !re.MatchString(key) {
			return fmt.Errorf("key %q does not match pattern %q", key, opts.KeyPattern)
		}
	}
	return nil
}

func ValidateValue(value string, opts ValidateOptions) error {
	if len(value) < opts.MinLength {
		return fmt.Errorf("value length %d is below minimum %d", len(value), opts.MinLength)
	}
	if opts.MaxLength > 0 && len(value) > opts.MaxLength {
		return fmt.Errorf("value length %d exceeds maximum %d", len(value), opts.MaxLength)
	}
	return nil
}

func ValidateMap(secrets map[string]string, opts ValidateOptions) map[string]error {
	errs := make(map[string]error)
	for k, v := range secrets {
		if err := ValidateKey(k, opts); err != nil {
			errs[k] = err
			continue
		}
		if err := ValidateValue(v, opts); err != nil {
			errs[k] = err
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
