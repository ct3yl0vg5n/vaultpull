package secret

import "fmt"

// LintOptions configures linting rules for secret keys and values.
type LintOptions struct {
	DisallowedPrefixes []string
	DisallowedSuffixes []string
	WarnOnLowercase    bool
	WarnOnWhitespace   bool
}

// LintResult holds a single lint warning for a key.
type LintResult struct {
	Key     string
	Message string
}

// DefaultLintOptions returns sensible lint defaults.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		DisallowedPrefixes: []string{"TMP_", "DEBUG_"},
		DisallowedSuffixes: []string{"_OLD", "_BAK"},
		WarnOnLowercase:    true,
		WarnOnWhitespace:   true,
	}
}

// LintMap runs lint checks against all keys and values in a map.
func LintMap(secrets map[string]string, opts LintOptions) []LintResult {
	var results []LintResult
	for k, v := range secrets {
		for _, p := range opts.DisallowedPrefixes {
			if len(k) >= len(p) && k[:len(p)] == p {
				results = append(results, LintResult{Key: k, Message: fmt.Sprintf("key has disallowed prefix %q", p)})
			}
		}
		for _, s := range opts.DisallowedSuffixes {
			if len(k) >= len(s) && k[len(k)-len(s):] == s {
				results = append(results, LintResult{Key: k, Message: fmt.Sprintf("key has disallowed suffix %q", s)})
			}
		}
		if opts.WarnOnLowercase {
			for _, c := range k {
				if c >= 'a' && c <= 'z' {
					results = append(results, LintResult{Key: k, Message: "key contains lowercase characters"})
					break
				}
			}
		}
		if opts.WarnOnWhitespace {
			for _, c := range v {
				if c == ' ' || c == '\t' {
					results = append(results, LintResult{Key: k, Message: "value contains leading/embedded whitespace"})
					break
				}
			}
		}
	}
	return results
}
