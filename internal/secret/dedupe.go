package secret

import (
	"fmt"
	"sort"
	"strings"
)

// DedupeOptions controls deduplication behaviour.
type DedupeOptions struct {
	// CaseSensitive determines whether key comparison is case-sensitive.
	CaseSensitive bool
	// PreferLast keeps the last occurrence of a duplicate key instead of the first.
	PreferLast bool
	// ReportOnly returns the duplicates without modifying the map.
	ReportOnly bool
}

// DefaultDedupeOptions returns sensible defaults.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{
		CaseSensitive: true,
		PreferLast:    false,
		ReportOnly:    false,
	}
}

// DedupeResult holds the output of a Dedupe operation.
type DedupeResult struct {
	Out        map[string]string
	Duplicates []string
}

// Dedupe removes duplicate keys from src according to opts.
// Because map iteration order is non-deterministic, callers should pass an
// ordered slice of keys via orderedKeys when PreferLast semantics matter.
func Dedupe(src map[string]string, orderedKeys []string, opts DedupeOptions) DedupeResult {
	seen := make(map[string]string) // normalised key -> original key
	out := make(map[string]string)
	var dupes []string

	keys := orderedKeys
	if len(keys) == 0 {
		for k := range src {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			continue
		}
		norm := k
		if !opts.CaseSensitive {
			norm = strings.ToUpper(k)
		}
		if orig, exists := seen[norm]; exists {
			dupes = append(dupes, k)
			if opts.PreferLast && !opts.ReportOnly {
				delete(out, orig)
				out[k] = v
				seen[norm] = k
			}
			continue
		}
		seen[norm] = k
		if !opts.ReportOnly {
			out[k] = v
		}
	}

	if opts.ReportOnly {
		out = src
	}

	return DedupeResult{Out: out, Duplicates: dupes}
}

// FormatDedupeReport returns a human-readable summary.
func FormatDedupeReport(r DedupeResult) string {
	if len(r.Duplicates) == 0 {
		return "no duplicate keys found"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "found %d duplicate key(s):\n", len(r.Duplicates))
	for _, d := range r.Duplicates {
		fmt.Fprintf(&sb, "  - %s\n", d)
	}
	return strings.TrimRight(sb.String(), "\n")
}
