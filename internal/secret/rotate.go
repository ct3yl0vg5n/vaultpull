package secret

import (
	"fmt"
	"strings"
	"time"
)

// RotateOptions controls secret rotation behaviour.
type RotateOptions struct {
	MaxAgeDays int
	DryRun     bool
}

// DefaultRotateOptions returns sensible defaults.
func DefaultRotateOptions() RotateOptions {
	return RotateOptions{
		MaxAgeDays: 90,
		DryRun:     false,
	}
}

// RotateEntry represents a single secret that may need rotation.
type RotateEntry struct {
	Key       string
	CreatedAt time.Time
}

// RotateResult holds the outcome for a single entry.
type RotateResult struct {
	Key     string
	Rotated bool
	Reason  string
}

// CheckRotation evaluates which entries exceed the max age and marks them for rotation.
func CheckRotation(entries []RotateEntry, opts RotateOptions) []RotateResult {
	cutoff := time.Now().AddDate(0, 0, -opts.MaxAgeDays)
	results := make([]RotateResult, 0, len(entries))
	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		r := RotateResult{Key: e.Key}
		if e.CreatedAt.Before(cutoff) {
			r.Rotated = !opts.DryRun
			age := int(time.Since(e.CreatedAt).Hours() / 24)
			r.Reason = fmt.Sprintf("age %d days exceeds limit %d", age, opts.MaxAgeDays)
		}
		results = append(results, r)
	}
	return results
}

// FormatRotateReport returns a human-readable summary of rotation results.
func FormatRotateReport(results []RotateResult) string {
	var sb strings.Builder
	rotated, skipped := 0, 0
	for _, r := range results {
		if r.Rotated {
			rotated++
			fmt.Fprintf(&sb, "  [ROTATE] %s — %s\n", r.Key, r.Reason)
		} else if r.Reason != "" {
			fmt.Fprintf(&sb, "  [DRY-RUN] %s — %s\n", r.Key, r.Reason)
		} else {
			skipped++
		}
	}
	fmt.Fprintf(&sb, "\nRotated: %d  Skipped: %d\n", rotated, skipped)
	return sb.String()
}
