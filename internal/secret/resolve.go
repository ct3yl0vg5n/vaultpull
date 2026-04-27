package secret

import (
	"fmt"
	"strings"
)

// ResolveOptions controls how secret references are resolved.
type ResolveOptions struct {
	// Prefix marks a value as a reference, e.g. "ref:"
	RefPrefix string
	// Strict causes Resolve to return an error if a referenced key is missing.
	Strict bool
}

// DefaultResolveOptions returns sensible defaults.
func DefaultResolveOptions() ResolveOptions {
	return ResolveOptions{
		RefPrefix: "ref:",
		Strict:    true,
	}
}

// ResolveResult holds the outcome for a single key.
type ResolveResult struct {
	Key      string
	Original string
	Resolved string
	WasRef   bool
	Missing  bool
}

// Resolve walks src and replaces values that start with RefPrefix with the
// value of the referenced key, also looked up in src.
func Resolve(src map[string]string, opts ResolveOptions) (map[string]string, []ResolveResult, error) {
	out := make(map[string]string, len(src))
	results := make([]ResolveResult, 0)

	for k, v := range src {
		if !strings.HasPrefix(v, opts.RefPrefix) {
			out[k] = v
			results = append(results, ResolveResult{Key: k, Original: v, Resolved: v, WasRef: false})
			continue
		}

		refKey := strings.TrimPrefix(v, opts.RefPrefix)
		resolved, ok := src[refKey]
		if !ok {
			if opts.Strict {
				return nil, nil, fmt.Errorf("resolve: key %q references missing key %q", k, refKey)
			}
			out[k] = v
			results = append(results, ResolveResult{Key: k, Original: v, Resolved: v, WasRef: true, Missing: true})
			continue
		}

		out[k] = resolved
		results = append(results, ResolveResult{Key: k, Original: v, Resolved: resolved, WasRef: true})
	}

	return out, results, nil
}

// FormatResolveReport returns a human-readable summary of resolve results.
func FormatResolveReport(results []ResolveResult) string {
	var sb strings.Builder
	refs, missing := 0, 0
	for _, r := range results {
		if r.WasRef {
			refs++
		}
		if r.Missing {
			missing++
			sb.WriteString(fmt.Sprintf("  [missing]  %s -> %s\n", r.Key, r.Original))
		} else if r.WasRef {
			sb.WriteString(fmt.Sprintf("  [resolved] %s -> %s\n", r.Key, r.Resolved))
		}
	}
	header := fmt.Sprintf("Resolve: %d reference(s), %d missing\n", refs, missing)
	return header + sb.String()
}
