package filter

import (
	"strings"
)

// Rule defines how secrets should be filtered before writing to .env files.
type Rule struct {
	// Prefix, if set, only includes keys that start with this prefix.
	Prefix string
	// Exclude is a list of exact key names to omit.
	Exclude []string
	// StripPrefix removes the prefix from key names when writing.
	StripPrefix bool
}

// Apply filters a map of secrets according to the rule and returns a new map.
func (r *Rule) Apply(secrets map[string]string) map[string]string {
	excludeSet := make(map[string]struct{}, len(r.Exclude))
	for _, k := range r.Exclude {
		excludeSet[k] = struct{}{}
	}

	result := make(map[string]string)
	for k, v := range secrets {
		if _, excluded := excludeSet[k]; excluded {
			continue
		}
		if r.Prefix != "" && !strings.HasPrefix(k, r.Prefix) {
			continue
		}
		outKey := k
		if r.StripPrefix && r.Prefix != "" {
			outKey = strings.TrimPrefix(k, r.Prefix)
			if outKey == "" {
				continue
			}
		}
		result[outKey] = v
	}
	return result
}

// IsEmpty returns true when the rule has no filtering constraints.
func (r *Rule) IsEmpty() bool {
	return r.Prefix == "" && len(r.Exclude) == 0
}
