package secret

import (
	"fmt"
	"strings"
)

// ImportFormat represents the source format for importing secrets.
type ImportFormat string

const (
	FormatEnv    ImportFormat = "env"
	FormatJSON   ImportFormat = "json"
	FormatExport ImportFormat = "export"
)

// ImportOptions controls how secrets are imported from external formats.
type ImportOptions struct {
	Format      ImportFormat
	StripExport bool // strip "export " prefix from shell export statements
	Overwrite   bool // overwrite existing keys
	IgnoreKeys  []string
}

// DefaultImportOptions returns sensible defaults for importing secrets.
func DefaultImportOptions() ImportOptions {
	return ImportOptions{
		Format:      FormatEnv,
		StripExport: true,
		Overwrite:   false,
	}
}

// ImportResult holds the outcome of an import operation.
type ImportResult struct {
	Imported int
	Skipped  int
	Errors   []string
}

// Import parses raw input lines according to opts and merges them into dst.
// Returns an ImportResult summarising what happened.
func Import(dst map[string]string, lines []string, opts ImportOptions) (ImportResult, error) {
	ignore := make(map[string]bool, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignore[k] = true
	}

	var result ImportResult
	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if opts.StripExport {
			line = strings.TrimPrefix(line, "export ")
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			result.Errors = append(result.Errors, fmt.Sprintf("line %d: invalid format %q", i+1, raw))
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"`)

		if ignore[key] {
			result.Skipped++
			continue
		}
		if _, exists := dst[key]; exists && !opts.Overwrite {
			result.Skipped++
			continue
		}
		dst[key] = val
		result.Imported++
	}
	return result, nil
}

// FormatImportReport returns a human-readable summary of an ImportResult.
func FormatImportReport(r ImportResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "import: %d imported, %d skipped", r.Imported, r.Skipped)
	if len(r.Errors) > 0 {
		fmt.Fprintf(&sb, ", %d error(s):\n", len(r.Errors))
		for _, e := range r.Errors {
			fmt.Fprintf(&sb, "  - %s\n", e)
		}
	} else {
		sb.WriteByte('\n')
	}
	return sb.String()
}
