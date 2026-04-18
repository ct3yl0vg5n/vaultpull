package secret

import (
	"fmt"
	"strings"
)

// DiffReportOptions controls how secret diffs are reported.
type DiffReportOptions struct {
	RedactValues bool
	RevealChars  int
}

// DefaultDiffReportOptions returns sensible defaults.
func DefaultDiffReportOptions() DiffReportOptions {
	return DiffReportOptions{
		RedactValues: true,
		RevealChars:  4,
	}
}

// DiffEntry represents a single key change between two secret maps.
type DiffEntry struct {
	Key    string
	Old    string
	New    string
	Status string // added, removed, modified, unchanged
}

// BuildDiffReport compares two maps and returns a slice of DiffEntry.
func BuildDiffReport(old, next map[string]string, opts DiffReportOptions) []DiffEntry {
	seen := map[string]bool{}
	var entries []DiffEntry

	for k, newVal := range next {
		seen[k] = true
		oldVal, exists := old[k]
		switch {
		case !exists:
			entries = append(entries, DiffEntry{Key: k, Old: "", New: maybeRedact(newVal, opts), Status: "added"})
		case oldVal != newVal:
			entries = append(entries, DiffEntry{Key: k, Old: maybeRedact(oldVal, opts), New: maybeRedact(newVal, opts), Status: "modified"})
		default:
			entries = append(entries, DiffEntry{Key: k, Old: maybeRedact(oldVal, opts), New: maybeRedact(newVal, opts), Status: "unchanged"})
		}
	}

	for k, oldVal := range old {
		if !seen[k] {
			entries = append(entries, DiffEntry{Key: k, Old: maybeRedact(oldVal, opts), New: "", Status: "removed"})
		}
	}
	return entries
}

// FormatDiffReport renders a human-readable diff report.
func FormatDiffReport(entries []DiffEntry) string {
	if len(entries) == 0 {
		return "no differences found"
	}
	var sb strings.Builder
	for _, e := range entries {
		switch e.Status {
		case "added":
			sb.WriteString(fmt.Sprintf("+ %s=%s\n", e.Key, e.New))
		case "removed":
			sb.WriteString(fmt.Sprintf("- %s=%s\n", e.Key, e.Old))
		case "modified":
			sb.WriteString(fmt.Sprintf("~ %s: %s -> %s\n", e.Key, e.Old, e.New))
		case "unchanged":
			sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, e.New))
		}
	}
	return sb.String()
}

func maybeRedact(val string, opts DiffReportOptions) string {
	if !opts.RedactValues {
		return val
	}
	return Redact(val, DefaultRedactOptions())
}
