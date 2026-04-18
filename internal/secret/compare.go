package secret

import "time"

// CompareResult represents the result of comparing two secret values.
type CompareResult struct {
	Key       string
	OldValue  string
	NewValue  string
	Changed   bool
	ChangedAt time.Time
}

// CompareOptions controls comparison behavior.
type CompareOptions struct {
	CaseSensitive bool
	IgnoreKeys    []string
}

// DefaultCompareOptions returns sensible defaults.
func DefaultCompareOptions() CompareOptions {
	return CompareOptions{
		CaseSensitive: true,
	}
}

// CompareMap compares two maps of secrets and returns per-key results.
func CompareMap(old, new map[string]string, opts CompareOptions) []CompareResult {
	ignore := make(map[string]bool, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignore[k] = true
	}

	seen := make(map[string]bool)
	var results []CompareResult
	now := time.Now()

	for k, newVal := range new {
		if ignore[k] {
			continue
		}
		seen[k] = true
		oldVal, exists := old[k]
		changed := !exists || (opts.CaseSensitive && oldVal != newVal)
		results = append(results, CompareResult{
			Key:       k,
			OldValue:  oldVal,
			NewValue:  newVal,
			Changed:   changed,
			ChangedAt: now,
		})
	}

	for k, oldVal := range old {
		if ignore[k] || seen[k] {
			continue
		}
		results = append(results, CompareResult{
			Key:      k,
			OldValue: oldVal,
			Changed:  true,
			ChangedAt: now,
		})
	}

	return results
}
