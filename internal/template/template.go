package template

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Rule defines how to transform a secret key into an env var name.
type Rule struct {
	// Prefix adds a string before the key.
	Prefix string
	// Upper forces the key to uppercase.
	Upper bool
	// Replace maps substrings: each pair [from, to].
	Replace [][2]string
}

var validKey = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Apply transforms a raw secret key according to the Rule.
func Apply(key string, r Rule) (string, error) {
	result := key

	for _, pair := range r.Replace {
		result = strings.ReplaceAll(result, pair[0], pair[1])
	}

	if r.Upper {
		result = strings.ToUpper(result)
	}

	if r.Prefix != "" {
		result = r.Prefix + result
	}

	if !validKey.MatchString(result) {
		return "", fmt.Errorf("template: resulting key %q is not a valid env var name", result)
	}

	return result, nil
}

// ApplyMap transforms all keys in a map according to the Rule.
// Keys that fail validation are skipped and a warning is printed to stderr.
func ApplyMap(secrets map[string]string, r Rule) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		newKey, err := Apply(k, r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: %v — skipping\n", err)
			continue
		}
		out[newKey] = v
	}
	return out
}
