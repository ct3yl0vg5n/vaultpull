package secret

import (
	"fmt"
	"strings"
)

// RenameOptions controls key renaming behaviour.
type RenameOptions struct {
	// Rules maps old key names to new key names.
	Rules map[string]string
	// DryRun reports what would change without mutating dst.
	DryRun bool
	// ErrorOnMissing returns an error when a source key is not found.
	ErrorOnMissing bool
}

// DefaultRenameOptions returns sensible defaults.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		Rules:          map[string]string{},
		DryRun:         false,
		ErrorOnMissing: false,
	}
}

// RenameResult describes a single rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Applied bool
	Skipped bool // key not found in src
}

// Rename applies renaming rules to src, returning a new map and a report.
func Rename(src map[string]string, opts RenameOptions) (map[string]string, []RenameResult, error) {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}

	var results []RenameResult

	for oldKey, newKey := range opts.Rules {
		val, exists := dst[oldKey]
		if !exists {
			if opts.ErrorOnMissing {
				return nil, nil, fmt.Errorf("rename: key %q not found in source", oldKey)
			}
			results = append(results, RenameResult{OldKey: oldKey, NewKey: newKey, Skipped: true})
			continue
		}
		if !opts.DryRun {
			delete(dst, oldKey)
			dst[newKey] = val
		}
		results = append(results, RenameResult{OldKey: oldKey, NewKey: newKey, Applied: !opts.DryRun})
	}

	return dst, results, nil
}

// FormatRenameReport returns a human-readable summary of rename results.
func FormatRenameReport(results []RenameResult) string {
	if len(results) == 0 {
		return "no rename rules applied"
	}
	var sb strings.Builder
	applied, skipped := 0, 0
	for _, r := range results {
		switch {
		case r.Skipped:
			fmt.Fprintf(&sb, "  SKIP    %s -> %s (key not found)\n", r.OldKey, r.NewKey)
			skipped++
		case r.Applied:
			fmt.Fprintf(&sb, "  RENAMED %s -> %s\n", r.OldKey, r.NewKey)
			applied++
		default:
			fmt.Fprintf(&sb, "  DRY-RUN %s -> %s\n", r.OldKey, r.NewKey)
		}
	}
	fmt.Fprintf(&sb, "summary: %d renamed, %d skipped", applied, skipped)
	return sb.String()
}
