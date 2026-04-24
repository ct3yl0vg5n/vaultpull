package secret

import (
	"fmt"
	"strings"
	"time"
)

// WatchOptions configures the watch/poll behaviour.
type WatchOptions struct {
	Interval    time.Duration
	MaxChecks   int // 0 = unlimited
	AlertOnNew  bool
	AlertOnDiff bool
}

// WatchEvent represents a single detected change during a watch cycle.
type WatchEvent struct {
	Key       string
	Kind      string // "added", "removed", "modified"
	OldValue  string
	NewValue  string
	DetectedAt time.Time
}

// DefaultWatchOptions returns sensible defaults.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Interval:    30 * time.Second,
		MaxChecks:   0,
		AlertOnNew:  true,
		AlertOnDiff: true,
	}
}

// DetectChanges compares a previous snapshot with the current secrets and
// returns a slice of WatchEvents describing what changed.
func DetectChanges(prev, curr map[string]string, opts WatchOptions) []WatchEvent {
	now := time.Now().UTC()
	var events []WatchEvent

	for k, newVal := range curr {
		oldVal, existed := prev[k]
		switch {
		case !existed && opts.AlertOnNew:
			events = append(events, WatchEvent{Key: k, Kind: "added", NewValue: newVal, DetectedAt: now})
		case existed && oldVal != newVal && opts.AlertOnDiff:
			events = append(events, WatchEvent{Key: k, Kind: "modified", OldValue: oldVal, NewValue: newVal, DetectedAt: now})
		}
	}

	for k, oldVal := range prev {
		if _, exists := curr[k]; !exists {
			events = append(events, WatchEvent{Key: k, Kind: "removed", OldValue: oldVal, DetectedAt: now})
		}
	}

	return events
}

// FormatWatchReport renders a human-readable report of watch events.
func FormatWatchReport(events []WatchEvent) string {
	if len(events) == 0 {
		return "no changes detected"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d change(s) detected:\n", len(events)))
	for _, e := range events {
		switch e.Kind {
		case "added":
			sb.WriteString(fmt.Sprintf("  + %s (added)\n", e.Key))
		case "removed":
			sb.WriteString(fmt.Sprintf("  - %s (removed)\n", e.Key))
		case "modified":
			sb.WriteString(fmt.Sprintf("  ~ %s (modified)\n", e.Key))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
