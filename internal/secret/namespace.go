package secret

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultNamespaceOptions returns sensible defaults for namespace operations.
func DefaultNamespaceOptions() NamespaceOptions {
	return NamespaceOptions{
		Separator: "__",
		Uppercase: true,
	}
}

// NamespaceOptions configures how secrets are namespaced.
type NamespaceOptions struct {
	Namespace string
	Separator string
	Uppercase bool
	StripOnRead bool
}

// NamespaceResult holds the outcome of a namespace operation.
type NamespaceResult struct {
	Key    string
	NewKey string
	Value  string
}

// ApplyNamespace prefixes all keys in src with the given namespace.
// If StripOnRead is true, it instead strips the namespace prefix from keys.
func ApplyNamespace(src map[string]string, opts NamespaceOptions) (map[string]string, []NamespaceResult, error) {
	if opts.Namespace == "" {
		return nil, nil, fmt.Errorf("namespace: namespace must not be empty")
	}
	if opts.Separator == "" {
		opts.Separator = "__"
	}

	prefix := opts.Namespace
	if opts.Uppercase {
		prefix = strings.ToUpper(prefix)
	}
	full := prefix + opts.Separator

	out := make(map[string]string, len(src))
	var results []NamespaceResult

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := src[k]
		var newKey string
		if opts.StripOnRead {
			if !strings.HasPrefix(k, full) {
				continue
			}
			newKey = strings.TrimPrefix(k, full)
		} else {
			newKey = full + k
		}
		out[newKey] = v
		results = append(results, NamespaceResult{Key: k, NewKey: newKey, Value: v})
	}
	return out, results, nil
}

// FormatNamespaceReport returns a human-readable summary of namespace results.
func FormatNamespaceReport(results []NamespaceResult, strip bool) string {
	if len(results) == 0 {
		return "namespace: no keys affected\n"
	}
	var sb strings.Builder
	action := "namespaced"
	if strip {
		action = "stripped"
	}
	fmt.Fprintf(&sb, "namespace: %d key(s) %s\n", len(results), action)
	for _, r := range results {
		fmt.Fprintf(&sb, "  %s -> %s\n", r.Key, r.NewKey)
	}
	return sb.String()
}
