package secret

import (
	"regexp"
	"sort"
	"strings"
)

// SearchOptions controls how secrets are searched.
type SearchOptions struct {
	// CaseSensitive controls whether key/value matching is case-sensitive.
	CaseSensitive bool
	// SearchValues includes secret values in the search in addition to keys.
	SearchValues bool
	// Regex treats the query as a regular expression.
	Regex bool
}

// DefaultSearchOptions returns sensible defaults.
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		CaseSensitive: false,
		SearchValues:  false,
		Regex:         false,
	}
}

// SearchResult holds a single matched secret.
type SearchResult struct {
	Key        string
	Value      string
	MatchedKey bool
	MatchedVal bool
}

// Search finds secrets whose keys (and optionally values) match the query.
func Search(secrets map[string]string, query string, opts SearchOptions) ([]SearchResult, error) {
	var matchFn func(s string) bool

	if opts.Regex {
		flags := ""
		if !opts.CaseSensitive {
			flags = "(?i)"
		}
		re, err := regexp.Compile(flags + query)
		if err != nil {
			return nil, err
		}
		matchFn = re.MatchString
	} else {
		needle := query
		if !opts.CaseSensitive {
			needle = strings.ToLower(query)
		}
		matchFn = func(s string) bool {
			if !opts.CaseSensitive {
				s = strings.ToLower(s)
			}
			return strings.Contains(s, needle)
		}
	}

	var results []SearchResult
	for k, v := range secrets {
		result := SearchResult{Key: k, Value: v}
		result.MatchedKey = matchFn(k)
		if opts.SearchValues {
			result.MatchedVal = matchFn(v)
		}
		if result.MatchedKey || result.MatchedVal {
			results = append(results, result)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results, nil
}

// FormatSearchReport returns a human-readable summary of search results.
func FormatSearchReport(results []SearchResult, query string) string {
	if len(results) == 0 {
		return "no matches found for: " + query
	}
	var sb strings.Builder
	sb.WriteString("search results for: " + query + "\n")
	for _, r := range results {
		tag := "[key]"
		if r.MatchedVal && !r.MatchedKey {
			tag = "[value]"
		} else if r.MatchedKey && r.MatchedVal {
			tag = "[key+value]"
		}
		sb.WriteString("  " + tag + " " + r.Key + "\n")
	}
	return sb.String()
}
