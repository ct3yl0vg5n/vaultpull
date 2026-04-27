package secret

import (
	"fmt"
	"strings"
)

// TruncateOptions controls how values are truncated.
type TruncateOptions struct {
	MaxLength  int
	Suffix     string
	OnlyValues bool
	SkipKeys   []string
}

// DefaultTruncateOptions returns sensible defaults.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		MaxLength:  64,
		Suffix:     "...",
		OnlyValues: true,
	}
}

// TruncateResult holds a single truncation result.
type TruncateResult struct {
	Key       string
	Original  string
	Truncated string
	Changed   bool
}

// TruncateMap truncates values (and optionally keys) in the provided map.
func TruncateMap(src map[string]string, opts TruncateOptions) (map[string]string, []TruncateResult) {
	skip := make(map[string]bool, len(opts.SkipKeys))
	for _, k := range opts.SkipKeys {
		skip[k] = true
	}

	out := make(map[string]string, len(src))
	var results []TruncateResult

	for k, v := range src {
		if skip[k] {
			out[k] = v
			continue
		}
		truncated := truncateString(v, opts.MaxLength, opts.Suffix)
		changed := truncated != v
		out[k] = truncated
		results = append(results, TruncateResult{
			Key:       k,
			Original:  v,
			Truncated: truncated,
			Changed:   changed,
		})
	}
	return out, results
}

func truncateString(s string, max int, suffix string) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	cutAt := max - len(suffix)
	if cutAt < 0 {
		cutAt = 0
	}
	return s[:cutAt] + suffix
}

// FormatTruncateReport formats a human-readable summary of truncation results.
func FormatTruncateReport(results []TruncateResult) string {
	var sb strings.Builder
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
			sb.WriteString(fmt.Sprintf("  ~ %s: %q -> %q\n", r.Key, r.Original, r.Truncated))
		}
	}
	if changed == 0 {
		return "No values truncated.\n"
	}
	header := fmt.Sprintf("Truncated %d value(s):\n", changed)
	return header + sb.String()
}
