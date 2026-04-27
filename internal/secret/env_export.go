package secret

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultExportOptions returns sensible defaults for env export.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Prefix:      "",
		ExportDecl:  false,
		QuoteValues: false,
		SkipEmpty:   false,
		SortKeys:    true,
	}
}

// ExportOptions controls how secrets are rendered as env output.
type ExportOptions struct {
	Prefix      string
	ExportDecl  bool
	QuoteValues bool
	SkipEmpty   bool
	SortKeys    bool
}

// ExportResult holds a single rendered line and its source key.
type ExportResult struct {
	Key  string
	Line string
}

// Export converts a map of secrets into env-formatted lines.
func Export(src map[string]string, opts ExportOptions) []ExportResult {
	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	if opts.SortKeys {
		sort.Strings(keys)
	}

	results := make([]ExportResult, 0, len(keys))
	for _, k := range keys {
		v := src[k]
		if opts.SkipEmpty && v == "" {
			continue
		}
		finalKey := opts.Prefix + k
		var line string
		if opts.QuoteValues {
			v = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
		}
		if opts.ExportDecl {
			line = fmt.Sprintf("export %s=%s", finalKey, v)
		} else {
			line = fmt.Sprintf("%s=%s", finalKey, v)
		}
		results = append(results, ExportResult{Key: finalKey, Line: line})
	}
	return results
}

// FormatExport renders ExportResults as a newline-joined string.
func FormatExport(results []ExportResult) string {
	lines := make([]string, len(results))
	for i, r := range results {
		lines[i] = r.Line
	}
	return strings.Join(lines, "\n")
}
