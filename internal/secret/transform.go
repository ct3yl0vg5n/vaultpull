package secret

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a transformation to apply to a secret key or value.
type Rule struct {
	UppercaseKeys bool
	KeyPrefix     string
	KeyReplace    *KeyReplace
	RedactValues  bool
}

// KeyReplace defines a find/replace on key names.
type KeyReplace struct {
	From string
	To   string
}

var invalidKeyChar = regexp.MustCompile(`[^A-Z0-9_]`)

// ApplyKey transforms a secret key according to the rule.
func ApplyKey(key string, r Rule) (string, error) {
	if r.UppercaseKeys {
		key = strings.ToUpper(key)
	}
	if r.KeyPrefix != "" {
		key = r.KeyPrefix + key
	}
	if r.KeyReplace != nil {
		key = strings.ReplaceAll(key, r.KeyReplace.From, r.KeyReplace.To)
	}
	if invalidKeyChar.MatchString(key) {
		return "", fmt.Errorf("invalid env key after transform: %q", key)
	}
	return key, nil
}

// ApplyValue optionally redacts a secret value.
func ApplyValue(value string, r Rule) string {
	if r.RedactValues {
		return "***REDACTED***"
	}
	return value
}

// ApplyMap transforms a map of secrets using the given rule.
func ApplyMap(secrets map[string]string, r Rule) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		newKey, err := ApplyKey(k, r)
		if err != nil {
			return nil, err
		}
		out[newKey] = ApplyValue(v, r)
	}
	return out, nil
}
