package secret

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// RollbackOptions configures rollback behaviour.
type RollbackOptions struct {
	DryRun    bool
	MaxAge    time.Duration
	IgnoreKeys []string
}

// DefaultRollbackOptions returns sensible defaults.
func DefaultRollbackOptions() RollbackOptions {
	return RollbackOptions{
		DryRun: false,
		MaxAge: 7 * 24 * time.Hour,
	}
}

// RollbackEntry describes a single key being rolled back.
type RollbackEntry struct {
	Key      string
	OldValue string
	NewValue string
}

// RollbackResult holds the outcome of a rollback operation.
type RollbackResult struct {
	Applied []RollbackEntry
	Skipped []string
	DryRun  bool
}

// Rollback restores secrets in dst to the values found in snapshot,
// returning a RollbackResult describing what changed.
func Rollback(dst map[string]string, snapshot map[string]string, opts RollbackOptions) (map[string]string, RollbackResult) {
	ignore := make(map[string]bool, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignore[k] = true
	}

	result := RollbackResult{DryRun: opts.DryRun}
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	for key, snapVal := range snapshot {
		if ignore[key] {
			result.Skipped = append(result.Skipped, key)
			continue
		}
		current, exists := dst[key]
		if !exists || current != snapVal {
			entry := RollbackEntry{Key: key, OldValue: current, NewValue: snapVal}
			result.Applied = append(result.Applied, entry)
			if !opts.DryRun {
				out[key] = snapVal
			}
		}
	}

	sort.Slice(result.Applied, func(i, j int) bool {
		return result.Applied[i].Key < result.Applied[j].Key
	})
	sort.Strings(result.Skipped)
	return out, result
}

// FormatRollbackReport formats a RollbackResult as a human-readable string.
func FormatRollbackReport(r RollbackResult) string {
	var sb strings.Builder
	if r.DryRun {
		sb.WriteString("[dry-run] rollback preview\n")
	}
	if len(r.Applied) == 0 {
		sb.WriteString("no changes to roll back\n")
		return sb.String()
	}
	for _, e := range r.Applied {
		action := "applied"
		if r.DryRun {
			action = "would apply"
		}
		sb.WriteString(fmt.Sprintf("  %s  %s: %q -> %q\n", action, e.Key, e.OldValue, e.NewValue))
	}
	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("  skipped: %s\n", strings.Join(r.Skipped, ", ")))
	}
	return sb.String()
}
