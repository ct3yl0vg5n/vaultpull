package secret

import (
	"fmt"
	"sort"
	"strings"
)

// ScopeOptions controls how secrets are partitioned into named scopes.
type ScopeOptions struct {
	// Scopes maps scope name -> key prefixes that belong to it.
	Scopes map[string][]string
	// DefaultScope receives keys that match no explicit scope.
	DefaultScope string
}

// DefaultScopeOptions returns sensible defaults.
func DefaultScopeOptions() ScopeOptions {
	return ScopeOptions{
		Scopes:       map[string][]string{},
		DefaultScope: "default",
	}
}

// ScopeResult holds the partitioned secrets.
type ScopeResult struct {
	Name    string
	Secrets map[string]string
}

// Partition groups secrets from src into scopes defined by opts.
// Keys that match multiple prefixes are assigned to the first matching scope
// (iteration order is alphabetical for determinism).
func Partition(src map[string]string, opts ScopeOptions) []ScopeResult {
	results := map[string]*ScopeResult{}

	// Collect ordered scope names for deterministic matching.
	scopeNames := make([]string, 0, len(opts.Scopes))
	for name := range opts.Scopes {
		scopeNames = append(scopeNames, name)
	}
	sort.Strings(scopeNames)

	for key, val := range src {
		assigned := false
		for _, name := range scopeNames {
			for _, prefix := range opts.Scopes[name] {
				if strings.HasPrefix(key, prefix) {
					if results[name] == nil {
						results[name] = &ScopeResult{Name: name, Secrets: map[string]string{}}
					}
					results[name].Secrets[key] = val
					assigned = true
					break
				}
			}
			if assigned {
				break
			}
		}
		if !assigned {
			def := opts.DefaultScope
			if results[def] == nil {
				results[def] = &ScopeResult{Name: def, Secrets: map[string]string{}}
			}
			results[def].Secrets[key] = val
		}
	}

	ordered := make([]ScopeResult, 0, len(results))
	names := make([]string, 0, len(results))
	for n := range results {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		ordered = append(ordered, *results[n])
	}
	return ordered
}

// FormatScopeReport returns a human-readable summary of scope partitioning.
func FormatScopeReport(results []ScopeResult) string {
	if len(results) == 0 {
		return "no secrets to partition\n"
	}
	var sb strings.Builder
	for _, r := range results {
		fmt.Fprintf(&sb, "[%s] %d key(s)\n", r.Name, len(r.Secrets))
		keys := make([]string, 0, len(r.Secrets))
		for k := range r.Secrets {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&sb, "  %s\n", k)
		}
	}
	return sb.String()
}
