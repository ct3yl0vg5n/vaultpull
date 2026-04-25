package secret

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultGroupOptions returns sensible defaults for grouping secrets.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		Delimiter: "_",
		MaxDepth:  1,
	}
}

// GroupOptions controls how secrets are grouped.
type GroupOptions struct {
	// Delimiter separates the prefix from the rest of the key.
	Delimiter string
	// MaxDepth controls how many delimiter-separated segments form the group key.
	MaxDepth int
}

// GroupResult holds the grouped secrets and a summary.
type GroupResult struct {
	Groups map[string]map[string]string
	Order  []string // sorted group names
}

// Group partitions a flat map of secrets into named groups based on key prefixes.
// Keys without a delimiter are placed in the "default" group.
func Group(secrets map[string]string, opts GroupOptions) GroupResult {
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}
	if opts.MaxDepth < 1 {
		opts.MaxDepth = 1
	}

	groups := make(map[string]map[string]string)

	for k, v := range secrets {
		groupKey := extractGroupKey(k, opts.Delimiter, opts.MaxDepth)
		if _, ok := groups[groupKey]; !ok {
			groups[groupKey] = make(map[string]string)
		}
		groups[groupKey][k] = v
	}

	order := make([]string, 0, len(groups))
	for g := range groups {
		order = append(order, g)
	}
	sort.Strings(order)

	return GroupResult{Groups: groups, Order: order}
}

func extractGroupKey(key, delimiter string, maxDepth int) string {
	parts := strings.SplitN(key, delimiter, maxDepth+1)
	if len(parts) <= 1 {
		return "default"
	}
	return strings.Join(parts[:maxDepth], delimiter)
}

// FormatGroupReport returns a human-readable summary of grouped secrets.
func FormatGroupReport(result GroupResult) string {
	if len(result.Groups) == 0 {
		return "no secrets to group\n"
	}
	var sb strings.Builder
	for _, g := range result.Order {
		keys := make([]string, 0, len(result.Groups[g]))
		for k := range result.Groups[g] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sb.WriteString(fmt.Sprintf("[%s] (%d keys)\n", g, len(keys)))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}
	return sb.String()
}
