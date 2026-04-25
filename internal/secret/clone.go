package secret

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultCloneOptions returns sensible defaults for cloning secrets.
func DefaultCloneOptions() CloneOptions {
	return CloneOptions{
		Overwrite:  false,
		DryRun:     false,
		KeyPrefix:  "",
		IgnoreKeys: []string{},
	}
}

// CloneOptions controls how secrets are cloned between maps.
type CloneOptions struct {
	Overwrite  bool
	DryRun     bool
	KeyPrefix  string
	IgnoreKeys []string
}

// CloneResult holds the outcome of a Clone operation.
type CloneResult struct {
	Copied  []string
	Skipped []string
	Dest    map[string]string
}

// Clone copies secrets from src into dst, applying prefix and overwrite rules.
func Clone(src, dst map[string]string, opts CloneOptions) CloneResult {
	ignored := make(map[string]bool, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignored[k] = true
	}

	result := CloneResult{
		Dest: make(map[string]string),
	}
	for k, v := range dst {
		result.Dest[k] = v
	}

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if ignored[k] {
			continue
		}
		destKey := opts.KeyPrefix + k
		if _, exists := result.Dest[destKey]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, destKey)
			continue
		}
		result.Copied = append(result.Copied, destKey)
		if !opts.DryRun {
			result.Dest[destKey] = src[k]
		}
	}
	return result
}

// FormatCloneReport returns a human-readable summary of a CloneResult.
func FormatCloneReport(r CloneResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Cloned: %d  Skipped: %d\n", len(r.Copied), len(r.Skipped)))
	for _, k := range r.Copied {
		sb.WriteString(fmt.Sprintf("  + %s\n", k))
	}
	for _, k := range r.Skipped {
		sb.WriteString(fmt.Sprintf("  ~ %s (skipped, already exists)\n", k))
	}
	return sb.String()
}
