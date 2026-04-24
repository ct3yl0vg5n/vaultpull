package secret

import (
	"fmt"
	"sort"
	"strings"
)

// MergeStrategy defines how conflicts are resolved when merging two secret maps.
type MergeStrategy string

const (
	MergeStrategyOurs   MergeStrategy = "ours"   // keep local value on conflict
	MergeStrategyTheirs MergeStrategy = "theirs" // take incoming value on conflict
	MergeStrategyError  MergeStrategy = "error"  // return error on conflict
)

// MergeOptions configures the Merge operation.
type MergeOptions struct {
	Strategy   MergeStrategy
	IgnoreKeys []string
}

// DefaultMergeOptions returns sensible defaults.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		Strategy: MergeStrategyTheirs,
	}
}

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Merged    map[string]string
	Conflicts []string
	Added     []string
	Skipped   []string
}

// Merge combines base and incoming secret maps according to the given options.
func Merge(base, incoming map[string]string, opts MergeOptions) (MergeResult, error) {
	ignore := make(map[string]struct{}, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignore[k] = struct{}{}
	}

	result := MergeResult{
		Merged: make(map[string]string, len(base)),
	}
	for k, v := range base {
		result.Merged[k] = v
	}

	for k, inVal := range incoming {
		if _, skip := ignore[k]; skip {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		baseVal, exists := base[k]
		if !exists {
			result.Merged[k] = inVal
			result.Added = append(result.Added, k)
			continue
		}
		if baseVal == inVal {
			continue
		}
		result.Conflicts = append(result.Conflicts, k)
		switch opts.Strategy {
		case MergeStrategyTheirs:
			result.Merged[k] = inVal
		case MergeStrategyOurs:
			// keep existing value already in Merged
		case MergeStrategyError:
			return MergeResult{}, fmt.Errorf("merge conflict on key %q", k)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Conflicts)
	sort.Strings(result.Skipped)
	return result, nil
}

// FormatMergeReport returns a human-readable summary of the merge result.
func FormatMergeReport(r MergeResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("merged: %d keys total\n", len(r.Merged)))
	if len(r.Added) > 0 {
		sb.WriteString(fmt.Sprintf("  added (%d): %s\n", len(r.Added), strings.Join(r.Added, ", ")))
	}
	if len(r.Conflicts) > 0 {
		sb.WriteString(fmt.Sprintf("  conflicts (%d): %s\n", len(r.Conflicts), strings.Join(r.Conflicts, ", ")))
	}
	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("  skipped (%d): %s\n", len(r.Skipped), strings.Join(r.Skipped, ", ")))
	}
	return strings.TrimRight(sb.String(), "\n")
}
