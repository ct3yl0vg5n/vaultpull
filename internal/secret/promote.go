package secret

import (
	"fmt"
	"sort"
	"strings"
)

// PromoteOptions controls how secrets are promoted between environments.
type PromoteOptions struct {
	// FromEnv is the source environment label (e.g. "staging").
	FromEnv string
	// ToEnv is the destination environment label (e.g. "production").
	ToEnv string
	// DryRun previews changes without applying them.
	DryRun bool
	// OverwriteExisting replaces keys that already exist in the destination.
	OverwriteExisting bool
	// IgnoreKeys is a set of keys to skip during promotion.
	IgnoreKeys []string
}

// DefaultPromoteOptions returns sensible defaults.
func DefaultPromoteOptions() PromoteOptions {
	return PromoteOptions{
		FromEnv:           "staging",
		ToEnv:             "production",
		DryRun:            false,
		OverwriteExisting: false,
	}
}

// PromoteResult describes the outcome of a single key promotion.
type PromoteResult struct {
	Key      string
	Action   string // "added", "overwritten", "skipped", "ignored"
	OldValue string
	NewValue string
}

// Promote copies keys from src into dst according to opts.
// It returns the merged destination map and a slice of results.
func Promote(src, dst map[string]string, opts PromoteOptions) (map[string]string, []PromoteResult) {
	ignored := make(map[string]struct{}, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignored[k] = struct{}{}
	}

	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var results []PromoteResult
	keys := sortedMapKeys(src)
	for _, k := range keys {
		v := src[k]
		if _, skip := ignored[k]; skip {
			results = append(results, PromoteResult{Key: k, Action: "ignored"})
			continue
		}
		existing, exists := out[k]
		if exists && !opts.OverwriteExisting {
			results = append(results, PromoteResult{Key: k, Action: "skipped", OldValue: existing, NewValue: v})
			continue
		}
		action := "added"
		if exists {
			action = "overwritten"
		}
		if !opts.DryRun {
			out[k] = v
		}
		results = append(results, PromoteResult{Key: k, Action: action, OldValue: existing, NewValue: v})
	}
	return out, results
}

// FormatPromoteReport renders a human-readable summary of promotion results.
func FormatPromoteReport(from, to string, results []PromoteResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Promote: %s → %s\n", from, to)
	counts := map[string]int{}
	for _, r := range results {
		counts[r.Action]++
		switch r.Action {
		case "added":
			fmt.Fprintf(&sb, "  [+] %s\n", r.Key)
		case "overwritten":
			fmt.Fprintf(&sb, "  [~] %s\n", r.Key)
		case "skipped":
			fmt.Fprintf(&sb, "  [-] %s (skipped, already exists)\n", r.Key)
		case "ignored":
			fmt.Fprintf(&sb, "  [x] %s (ignored)\n", r.Key)
		}
	}
	fmt.Fprintf(&sb, "Summary: %d added, %d overwritten, %d skipped, %d ignored\n",
		counts["added"], counts["overwritten"], counts["skipped"], counts["ignored"])
	return sb.String()
}

func sortedMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
